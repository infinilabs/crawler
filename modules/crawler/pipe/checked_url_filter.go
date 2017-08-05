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
	"github.com/infinitbyte/gopa/core/filter"
	"github.com/infinitbyte/gopa/core/model"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/modules/config"
	"regexp"
)

const UrlCheckFilter JointKey = "url_check_filter"

type UrlCheckFilterJoint struct {
	Parameters
	//ignore files end with js,css,apk,zip
	SkipPageParsePattern *regexp.Regexp
}

func (joint UrlCheckFilterJoint) Name() string {
	return string(UrlCheckFilter)
}

func (joint UrlCheckFilterJoint) Process(context *Context) error {

	task := context.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)

	url := task.Url

	//the url input here should not be a relative path
	b, err := filter.CheckThenAdd(config.CheckFilter, []byte(url))
	log.Trace("cheking url:", url, ",hit:", b)

	//checking
	if b {
		stats.Increment("checker.url", "duplicated")
		log.Trace("duplicated url,already checked,  url:", url)
		context.ErrorExit("duplicated url,already checked,  url:" + url)
		return nil
	}
	if err != nil {
		log.Error(err)
		context.Break("check url error, url: " + url + ", " + err.Error())
	}

	return nil
}
