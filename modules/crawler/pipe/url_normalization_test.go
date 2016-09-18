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
	"testing"
	"github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/pipeline"
)

func TestNormailzeLinks(t *testing.T) {

	context:= pipeline.Context{Env:env.EmptyEnv()}
	context.Data=map[string]interface{}{}
	context.Set(CONTEXT_URL.String(),[]byte("http://elasticsearch.cn/"))
	context.Set(CONTEXT_DEPTH.String(),1)
	parse:=UrlNormalizationJoint{}
	parse.Process(&context)


	//assert.Equal(t,"baidu",links["baidu.com"])
	//assert.Equal(t,"/wiki/Marking/Users",links["http://elasticsearch.cn/wiki/Marking/Users"])

}
