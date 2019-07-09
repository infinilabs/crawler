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
	"crypto/sha1"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/gopa/model"
)

type HashJoint struct {
	pipeline.Parameters
}

func (joint HashJoint) Name() string {
	return "hash"
}

func (joint HashJoint) Process(context *pipeline.Context) error {

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	if snapshot.Payload == nil || len(snapshot.Payload) == 0 {
		log.Trace("snapshot payload is empty, skip hash,", snapshot.Payload)
		return nil
	}

	if global.Env().IsDebug {
		log.Trace("cal hash,", snapshot.Payload)
	}
	h := sha1.New()
	h.Write(snapshot.Payload)
	bs := h.Sum(nil)

	snapshot.Hash = fmt.Sprintf("%x", bs)

	if global.Env().IsDebug {
		log.Trace("get hash,", snapshot.Hash)
	}

	return nil
}
