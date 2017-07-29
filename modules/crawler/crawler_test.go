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
	. "github.com/infinitbyte/gopa/core/env"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/util"
	db "github.com/infinitbyte/gopa/modules/database"
	f "github.com/infinitbyte/gopa/modules/filter"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test1(t *testing.T) {
	env1 := EmptyEnv()
	env1.SystemConfig.PathConfig.Data = "/tmp/filter_" + util.RandomPickName()
	os.RemoveAll(env1.SystemConfig.PathConfig.Data)
	env1.IsDebug = true
	global.RegisterEnv(env1)

	checker := CheckerModule{}

	filter := f.FilterModule{}
	filter.Start(GetModuleConfig(filter.Name()))

	db.DatabaseModule{}.Start(GetModuleConfig("database"))

	task := model.Task{}
	task.Url = "http://elasticsearch.cn"

	pipeline := checker.runPipe(true, &task)
	fmt.Println(util.ToJson(pipeline.GetContext(), true))
	assert.Equal(t, false, pipeline.GetContext().IsErrorExit())

	for i := 0; i < 10; i++ {
		pipeline := checker.runPipe(true, &task)
		assert.Equal(t, true, pipeline.GetContext().IsErrorExit())
		if !pipeline.GetContext().IsErrorExit() {
			fmt.Print("not exists")
		}
	}
}
