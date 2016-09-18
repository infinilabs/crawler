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

package stats

import (
	"encoding/json"
	"runtime"
	"sync"
)

var data map[string]map[string]int
var inited bool
var l sync.RWMutex

func initData(domain, key string) {
	l.Lock()
	if !inited {
		data = make(map[string]map[string]int)
		inited = true
	}

	_, ok := data[domain]
	if !ok {
		data[domain] = make(map[string]int)
	}

	_, ok1 := data[domain][key]
	if !ok1 {
		data[domain][key] = 0
	}
	l.Unlock()
	runtime.Gosched()
}

func Increment(domain, key string) {
	IncrementBy(domain, key, 1)
}

func IncrementBy(domain, key string, value int) {
	initData(domain, key)
	l.Lock()
	data[domain][key] += value
	l.Unlock()
	runtime.Gosched()
}

func Decrement(domain, key string) {
	DecrementBy(domain, key, 1)
}

func DecrementBy(domain, key string, value int) {
	initData(domain, key)
	l.Lock()
	data[domain][key] -= value
	l.Unlock()
	runtime.Gosched()
}

func Stat(domain, key string) int {
	initData(domain, key)
	l.RLock()
	v := (data[domain][key])
	l.RUnlock()
	return v
}

func StatsAll() []byte {
	l.RLock()
	defer l.RUnlock()
	b, _ := json.MarshalIndent(data, "", " ")
	return b
}
