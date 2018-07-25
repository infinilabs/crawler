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
	"fmt"
	"github.com/infinitbyte/framework/core/env"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/gopa/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestExtractBlock(t *testing.T) {
	global.RegisterEnv(env.EmptyEnv())

	body, e := ioutil.ReadFile("../../../test/samples/default.html")
	if e != nil {
		panic(e)
	}

	context := pipeline.Context{}
	context.Set(model.CONTEXT_TASK_Depth, 0)
	context.Set(model.CONTEXT_TASK_Breadth, 0)
	context.Set(model.CONTEXT_TASK_URL, "http://elasticsearch.cn/")
	context.Set(model.CONTEXT_TASK_Host, "elasticsearch.cn/")
	parse := ExtractJoint{}
	m := map[string]interface{}{}
	m["h1"] = "h1"
	m["h2"] = "h2"
	m["pre"] = "pre"
	parse.Set(htmlBlock, m)
	snapshot := model.Snapshot{}
	snapshot.ContentType = "text/html"

	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)
	snapshot.Payload = []byte(body)
	parse.Process(&context)
	fmt.Println(snapshot.EnrichedFeatures)
	assert.Equal(t, "H1 title", (*snapshot.EnrichedFeatures)["h1"].(string))
	assert.Equal(t, "H2 title", (*snapshot.EnrichedFeatures)["h2"].(string))
	assert.Equal(t, "code block", (*snapshot.EnrichedFeatures)["pre"].(string))

}
