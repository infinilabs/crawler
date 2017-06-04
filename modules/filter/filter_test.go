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
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/util"
	"github.com/medcl/gopa/modules/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test(t *testing.T) {

	env1 := EmptyEnv()
	env1.SystemConfig.PathConfig.Data = "/tmp/filter_" + util.RandomPickName()
	os.RemoveAll(env1.SystemConfig.PathConfig.Data)
	env1.IsDebug = true
	global.RegisterEnv(env1)

	filter := FilterModule{}
	filter.Start(GetModuleConfig(filter.Name()))
	b, _ := filter.CheckThenAdd(config.CheckFilter, []byte("key"))
	assert.Equal(t, false, b)
	for i := 0; i < 1000; i++ {
		b, _ := filter.CheckThenAdd(config.CheckFilter, []byte("key"))
		assert.Equal(t, true, b)
		b = filter.Exists(config.CheckFilter, []byte("key"))
		assert.Equal(t, true, b)
		if !b {
			fmt.Print("not exists")
		}
	}

}
