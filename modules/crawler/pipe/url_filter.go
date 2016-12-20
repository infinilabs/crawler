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
	. "github.com/medcl/gopa/core/pipeline"
	"strings"
	log "github.com/cihub/seelog"
	"regexp"
)

type UrlFilterJoint struct {
	//ignore files end with js,css,apk,zip
	SkipPageParsePattern *regexp.Regexp
}

func (this UrlFilterJoint) Name() string {
	return "url_filter"
}

func (this UrlFilterJoint) Process(context *Context) (*Context, error) {
	this.SkipPageParsePattern = regexp.MustCompile(".*?\\.((js)|(css)|(rar)|(gz)|(zip)|(exe)|(bmp)|(jpeg)|(gif)|(png)|(jpg)|(apk))\\b")
	url := context.MustGetString(CONTEXT_URL)
	orgUrl := context.MustGetString(CONTEXT_ORIGINAL_URL)

	if orgUrl == "" {
		orgUrl = url
	}

	if (!this.valid(orgUrl)) || (url != orgUrl && (!this.valid(url))) {
		context.Break("invalid url,"+url)
	}

	return context, nil
}

func (this UrlFilterJoint) valid(url string) bool {
	if url == "" {
		return false
	}

	if strings.HasPrefix(url, "mailto:") {
		log.Trace("filteredUrl started with: mailto: , invalid")
		return false
	}

	if strings.Contains(url, "data:image/") {
		log.Trace("filteredUrl started with: data:image/ , invalid")
		return false
	}

	if strings.HasPrefix(url, "#") {
		log.Trace("filteredUrl started with: # , invalid")
		return false
	}

	if strings.HasPrefix(url, "javascript:") {
		log.Trace("filteredUrl started with: javascript: , invalid")
		return false
	}


	if this.SkipPageParsePattern.Match([]byte(url)) {
		log.Trace("hit SkipPattern pattern,", url)
		return false
	}

	return true
}


//func checkIfUrlWillBeSave(taskConfig *config.TaskConfig, url []byte) bool {
//
//	requestUrl := string(url)
//
//	log.Debug("started check savingUrlPattern,", taskConfig.SavingUrlPattern, ",", string(url))
//	if taskConfig.SavingUrlPattern.Match(url) {
//
//		log.Debug("match saving url pattern,", requestUrl)
//		if len(taskConfig.SavingUrlMustNotContain) > 0 {
//			if util.ContainStr(requestUrl, taskConfig.SavingUrlMustNotContain) {
//				log.Debug("hit SavingUrlMustNotContain,ignore,", requestUrl, " , ", taskConfig.SavingUrlMustNotContain)
//				return false
//			}
//		}
//
//		if len(taskConfig.SavingUrlMustContain) > 0 {
//			if !util.ContainStr(requestUrl, taskConfig.SavingUrlMustContain) {
//				log.Debug("not hit SavingUrlMustContain,ignore,", requestUrl, " , ", taskConfig.SavingUrlMustContain)
//				return false
//			}
//		}
//
//		return true
//
//	} else {
//		log.Debug("does not hit SavingUrlPattern ignoring,", requestUrl)
//	}
//	return false
//}
