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
	"github.com/medcl/gopa/core/model"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/modules/config"
	"strings"
)

const ParsePage JointKey = "parse"

type ParsePageJoint struct {
	links         map[string]interface{}
	DispatchLinks bool
	MaxDepth      int
}

func (this ParsePageJoint) Name() string {
	return string(ParsePage)
}

func (this ParsePageJoint) Process(s *Context) (*Context, error) {

	refUrl := s.MustGetString(CONTEXT_URL)
	depth := s.MustGetInt(CONTEXT_DEPTH)
	fileContent := s.MustGetBytes(CONTEXT_PAGE_BODY_BYTES)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(fileContent))
	if err != nil {
		panic(err)
	}

	title := doc.Find("title").Text()

	selected := doc.Find("body")
	body, err := selected.Html()
	if err != nil {
		panic(err)
	} else {
		s.Set(CONTEXT_PAGE_BODY_BYTES, []byte(body))
	}

	metadata := map[string]interface{}{}
	if len(title) > 0 {
		metadata["title"] = title
	}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, exist := s.Attr("name")
		name = strings.TrimSpace(name)
		if exist && len(name) > 0 {
			content, exist := s.Attr("content")
			if exist {
				metadata[name] = content
			}
		}

	})

	if len(metadata) > 0 {
		s.Set(CONTEXT_PAGE_METADATA, metadata)
	}

	this.links = map[string]interface{}{}
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
				this.links[href] = text
			}
		}

	})

	s.Set(CONTEXT_PAGE_LINKS, this.links)

	//if reach max depth, skip for future fetch
	if depth >= this.MaxDepth {
		return s, nil
	}

	//dispatch links
	for url := range this.links {
		if this.DispatchLinks {
			queue.Push(config.CheckChannel, model.NewTaskSeed(url, refUrl, depth+1).MustGetBytes())
		}
	}

	return s, nil
}
