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
	"strings"
	"github.com/medcl/gopa/core/types"
)

type ParserJoint struct {
	links map[string]interface{}
	DispatchLinks bool
}

func (this ParserJoint) Process(s *Context) (*Context, error) {

	refUrl := s.MustGetString(CONTEXT_URL)
	depth := s.MustGetInt(CONTEXT_DEPTH)
	fileContent := s.Get(CONTEXT_PAGE_BODY_BYTES).([]byte)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(fileContent))
	if err != nil {
		panic(err)
	}



	title:=doc.Find("title").Text()

	metadata:=map[string]interface{}{}
	if(len(title)>0){
		metadata["title"]=title
	}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		name,exist := s.Attr("name")
		name=strings.TrimSpace(name)
		if(exist&&len(name)>0){
			content,exist := s.Attr("content")
			if(exist){
				metadata[name]=content
			}
		}

	})

	if(len(metadata)>0){
		s.Set(CONTEXT_PAGE_METADATA,metadata)
	}


	this.links = map[string]interface{}{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		href = strings.TrimSpace(href)
		if exist && len(href) > 0 && !(strings.HasPrefix(href,"javascript"))&& !(strings.HasPrefix(href,"#")) && href!="/" {
			if(strings.Contains(href,"#")){
				hrefs:=strings.Split(href,"#")
				href=hrefs[0]
			}
			text := strings.TrimSpace(s.Text())
			strings.Replace(text, "\t", "", -1)

			if len(text) > 0 {
				log.Debug("get link: ", text, " , ", href)
				this.links[href] = text
			}
		}

	})

	s.Set(CONTEXT_PAGE_LINKS, this.links)

	//dispatch links
	for url,_ := range this.links{
		if(this.DispatchLinks){
			s.Env.Channels.PushUrlToCheck(types.NewPageTask(url,refUrl,depth+1))
		}
	}




	return s, nil
}
