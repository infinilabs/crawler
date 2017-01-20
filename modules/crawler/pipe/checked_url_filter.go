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
	"github.com/medcl/gopa/core/filter"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/modules/config"
	"regexp"
)

const UrlCheckedFilter JointKey = "url_checked_filter"

type UrlCheckedFilterJoint struct {
	Parameters
	//ignore files end with js,css,apk,zip
	SkipPageParsePattern *regexp.Regexp
}

func (this UrlCheckedFilterJoint) Name() string {
	return string(UrlCheckedFilter)
}

func (this UrlCheckedFilterJoint) Process(context *Context) (*Context, error) {
	url := context.MustGetString(CONTEXT_URL)
	//统一 url 格式 , url 此处应该不能是相对路径

	log.Trace("cheking url:", url)

	b, err := filter.CheckThenAdd(config.CheckFilter, []byte(url))
	//checking
	if b {
		stats.Increment("checker.url", "duplicated")
		log.Debug("duplicated url,already checked,  url:", url)
		context.Exit("duplicated url,already checked,  url:" + url)
		return context, nil
	}
	if err != nil {
		log.Error(err)
		panic(err)
		context.Break("check url error, url: " + url + ", " + err.Error())
	}

	return context, nil
}
