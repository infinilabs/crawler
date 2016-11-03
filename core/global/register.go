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

import "sync"

type Registrar struct {
	values map[string]interface{}
	sync.Mutex
}

func (r *Registrar) Register(k string, v interface{}) {
	if r == nil {
		return
	}

	r.Lock()
	defer r.Unlock()
	r.values[k] = v
}

func (r *Registrar) Lookup(k string) interface{} {
	if r == nil {
		return nil
	}

	r.Lock()
	defer r.Unlock()
	return r.values[k]
}
