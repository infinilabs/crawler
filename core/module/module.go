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
	. "github.com/medcl/gopa/core/env"
)

type Modules struct {
	env     *Env
	modules []Module
}

var m *Modules

func New(env *Env) {
	mod := Modules{}
	mod.env = env
	m= &mod
}

func Register(mod Module)  {
	m.modules = append(m.modules,mod)
}


func Start() {

	log.Trace("start to start modules")
	for _, v := range m.modules{
		v.Start(m.env)
	}
	log.Trace("all modules started")
}

func Stop() {
	log.Trace("start to stop modules")
	for _, v := range m.modules{
		v.Stop()
	}
	log.Info("all modules stopeed")
}
