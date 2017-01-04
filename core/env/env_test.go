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

package env

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestLoad(t *testing.T) {
	env:=Environment(SystemConfig{ConfigFile:"../../gopa.yml",LogLevel:"debug"})
	fmt.Sprintln(env)
	fmt.Sprintln(env.SystemConfig)
	config:=env.RuntimeConfig
	assert.Equal(t,"gopa",config.ClusterConfig.Name)
	assert.Equal(t,"debug",config.LoggingConfig.Level)
	assert.Equal(t,"http://eshost:9200",config.IndexingConfig.Host)
	assert.Equal(t,"gopa",config.IndexingConfig.Index)
	//assert.Equal(t,"data",config.PathConfig.Data)
	//assert.Equal(t,"log",config.PathConfig.Log)

	assert.Equal(t,true,config.ParserConfig.Enabled)
	assert.Equal(t,true,config.CrawlerConfig.Enabled)

}
