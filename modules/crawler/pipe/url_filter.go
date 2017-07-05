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
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"regexp"
	"strings"

	"github.com/infinitbyte/gopa/core/model"
)

const UrlExtFilter JointKey = "url_ext_filter"

type UrlExtFilterJoint struct {
	//ignore files end with js,css,apk,zip
	SkipPageParsePattern *regexp.Regexp
}

func (this UrlExtFilterJoint) Name() string {
	return string(UrlExtFilter)
}

func (this UrlExtFilterJoint) Process(context *Context) error {
	this.SkipPageParsePattern = regexp.MustCompile(".*?\\.((js)|(css)|(rar)|(gz)|(zip)|(exe)|(bmp)|(jpeg)|(gif)|(png)|(jpg)|(apk))\\b")

	task := context.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)

	if task.OriginalUrl != "" && !this.valid(task.OriginalUrl) {
		context.ErrorExit("invalid url ext, " + task.OriginalUrl)
	}

	if task.Url != "" && !this.valid(task.Url) {
		context.ErrorExit("invalid url ext, " + task.Url)
	}

	return nil
}

func (this UrlExtFilterJoint) valid(url string) bool {
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
