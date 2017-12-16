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

package joint

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/filter"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	"strings"
)

type ParsePageJoint struct {
	model.Parameters
	MaxPageOfBreadth map[int]int //max page to fetch in each level's breadth, eg: 1:100;2:50;3:5;4:1
	//TODO support save link,script
}

const dispatchLinks model.ParaKey = "dispatch_links"
const saveImages model.ParaKey = "save_images"
const maxDepth model.ParaKey = "max_depth"
const maxBreadth model.ParaKey = "max_breadth"
const replaceNoscript model.ParaKey = "replace_noscript"

func (joint ParsePageJoint) Name() string {
	return "parse"
}

func (joint ParsePageJoint) Process(context *model.Context) error {

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	if !util.PrefixStr(snapshot.ContentType, "text/html") {
		log.Debugf("snapshot is not html, %s, %s , %s", snapshot.ID, snapshot.Url, snapshot.ContentType)
		return nil
	}

	refUrl := context.MustGetString(model.CONTEXT_TASK_URL)
	refHost := context.MustGetString(model.CONTEXT_TASK_Host)
	depth := context.MustGetInt(model.CONTEXT_TASK_Depth)
	breadth := context.MustGetInt(model.CONTEXT_TASK_Breadth)
	fileContent := snapshot.Payload

	//replace noscript to div
	if context.GetBool(replaceNoscript, true) {
		fileContent = util.ReplaceByte(snapshot.Payload, []byte("noscript"), []byte("div     "))
	}

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
				arr := strings.Split(strings.ToLower(content), "url=")
				if len(arr) == 2 {
					url := arr[1]
					links[url] = "http-equiv-refresh"
				} else {
					log.Error("unexpected http-equiv, ", content)
				}
				context.Set(model.CONTEXT_TASK_Status, model.TaskRedirected)
				context.End(fmt.Sprintf("redirected to: %v", content))
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
		context.Set(model.CONTEXT_PAGE_LINKS, links)
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
			if host != "" && host != refHost {
				snapshot.Links.External = append(snapshot.Links.External, l)
			} else {
				snapshot.Links.Internal = append(snapshot.Links.Internal, l)
			}
		}
	}

	snapshot.H1 = parseTagText(doc, "h1")
	snapshot.H2 = parseTagText(doc, "h2")
	snapshot.H3 = parseTagText(doc, "h3")
	snapshot.H4 = parseTagText(doc, "h4")
	snapshot.Bold = parseTagText(doc, "b")
	snapshot.Italic = parseTagText(doc, "i")

	images := map[string]string{}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exist := s.Attr("src")
		src = strings.TrimSpace(src)
		if exist {
			alt, _ := s.Attr("alt")
			alt = strings.TrimSpace(alt)
			images[src] = util.XSSHandle(alt)
		}

		if joint.GetBool(saveImages, false) {
			// save images
			context := model.Context{IgnoreBroken: true}
			context.Set(model.CONTEXT_TASK_URL, src)
			context.Set(model.CONTEXT_TASK_Reference, refUrl)
			context.Set(model.CONTEXT_TASK_Depth, 0)
			context.Set(model.CONTEXT_TASK_Breadth, 0)
			queue.Push(config.CheckChannel, util.ToJSONBytes(context))
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
		if host != "" && host != refHost {
			snapshot.Images.External = append(snapshot.Images.External, l)
		} else {
			snapshot.Images.Internal = append(snapshot.Images.Internal, l)
		}
	}

	log.Trace("depth:", depth, ", breath:", breadth, ",", joint.GetIntOrDefault(maxDepth, 10), ",", joint.GetIntOrDefault(maxBreadth, 10), ",url:", refUrl)

	//if reach max depth, skip for future fetch
	if depth > joint.GetIntOrDefault(maxDepth, 10) {
		log.Debug("skip while reach max depth, ", depth, ", ", refUrl)
		context.End(fmt.Sprintf("skip while reach max depth: %v", depth))
		return nil
	}
	//if reach max breadth, skip for future fetch
	if breadth > joint.GetIntOrDefault(maxBreadth, 10) {
		log.Debug("skip while reach max breadth, ", breadth, ", ", refUrl)
		context.End(fmt.Sprintf("skip while reach max breadth: %v", breadth))
		return nil
	}

	//dispatch links
	if joint.GetBool(dispatchLinks, true) && len(links) > 0 {

		for url := range links {
			if !filter.Exists(config.CheckFilter, []byte(url)) {
				host := util.GetHost(url)
				b := breadth
				d := depth
				if host != "" && refHost != host {
					b++
					d++
					log.Trace("auto incre breadth, ", b, ", ", refUrl, "->", url)
				} else {
					d++
				}

				context := model.Context{IgnoreBroken: true}
				context.Set(model.CONTEXT_TASK_URL, url)
				context.Set(model.CONTEXT_TASK_Reference, refUrl)
				context.Set(model.CONTEXT_TASK_Depth, d)
				context.Set(model.CONTEXT_TASK_Breadth, b)
				queue.Push(config.CheckChannel, util.ToJSONBytes(context))
			}

		}
	}

	return nil
}

func parseTagText(doc *goquery.Document, tag string) []string {
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
