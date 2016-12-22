package statsd

import (
	log "github.com/cihub/seelog"
	"sync"
	"github.com/medcl/gopa/core/stats"
	"runtime"
	"github.com/medcl/gopa/core/tasks"
	"encoding/json"
	."github.com/medcl/gopa/core/env"
)

func (this StatsStoreModule) Name() string {
	return "StatsStore"
}

func (this StatsStoreModule) Start(env *Env) {
	initStats()
	stats.Register(this)
}

func (this StatsStoreModule) Stop() error {
	tasks.Save(s)
	log.Trace("save stats db,", s.ID)
	return nil
}

type StatsStoreModule struct {
}

var s *stats.Stats
var inited bool
var l sync.RWMutex

func initData(category, key string) {
	initStats()

	l.Lock()
	_, ok := (*s.Data)[category]
	if !ok {
		(*s.Data)[category] = make(map[string]int64)
	}
	_, ok1 := (*s.Data)[category][key]
	if !ok1 {
		(*s.Data)[category][key] = 0
	}
	l.Unlock()
	runtime.Gosched()
}

func (this StatsStoreModule) Increment(category, key string) {
	this.IncrementBy(category, key, 1)
}

func  (this StatsStoreModule)IncrementBy(category, key string, value int64) {
	initData(category, key)
	l.Lock()
	(*s.Data)[category][key] += value
	l.Unlock()
	runtime.Gosched()
}

func  (this StatsStoreModule)Decrement(category, key string) {
	this.DecrementBy(category, key, 1)
}

func  (this StatsStoreModule)DecrementBy(category, key string, value int64) {
	initData(category, key)
	l.Lock()
	(*s.Data)[category][key] -= value
	l.Unlock()
	runtime.Gosched()
}

func  (this StatsStoreModule)Timing(category, key string, v int64) {

}

func  (this StatsStoreModule)Gauge(category, key string, v int64) {

}

func  (this StatsStoreModule)Stat(category, key string) int64 {
	initData(category, key)
	l.RLock()
	v := ((*s.Data)[category][key])
	l.RUnlock()
	return v
}

func  (this StatsStoreModule)StatsAll() *[]byte {
	initStats()
	l.RLock()
	defer l.RUnlock()
	b, _ := json.MarshalIndent((*s.Data), "", " ")
	return &b
}

func initStats() {
	if inited {
		return
	}
	l.Lock()
	defer l.Unlock()
	if s == nil {
		s = &stats.Stats{}
		s.ID = "statsd"
		err := tasks.Get("ID", "statsd", s)
		if err != nil {
		}
	}

	if s.Data == nil {
		s.Data = &map[string]map[string]int64{}
		log.Trace("inited stats map")
	}
	inited = true
}