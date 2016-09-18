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
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/stats"
)

//load metadata from db
type LoadMetadataJoint struct {
}

func (this LoadMetadataJoint) Process(context *Context) (*Context, error) {

	//url:=context.MustGetString(CONTEXT_URL)

	return context, nil
}



type IgnoreTimeoutJoint struct {
	IgnoreTimeoutAfterCount int
}

func (this IgnoreTimeoutJoint) Process(context *Context) (*Context, error) {


	host:=context.MustGetString(CONTEXT_HOST)
	timeoutCount:=stats.Stat(host,stats.STATS_FETCH_TIMEOUT_COUNT)
	if(timeoutCount>this.IgnoreTimeoutAfterCount){
		context.Break()
		log.Warnf("hit timeout host, %s , ignore after,%d ",host,timeoutCount)
	}
	return context, nil
}


