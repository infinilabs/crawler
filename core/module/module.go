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
	"github.com/medcl/gopa/core/env"
)

type Modules struct {
	modules []Module
	configs map[string]interface{}
}

var m *Modules

func New() {
	mod := Modules{}
	m = &mod
}

func Register(mod Module) {
	m.modules = append(m.modules, mod)
}

func Start() {

	log.Trace("start to start modules")
	for _, v := range m.modules {

		cfg := env.GetModuleConfig(v.Name())

		log.Trace("module: ", v.Name(), ", enabled: ", cfg.Enabled())

		if cfg.Enabled() {
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
		if cfg.Enabled() {
			log.Trace("stopping module: ", v.Name())
			v.Stop()
			log.Debug("stoped module: ", v.Name())
		}
	}
	log.Debug("all modules stopeed")
}
