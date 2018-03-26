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

package module

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/env"
)

const System = "system"
const Service = "service"
const Tools = "tools"
const PipelineJoint = "joint"
const Stats = "stats"
const Filter = "filter"
const Database = "database"
const KVStore = "kv"
const Queue = "queue"
const Storage = "storage"
const Index = "index"
const Logger = "logger"
const UI = "ui"
const API = "api"

type Modules struct {
	modules []Module
	plugins []Module
	configs map[string]interface{}
}

var m *Modules

func New() {
	mod := Modules{}
	m = &mod
}

func Register(moduleType string, mod Module) {
	m.modules = append(m.modules, mod)
}

func RegisterPlugin(moduleType string, mod Module) {
	m.plugins = append(m.plugins, mod)
}

func Start() {

	log.Trace("start to load plugins")
	for _, v := range m.plugins {

		cfg := env.GetPluginConfig(v.Name())

		log.Trace("plugin: ", v.Name(), ", enabled: ", cfg.Enabled(true))

		if cfg.Enabled(true) {
			log.Trace("starting plugin: ", v.Name())
			v.Start(cfg)
			log.Debug("started plugin: ", v.Name())
		}

	}
	log.Debug("all plugins loaded")

	log.Trace("start to start modules")
	for _, v := range m.modules {

		cfg := env.GetModuleConfig(v.Name())

		log.Trace("module: ", v.Name(), ", enabled: ", cfg.Enabled(true))

		if cfg.Enabled(true) {
			log.Trace("starting module: ", v.Name())
			v.Start(cfg)
			log.Debug("started module: ", v.Name())
		}

	}
	log.Debug("all modules started")
}

func Stop() {
	log.Trace("start to stop modules")
	for i := len(m.modules) - 1; i >= 0; i-- {
		v := m.modules[i]
		cfg := env.GetModuleConfig(v.Name())
		if cfg.Enabled(true) {
			log.Trace("stopping module: ", v.Name())
			v.Stop()
			log.Debug("stopped module: ", v.Name())
		}
	}
	log.Debug("all modules stopped")

	log.Trace("start to unload plugins")
	for i := len(m.plugins) - 1; i >= 0; i-- {
		v := m.plugins[i]
		cfg := env.GetPluginConfig(v.Name())
		if cfg.Enabled(true) {
			log.Trace("stopping plugin: ", v.Name())
			v.Stop()
			log.Debug("stopped plugin: ", v.Name())
		}
	}
	log.Debug("all plugins unloaded")
}
