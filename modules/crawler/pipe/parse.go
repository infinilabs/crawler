/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pipe

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/filter"
	"github.com/infinitbyte/gopa/core/model"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	"strings"
)

const ParsePage JointKey = "parse"

type ParsePageJoint struct {
	Parameters
	MaxPageOfBreadth map[int]int //max page to fetch in each level's breadth, eg: 1:100;2:50;3:5;4:1
	//TODO support save link,script
}

const dispatchLinks ParaKey = "dispatch_links"
const maxDepth ParaKey = "max_depth"
const maxBreadth ParaKey = "max_breadth"

func (this ParsePageJoint) Name() string {
	return string(ParsePage)
}

func (this ParsePageJoint) Process(context *Context) error {

	task := context.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := context.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	refUrl := task.Url
	refHost := task.Host
	depth := task.Depth
	breadth := task.Breadth
	fileContent := snapshot.Payload
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(fileContent))
	if err != nil {
		panic(err)
	}

	title := doc.Find("title").Text()
	if len(title) > 0 {
		snapshot.Title = util.NoWordBreak(util.XSSHandle(title))
	}

	links := map[string]string{}

	metadata := map[string]interface{}{}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, exist := s.Attr("name")
		name = strings.TrimSpace(name)
		if exist && len(name) > 0 {
			content, exist := s.Attr("content")
			if exist {
				metadata[strings.ToLower(name)] = content
			}
		}

		//check meta refresh
		equiv, exist := s.Attr("http-equiv")
		equiv = strings.TrimSpace(strings.ToLower(equiv))
		if exist && len(equiv) > 0 && equiv == "refresh" {
			content, exist := s.Attr("content")
			if exist {
				//0; url=/2016/beijing.html
				arr := strings.Split(content, "=")
				if len(arr) == 2 {
					url := arr[1]
					links[url] = "http-equiv-refresh"
				} else {
					log.Error("unexpected http-equiv", content)
				}

			}
		}

	})

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		href = strings.TrimSpace(href)
		if exist && len(href) > 0 && !(strings.HasPrefix(href, "javascript")) && !(strings.HasPrefix(href, "#")) && href != "/" {
			if strings.Contains(href, "#") {
				hrefs := strings.Split(href, "#")
				href = hrefs[0]
			}
			text := strings.TrimSpace(s.Text())
			text = strings.Replace(text, "\t", "", -1)

			if len(text) > 0 {
				log.Trace("get link: ", text, " , ", href)
				links[href] = text
			}
		}

	})

	if len(links) > 0 {
		context.Set(CONTEXT_PAGE_LINKS, links)
		snapshot.Links = model.LinkGroup{
			Internal: []model.PageLink{},
			External: []model.PageLink{},
		}

		for link, label := range links {
			host := util.GetHost(link)
			l := model.PageLink{
				Label: util.XSSHandle(label),
				Url:   link,
			}
			if host != "" && host != task.Host {
				snapshot.Links.External = append(snapshot.Links.External, l)
			} else {
				snapshot.Links.Internal = append(snapshot.Links.Internal, l)
			}
		}
	}

	snapshot.H1 = parseTag(doc, "h1")
	snapshot.H2 = parseTag(doc, "h2")
	snapshot.H3 = parseTag(doc, "h3")
	snapshot.H4 = parseTag(doc, "h4")
	snapshot.Bold = parseTag(doc, "b")
	snapshot.Italic = parseTag(doc, "i")

	images := map[string]string{}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exist := s.Attr("src")
		src = strings.TrimSpace(src)
		if exist {
			alt, _ := s.Attr("alt")
			alt = strings.TrimSpace(alt)
			images[src] = util.XSSHandle(alt)
		}
	})

	snapshot.Images = model.LinkGroup{
		Internal: []model.PageLink{},
		External: []model.PageLink{},
	}
	for link, label := range images {
		host := util.GetHost(link)
		l := model.PageLink{
			Label: util.XSSHandle(label),
			Url:   link,
		}
		if host != "" && host != task.Host {
			snapshot.Images.External = append(snapshot.Images.External, l)
		} else {
			snapshot.Images.Internal = append(snapshot.Images.Internal, l)
		}
	}

	//if reach max depth, skip for future fetch
	if depth > this.GetIntOrDefault(maxDepth, 10) {
		log.Trace("skip while reach max depth, ", depth, ", ", refUrl)
		context.Break(fmt.Sprintf("skip while reach max depth: %v", depth))
		return nil
	}
	//if reach max breadth, skip for future fetch
	if breadth > this.GetIntOrDefault(maxBreadth, 10) {
		log.Trace("skip while reach max breadth, ", breadth, ", ", refUrl)
		context.Break(fmt.Sprintf("skip while reach max breadth: %v", breadth))
		return nil
	}

	//dispatch links
	for url := range links {
		if this.GetBool(dispatchLinks, false) {
			if !filter.Exists(config.CheckFilter, []byte(url)) {
				host := util.GetHost(url)
				b := breadth
				if host != "" && refHost != host {
					b++
					log.Trace("auto incre breadth, ", b, ", ", refUrl, "->", url)
				}
				queue.Push(config.CheckChannel, model.NewTaskSeed(url, refUrl, depth+1, b).MustGetBytes())
			}
		}
	}

	return nil
}

func parseTag(doc *goquery.Document, tag string) []string {
	result := []string{}
	doc.Find(tag).Each(func(i int, s *goquery.Selection) {
		if len(s.Text()) > 0 {
			t := s.Text()
			t = util.NoWordBreak(t)
			t = strings.TrimSpace(t)
			if len(t) > 0 {
				result = append(result, util.XSSHandle(t))
			}
		}

	})
	return result
}
