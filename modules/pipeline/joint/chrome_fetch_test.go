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
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/util"
	"testing"
)

func ChromeFetch1(t *testing.T) {

	context := model.Context{}
	context.Set(model.CONTEXT_TASK_Host, "elasticsearch.cn")
	context.Set(model.CONTEXT_TASK_URL, "https://elasticsearch.cn/article/383")
	context.Set(model.CONTEXT_TASK_Depth, 0)
	parse := ChromeFetchJoint{}

	snapshot := model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(&context)
	fmt.Println(util.ToJson(context, true))

}
