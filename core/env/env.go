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
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	. "github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/util"
	"os"
	"path/filepath"
	"strings"
)

const VERSION = "0.9.0_SNAPSHOT"

var (
	LastCommitLog string
	BuildDate string
)


type Env struct {

	// static configs
	SystemConfig *SystemConfig

	// dynamic configs
	RuntimeConfig *RuntimeConfig

	IsDebug bool

	LoggingLevel string
}

func Environment(configFile string) *Env {

	env := Env{}
	sysConfig := loadSystemConfig(configFile)
	env.SystemConfig = &sysConfig

	var err error
	env.RuntimeConfig, err = env.loadRuntimeConfig()
	if err != nil {
		log.Error(err)
		panic(err)
	}

	return &env
}

var moduleConfig map[string]*Config

var (
	defaultSystemConfig = SystemConfig{
		ClusterConfig: ClusterConfig{
			Name: "gopa",
		},
		NetworkConfig: NetworkConfig{
			Host: "127.0.0.1",
		},
		NodeConfig: NodeConfig{
			Name: util.RandomPickName(),
		},
		PathConfig: PathConfig{
			Data: "data",
			Log:  "log",
			Cert: "cert",
		},

		APIBinding:         ":8001",
		HttpBinding:        ":9001",
		ClusterBinding:     ":13001",
		AllowMultiInstance: false,
	}
)

func loadSystemConfig(cfgFile string) SystemConfig {
	cfg := defaultSystemConfig
	cfg.ConfigFile = cfgFile
	if util.IsExist(cfgFile) {
		config, err := yaml.NewConfigWithFile(cfgFile, ucfg.PathSep("."))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = config.Unpack(&cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	os.MkdirAll(cfg.GetDataDir(), 0777)
	os.MkdirAll(cfg.PathConfig.Log, 0777)
	return cfg
}

var (
	defaultRuntimeConfig = RuntimeConfig{}
)

func (this *Env) loadRuntimeConfig() (*RuntimeConfig, error) {

	moduleConfig = map[string]*Config{}

	var configFile = "./gopa.yml"
	if this.SystemConfig != nil && len(this.SystemConfig.ConfigFile) > 0 {
		configFile = this.SystemConfig.ConfigFile
	}

	filename, _ := filepath.Abs(configFile)
	var cfg RuntimeConfig

	if util.FileExists(filename) {
		log.Debug("load configFile:", filename)
		cfg, err := LoadFile(filename)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		config := defaultRuntimeConfig

		if err := cfg.Unpack(&config); err != nil {
			log.Error(err)
			return nil, err
		}

		parseModuleConfig(config.Modules)

	} else {
		log.Debug("no config file was found")

		cfg = defaultRuntimeConfig
	}

	return &cfg, nil
}

func parseModuleConfig(cfgs []*Config) []*Config {
	results := []*Config{}
	for _, cfg := range cfgs {
		//set map for modules and module config
		log.Trace(getModuleName(cfg), ",", cfg.Enabled())
		config, err := NewConfigFrom(cfg)
		if err != nil {
			panic(err)
		}

		name := getModuleName(cfg)
		moduleConfig[name] = cfg

		results = append(results, config)
	}

	return results
}

func GetModuleConfig(name string) *Config {
	cfg := moduleConfig[strings.ToLower(name)]
	return cfg
}

func getModuleName(c *Config) string {
	cfgObj := struct {
		Module string `config:"module"`
	}{}

	if c == nil {
		return ""
	}
	if err := c.Unpack(&cfgObj); err != nil {
		return ""
	}

	return cfgObj.Module
}

func EmptyEnv() *Env {
	system := defaultSystemConfig
	return &Env{SystemConfig: &system, RuntimeConfig: &RuntimeConfig{}}
}
