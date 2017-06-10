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
	"fmt"
	"github.com/medcl/gopa/core/pipeline"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessHash(t *testing.T) {
	body := "Just some test content,你好"

	context := pipeline.Context{}
	context.Init()
	context.Set(CONTEXT_URL, "http://elasticsearch.cn/")
	context.Set(CONTEXT_DEPTH, 1)
	context.Set(CONTEXT_PAGE_BODY_PLAIN_TEXT, body)
	parse := HashJoint{DictRoot: "../../../"}
	parse.Process(&context)

	hash := context.MustGetString(CONTEXT_PAGE_HASH)
	assert.Equal(t, "b96aa2d91a6b69648250c8d4938b19c0750d7c50", hash)

	hash1 := context.MustGetString(CONTEXT_PAGE_SIMHASH_100)
	//assert.Equal(t, "13442536247490772857", hash1)

	body = "Just some test content,你好啊,!!"

	context = pipeline.Context{}
	context.Init()
	context.Set(CONTEXT_URL, "http://elasticsearch.cn/")
	context.Set(CONTEXT_DEPTH, 1)
	context.Set(CONTEXT_PAGE_BODY_PLAIN_TEXT, body)
	parse = HashJoint{DictRoot: "../../../"}
	parse.Process(&context)
	hash2 := context.MustGetString(CONTEXT_PAGE_SIMHASH_100)
	fmt.Println(hash1)
	fmt.Println(hash2)
	assert.Equal(t, hash2, hash1)
}
