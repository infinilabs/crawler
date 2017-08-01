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

package global

import (
	"errors"
	"github.com/infinitbyte/gopa/core/env"
	"runtime"
	"sync"
)

// RegisterKey is used to register custom value and retrieve back
type RegisterKey string

type registrar struct {
	values map[RegisterKey]interface{}
	sync.Mutex
}

var (
	r      *registrar
	l      sync.RWMutex
	inited bool
	e      *env.Env
)

func getRegistrar() *registrar {
	if !inited {
		l.Lock()
		if !inited {
			r = &registrar{values: map[RegisterKey]interface{}{}}
			inited = true
		}
		l.Unlock()
		runtime.Gosched()
	}
	return r
}

// Register is used to register your own key and value
func Register(k RegisterKey, v interface{}) {
	reg := getRegistrar()
	if reg == nil {
		return
	}

	reg.Lock()
	defer reg.Unlock()
	reg.values[k] = v
}

// Lookup is to lookup your own previous registered value
func Lookup(k RegisterKey) interface{} {
	reg := getRegistrar()
	if reg == nil {
		return nil
	}

	reg.Lock()
	defer reg.Unlock()
	return reg.values[k]
}

// RegisterEnv is used to register env to this register hub
func RegisterEnv(e1 *env.Env) {
	e = e1
}

// Env returns registered env, should be available globally
func Env() *env.Env {
	if e == nil {
		panic(errors.New("env is not inited"))
	}
	return e
}
