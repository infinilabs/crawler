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
	. "github.com/infinitbyte/gopa/core/env"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test(t *testing.T) {
	var filter FilterModule
	env1 := EmptyEnv()
	env1.SystemConfig.PathConfig.Data = "/tmp/filter_" + util.PickRandomName()
	os.RemoveAll(env1.SystemConfig.PathConfig.Data)
	env1.IsDebug = true
	global.RegisterEnv(env1)

	filter = FilterModule{}
	filter.Start(GetModuleConfig(filter.Name()))
	b, _ := filter.CheckThenAdd(config.CheckFilter, []byte("key"))
	assert.Equal(t, false, b)

	for i := 0; i < 1000; i++ {
		go run(&filter, t)
	}

}

func run(filter *FilterModule, t *testing.T) {
	for i := 0; i < 100000; i++ {
		b, _ := filter.CheckThenAdd(config.CheckFilter, []byte("key"))
		assert.Equal(t, true, b)
		b = filter.Exists(config.CheckFilter, []byte("key"))
		assert.Equal(t, true, b)
		if !b {
			fmt.Print("not exists")
		}
	}
}
