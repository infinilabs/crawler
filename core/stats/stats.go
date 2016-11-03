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
	log "github.com/cihub/seelog"
	"runtime"
	"sync"
	"time"
	"github.com/medcl/gopa/core/stats/statsd"
)

var data map[string]map[string]int64
var inited bool
var statsdInited bool
var statsdbuffer *statsd.StatsdBuffer
var l sync.RWMutex

func initStatsd() {

	prefix := "gopa."
	statsdclient := statsd.NewStatsdClient("statsdhost:8125", prefix) //TODO configable
	err := statsdclient.CreateSocket()
	if nil != err {
		log.Error(err)
		return
	}

	interval := time.Second * 1 // aggregate stats and flush every 1 seconds
	statsdbuffer = statsd.NewStatsdBuffer(interval, statsdclient)
	//defer statsdbuffer.Close()

	statsdInited = true

}

func initData(category, key string) {
	l.Lock()

	if !inited {
		data = make(map[string]map[string]int64)
		inited = true
	}

	if !statsdInited {
		initStatsd()
	}

	_, ok := data[category]
	if !ok {
		data[category] = make(map[string]int64)
	}

	_, ok1 := data[category][key]
	if !ok1 {
		data[category][key] = 0
	}
	l.Unlock()
	runtime.Gosched()
}

func Increment(category, key string) {
	IncrementBy(category, key, 1)
}

func IncrementBy(category, key string, value int64) {
	initData(category, key)

	if statsdInited {
		statsdbuffer.Incr(category+"."+key, value)
	}

	l.Lock()
	data[category][key] += value
	l.Unlock()
	runtime.Gosched()
}

func Decrement(category, key string) {
	DecrementBy(category, key, 1)
}

func DecrementBy(category, key string, value int64) {
	initData(category, key)

	if statsdInited {
		statsdbuffer.Decr(category+"."+key, value)
	}

	l.Lock()
	data[category][key] -= value
	l.Unlock()
	runtime.Gosched()
}

func Timing(category, key string, v int64) {
	if statsdInited {
		statsdbuffer.Timing(category+"."+key, v)
	}
}

func Gauge(category, key string, v int64) {
	if statsdInited {
		statsdbuffer.Gauge(category+"."+key, v)
	}
}

func Stat(category, key string) int64 {
	initData(category, key)
	l.RLock()
	v := (data[category][key])
	l.RUnlock()
	return v
}

func StatsAll() []byte {
	l.RLock()
	defer l.RUnlock()
	b, _ := json.MarshalIndent(data, "", " ")
	return b
}
