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
	"sync"
	"runtime"
)

type RegisterKey string

type Registrar struct {
	values map[RegisterKey]interface{}
	sync.Mutex
}

var(
	r *Registrar
	l sync.RWMutex
	inited bool
)

func GetRegistrar()*Registrar  {
	if !inited {
		l.Lock()
		if(!inited){
			r = &Registrar{values: map[RegisterKey]interface{}{}}
			inited = true
		}
		l.Unlock()
		runtime.Gosched()
	}
	return r
}

func Register(k RegisterKey, v interface{}) {
	reg:=GetRegistrar()
	if reg == nil {
		return
	}

	reg.Lock()
	defer reg.Unlock()
	reg.values[k] = v
}

func  Lookup(k RegisterKey) interface{} {
	reg:=GetRegistrar()
	if reg == nil {
		return nil
	}

	reg.Lock()
	defer reg.Unlock()
	return reg.values[k]
}






