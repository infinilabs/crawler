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
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/util"
	"path/filepath"
	"strings"
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
	sysConfig := LoadSystemConfig(configFile)
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
	system:=GetDefaultSystemConfig()
	return &Env{SystemConfig: &system, RuntimeConfig: &RuntimeConfig{}}
}
