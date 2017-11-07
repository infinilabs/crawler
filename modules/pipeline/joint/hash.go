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

package joint

import (
	"crypto/sha1"
	"fmt"
	"github.com/infinitbyte/gopa/core/model"
)

type HashJoint struct {
	model.Parameters
}

func (joint HashJoint) Name() string {
	return "hash"
}

func (joint HashJoint) Process(context *model.Context) error {

	snapshot := context.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	h := sha1.New()
	h.Write(snapshot.Payload)
	bs := h.Sum(nil)

	snapshot.Hash = fmt.Sprintf("%x", bs)

	return nil
}
