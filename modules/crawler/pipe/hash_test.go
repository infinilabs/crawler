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
	"github.com/infinitbyte/gopa/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessHash(t *testing.T) {
	body := "Just some test content,你好"

	context := model.Context{}
	context.Init()
	task := model.Task{}
	task.Url = "http://elasticsearch.cn/"
	task.Depth = 1

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	parse := HashJoint{}

	snapshot := model.Snapshot{}
	snapshot.Payload = []byte(body)
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(&context)

	hash := snapshot.Hash
	assert.Equal(t, "b96aa2d91a6b69648250c8d4938b19c0750d7c50", hash)

	hash1 := task.SnapshotSimHash
	//assert.Equal(t, "13442536247490772857", hash1)

	body = "Just some test content,你好啊,!!"

	task1 := model.Task{}
	task.Url = "http://elasticsearch.cn/"
	task.Depth = 1

	context = model.Context{}
	context.Init()
	context.Set(CONTEXT_CRAWLER_TASK, &task1)
	parse = HashJoint{}

	snapshot = model.Snapshot{}
	snapshot.Payload = []byte(body)
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(&context)
	hash2 := snapshot.SimHash
	fmt.Println(hash1)
	fmt.Println(hash2)
	assert.Equal(t, hash2, hash1)
}
