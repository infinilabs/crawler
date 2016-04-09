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
	. "github.com/medcl/gopa/core/config"
)

type Env struct {
	//Logger logging.Logger
	*Registrar
	SystemConfig  *SystemConfig
	RuntimeConfig *RuntimeConfig
	Channels      *Channels
}

func Environment(registrar *Registrar, sysConfig *SystemConfig, runtimeConfig *RuntimeConfig) *Env {
	//if logger == nil {
	//	logger = logging.NullLogger{}
	//}

	env := Env{}
	env.RuntimeConfig = runtimeConfig
	env.SystemConfig = sysConfig
	env.Channels = &Channels{}
	env.Channels.PendingFetchUrl = make(chan []byte)
	env.Registrar = registrar
	//env.Logger = logger

	return &env
}

func NullEnv() *Env {
	return Environment(nil, nil, nil)
}

type Channels struct {
	PendingFetchUrl chan []byte
}

type SystemConfig struct {
	Version string `0.0.1`
}
