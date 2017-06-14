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
	"github.com/PuerkitoBio/goquery"
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/filter"
	"github.com/medcl/gopa/core/model"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/util"
	"github.com/medcl/gopa/modules/config"
	"strings"
)

const ParsePage JointKey = "parse"

type ParsePageJoint struct {
	DispatchLinks    bool
	MaxDepth         int         //max depth of page to follow
	MaxBreadth       int         //max breadth of the domain to follow
	MaxPageOfBreadth map[int]int //max page to fetch in each level's breadth, eg: 1:100;2:50;3:5;4:1
	//TODO support save link,script
}

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
		snapshot.Title = title
	}

	links := map[string]interface{}{}

	metadata := map[string]interface{}{}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, exist := s.Attr("name")
		name = strings.TrimSpace(name)
		if exist && len(name) > 0 {
			content, exist := s.Attr("content")
			if exist {
				metadata[name] = content
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

	if len(metadata) > 0 {
		snapshot.Metadata = &metadata
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		href = strings.TrimSpace(href)
		if exist && len(href) > 0 && !(strings.HasPrefix(href, "javascript")) && !(strings.HasPrefix(href, "#")) && href != "/" {
			if strings.Contains(href, "#") {
				hrefs := strings.Split(href, "#")
				href = hrefs[0]
			}
			text := strings.TrimSpace(s.Text())
			strings.Replace(text, "\t", "", -1)

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

		//TODO parse internal and external links
		//snapshot.Links.External=this.links
		//snapshot.Links.Internal=this.links
	}

	//if reach max depth, skip for future fetch
	if depth > this.MaxDepth {
		log.Trace("skip while reach max depth, ", depth, ", ", refUrl)
		context.Break("skip while reach max depth")
		return nil
	}
	//if reach max breadth, skip for future fetch
	if breadth > this.MaxBreadth {
		log.Trace("skip while reach max breadth, ", breadth, ", ", refUrl)
		context.Break("skip while reach max breadth")
		return nil
	}

	//dispatch links
	for url := range links {
		if this.DispatchLinks {
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
