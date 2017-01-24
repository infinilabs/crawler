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
	"github.com/medcl/gopa/core/pipeline"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlFilter(t *testing.T) {

	context := &pipeline.Context{}
	context.Init()
	context.Set(CONTEXT_URL, "http://elasticsearch.cn/")
	context.Set(CONTEXT_ORIGINAL_URL, "http://elasticsearch.cn/")
	parse := UrlExtFilterJoint{}
	parse.Process(context)
	assert.Equal(t, false, context.IsBreak())

	context.Set(CONTEXT_URL, "mailto:gg@gg.com")
	parse.Process(context)
	assert.Equal(t, true, context.IsErrorExit())

	context.Set(CONTEXT_URL, "asfasdf.gif")
	parse.Process(context)
	assert.Equal(t, true, context.IsErrorExit())

}
