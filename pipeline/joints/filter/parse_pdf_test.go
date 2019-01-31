/*
Copyright Medcl (m AT medcl.net)

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
	"github.com/golang/go/src/pkg/io/ioutil"
	"github.com/infinitbyte/framework/core/env"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessPDF1(t *testing.T) {
	global.RegisterEnv(env.EmptyEnv())
	fileName := "../../../test/samples/sample.pdf"
	text := ParsePDF2Text(fileName)
	fmt.Println(text)
	assert.Equal(t, true, util.ContainStr(text, "A Simple PDF File"), "pdf parse failure")
	assert.Equal(t, true, util.ContainStr(text, "More, a little more text."), "pdf parse failure")

}

func TestProcessPDF2(t *testing.T) {
	global.RegisterEnv(env.EmptyEnv())

	snapshot := model.Snapshot{}
	fileName := "../../../test/samples/sample.pdf"
	snapshot.Payload, _ = ioutil.ReadFile(fileName)
	snapshot.ContentType = "application/pdf"
	context := pipeline.Context{}
	context.Set(model.CONTEXT_TASK_Depth, 0)
	context.Set(model.CONTEXT_TASK_Breadth, 0)
	context.Set(model.CONTEXT_TASK_URL, "http://localhost:8000/sample.pdf")
	context.Set(model.CONTEXT_TASK_Host, "elasticsearch.cn/")

	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse := ParsePDFJoint{}

	parse.Process(&context)

	assert.Equal(t, true, util.ContainStr(snapshot.Text, "A Simple PDF File"), "pdf parse failure")
	assert.Equal(t, true, util.ContainStr(snapshot.Text, "More, a little more text."), "pdf parse failure")

	fmt.Println(snapshot.Text)
}
