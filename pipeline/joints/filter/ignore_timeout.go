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
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/stats"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
)

func (joint IgnoreTimeoutJoint) Name() string {
	return "ignore_timeout"
}

const ignoreTimeoutAfterCount pipeline.ParaKey = "ignore_timeout_after_count"

type IgnoreTimeoutJoint struct {
	pipeline.Parameters
}

func (joint IgnoreTimeoutJoint) Process(context *pipeline.Context) error {

	//task := context.MustGet(model.CONTEXT_TASK).(*model.Task)

	//TODO ignore within time period, rather than total count
	host := context.MustGetString(model.CONTEXT_TASK_Host)
	timeoutCount := stats.Stat("host.stats", host+"."+config.STATS_FETCH_TIMEOUT_COUNT)
	if timeoutCount > joint.MustGetInt64(ignoreTimeoutAfterCount) {
		stats.Increment("host.stats", host+"."+config.STATS_FETCH_TIMEOUT_IGNORE_COUNT)
		context.End("too much timeout on this host, ignored " + host)
		log.Warnf("hit timeout host, %s , ignore after,%d ", host, timeoutCount)
	}
	return nil
}
