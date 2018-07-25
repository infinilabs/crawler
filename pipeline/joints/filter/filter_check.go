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

package filter

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/filter"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/stats"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
	"regexp"
)

// FilterCheckJointused to check the task url if it is already in the filter, if not in the filter, then add it to task filter, and make sure won't add it next time
type FilterCheckJoint struct {
	pipeline.Parameters
	//ignore files end with js,css,apk,zip
	SkipPageParsePattern *regexp.Regexp
}

// filter_key is the filter name used to check and filter
var filterKey pipeline.ParaKey = "filter_key"

// Name return: filter_check
func (joint FilterCheckJoint) Name() string {
	return "filter_check"
}

// Process the filtering and add it to the filter
func (joint FilterCheckJoint) Process(context *pipeline.Context) error {

	url := context.MustGetString(model.CONTEXT_TASK_URL)

	//key := joint.GetStringOrDefault(filterKey, "check_filter")
	//v := filter.Key(key)

	hash := util.Sha1Hash(url)

	//the url input here should not be a relative path
	b, err := filter.CheckThenAdd(config.CheckFilter, []byte(hash))
	log.Trace("cheking url:", url, ",hit:", b)
	if err != nil {
		log.Error(err)
		context.End("check url error, url: " + url + ", " + err.Error())
	}
	//checking
	if b {
		stats.Increment("checker.url", "duplicated")
		log.Trace("duplicated url,already checked,  url:", url)
		context.Exit("duplicated url,already checked,  url:" + url)
		return nil
	}

	return nil
}
