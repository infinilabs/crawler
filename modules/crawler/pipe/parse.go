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
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/types"
	"strings"
	"regexp"
	"github.com/medcl/gopa/core/util"
)

type ParserJoint struct {
	links         map[string]interface{}
	DispatchLinks bool
	MaxDepth int
}

func (this ParserJoint) Name() string {
	return "parse"
}


func (this ParserJoint) Process(s *Context) (*Context, error) {

	refUrl := s.MustGetString(CONTEXT_URL)
	depth := s.MustGetInt(CONTEXT_DEPTH)
	fileContent := s.Get(CONTEXT_PAGE_BODY_BYTES).([]byte)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(fileContent))
	if err != nil {
		panic(err)
	}

	title := doc.Find("title").Text()

	selected:=doc.Find("body")
	selected.RemoveFiltered("script")
	selected.RemoveFiltered("noscript")
	selected.RemoveFiltered("div[style*='display: none']")


	body,err:=selected.Html()
	if(err!=nil){
		panic(err)
	}
	src := body

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除META
	re, _ = regexp.Compile("\\<meta[\\S\\s]+?\\</meta\\>")
	src = re.ReplaceAllString(src, "")

	//去除注释
	re, _ = regexp.Compile("\\<!--[\\S\\s]+? --\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除NOSCRIPT
	re, _ = regexp.Compile("\\<noscript[\\S\\s]+?\\</noscript\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	src = strings.TrimSpace(util.MergeSpace(src))


	s.Set(CONTEXT_PAGE_BODY_PLAIN_TEXT,src)

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
	if(depth>=this.MaxDepth){
		return s,nil
	}

	//dispatch links
	for url, _ := range this.links {
		if this.DispatchLinks {
			s.Env.Channels.PushUrlToCheck(types.NewPageTask(url, refUrl, depth+1))
		}
	}

	return s, nil
}
