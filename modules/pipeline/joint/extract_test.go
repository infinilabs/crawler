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
	"fmt"
	"github.com/infinitbyte/gopa/core/env"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/model"
	"io/ioutil"
	"testing"
)

func TestExtractBlock(t *testing.T) {
	global.RegisterEnv(env.EmptyEnv())

	body, e := ioutil.ReadFile("../../../test/samples/shiti.html")
	if e != nil {
		panic(e)
	}

	context := model.Context{}
	context.Init()
	context.Set(model.CONTEXT_TASK_Depth, 0)
	context.Set(model.CONTEXT_TASK_Breadth, 0)
	context.Set(model.CONTEXT_TASK_URL, "http://zujuan.21cnjy.com/question/detail/5492219")
	context.Set(model.CONTEXT_TASK_Host, "zujuan.21cnjy.com")
	parse := ExtractJoint{}
	m := map[string]interface{}{}
	m["question_header"] = ".exam-head"
	m["question_con"] = ".exam-con"
	m["question_brick"] = ".analyticbox-brick"
	parse.Set(htmlBlock, m)
	snapshot := model.Snapshot{}
	snapshot.ContentType = "text/html"

	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)
	snapshot.Payload = []byte(body)
	parse.Process(&context)
	fmt.Println(snapshot.EnrichedFeatures)

}
