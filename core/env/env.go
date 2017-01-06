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
	"errors"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

type Env struct {
	SystemConfig  *SystemConfig
	RuntimeConfig *RuntimeConfig
	IsDebug       bool
	LoggingLevel  string
}

func Environment(sysConfig SystemConfig) *Env {

	env := Env{}

	env.SystemConfig = &sysConfig
	config, err := env.loadRuntimeConfig()
	if err != nil {
		panic(err)
	}
	env.RuntimeConfig = &config
	env.init()

	return &env
}

func (this *Env) loadRuntimeConfig() (RuntimeConfig, error) {

	var configFile = "./gopa.yml"
	if this.SystemConfig != nil && len(this.SystemConfig.ConfigFile) > 0 {
		configFile = this.SystemConfig.ConfigFile
	}

	//load external yaml config
	filename, _ := filepath.Abs(configFile)
	var config RuntimeConfig

	if util.FileExists(filename) {
		log.Debug("load configFile:", filename)

		yamlFile, err := ioutil.ReadFile(filename)

		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(err)
		}
	} else {
		//init default Config
		config = RuntimeConfig{}
		config.IndexingConfig = (&IndexingConfig{}).Init()
		config.ChannelConfig = (&ChannelConfig{}).Init()
		config.CrawlerConfig = (&CrawlerConfig{}).Init()
		config.ParserConfig = (&ParserConfig{}) //.Init()
		config.TaskConfig = (&TaskConfig{}).Init()
		config.RuledFetchConfig = (&RuledFetchConfig{}) //.Init()
	}

	config.TaskConfig.LinkUrlExtractRegex = regexp.MustCompile(config.TaskConfig.LinkUrlExtractRegexStr)
	config.TaskConfig.FetchUrlPattern = regexp.MustCompile(config.TaskConfig.FetchUrlPatternStr)
	config.TaskConfig.SavingUrlPattern = regexp.MustCompile(config.TaskConfig.SavingUrlPatternStr)
	config.TaskConfig.SkipPageParsePattern = regexp.MustCompile(config.TaskConfig.SkipPageParsePatternStr)

	return config, nil
}

func (this *Env) init() error {

	if this.RuntimeConfig.MaxGoRoutine < 1 {
		this.RuntimeConfig.MaxGoRoutine = 1
	}

	return nil
}

func EmptyEnv() *Env {
	return &Env{}
}

//high priority config, init from the environment or startup, can't be changed
type SystemConfig struct {
	ClusterName        string `cluster_name`
	NodeName           string `node_name`
	ConfigFile         string `gopa.yml`
	LogLevel           string `info`
	HttpBinding        string `http_bind`
	ClusterBinding     string `cluster_bind`
	ClusterSeeds       string `cluster_seeds`
	AllowMultiInstance bool   `multi_instance`
	Data               string `data`
	Log                string `log`
}

func (this *SystemConfig) Init() {
	if len(this.Data) == 0 {
		this.Data = "data"
	}
	if len(this.Log) == 0 {
		this.Log = "log"
	}
	if len(this.ClusterName) == 0 {
		this.ClusterName = "gopa"
	}
	if len(this.NodeName) == 0 {
		this.NodeName = util.RandomPickName()
	}
	if len(this.HttpBinding) == 0 {
		this.HttpBinding = ":8001"
	}

	if len(this.ClusterBinding) == 0 {
		this.ClusterBinding = ":13001"
	}

	this.AllowMultiInstance = false
	os.MkdirAll(this.GetDataDir(), 0777)
	os.MkdirAll(this.Log, 0777)
}

func (this *SystemConfig) GetDataDir() string {
	if this.AllowMultiInstance == false {
		return path.Join(this.Data, this.ClusterName, "nodes", "0")
	}
	//TODO auto select next nodes folder, eg: nodes/1 nodes/2
	panic(errors.New("not supported yet"))
}
