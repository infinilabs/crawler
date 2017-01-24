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

import "sync"

type Stats struct {
	ID   string                       `storm:"id,unique" json:"id" gorm:"not null;unique;primary_key"`
	Data *map[string]map[string]int64 `storm:"inline" json:"data,omitempty"`
}

type StatsInterface interface {
	Increment(category, key string)

	IncrementBy(category, key string, value int64)

	Decrement(category, key string)

	DecrementBy(category, key string, value int64)

	Timing(category, key string, v int64)

	Gauge(category, key string, v int64)

	Stat(category, key string) int64

	StatsAll() *[]byte
}

var handlers []StatsInterface

func Increment(category, key string) {
	IncrementBy(category, key, 1)
}

func IncrementBy(category, key string, value int64) {
	for _, v := range handlers {
		v.IncrementBy(category, key, value)
	}
}

func Decrement(category, key string) {
	DecrementBy(category, key, 1)
}

func DecrementBy(category, key string, value int64) {
	for _, v := range handlers {
		v.DecrementBy(category, key, value)
	}
}

func Timing(category, key string, value int64) {
	for _, v := range handlers {
		v.Timing(category, key, value)
	}
}

func Gauge(category, key string, value int64) {
	for _, v := range handlers {
		v.Gauge(category, key, value)
	}
}

func Stat(category, key string) int64 {
	for _, v := range handlers {
		b := v.Stat(category, key)
		if b > 0 {
			return b
		}
	}
	return 0
}

func StatsAll() *[]byte {
	for _, v := range handlers {
		b := v.StatsAll()
		if b != nil {
			return b
		}
	}
	return &[]byte{}
}

var inited bool
var l sync.RWMutex

func Init() {
	if inited {
		return
	}
	l.Lock()
	handlers = []StatsInterface{}
	inited = true
	l.Unlock()
}

func Register(h StatsInterface) {
	Init()
	handlers = append(handlers, h)
}
