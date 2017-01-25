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

package crawler

import (
	"fmt"
	"github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/core/util"
	db "github.com/medcl/gopa/modules/database"
	f "github.com/medcl/gopa/modules/filter"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test1(t *testing.T) {
	env := env.EmptyEnv()
	env.SystemConfig.Data = "/tmp/filter_" + util.RandomPickName()
	os.RemoveAll(env.SystemConfig.Data)
	env.IsDebug = true
	global.RegisterEnv(env)

	checker := CheckerModule{}

	filter := f.FilterModule{}
	filter.Start(env)

	db.DatabaseModule{}.Start(env)

	task := model.Task{}
	task.Url = "http://elasticsearch.cn"

	pipeline := checker.runPipe(true, &task)
	fmt.Println(util.ToJson(pipeline.GetContext(), true))
	assert.Equal(t, false, pipeline.GetContext().IsErrorExit())

	for i := 0; i < 10000; i++ {
		pipeline := checker.runPipe(true, &task)
		assert.Equal(t, true, pipeline.GetContext().IsErrorExit())
		if !pipeline.GetContext().IsErrorExit() {
			fmt.Print("not exists")
		}
	}
}
