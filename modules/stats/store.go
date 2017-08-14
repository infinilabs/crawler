package stats

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/modules/config"
	"runtime"
	"sync"
)

func (module StatsStoreModule) Name() string {
	return "StatsStore"
}

func (module StatsStoreModule) Start(cfg *Config) {
	initStats()
	stats.Register(module)
}

func (module StatsStoreModule) Stop() error {
	v, _ := json.Marshal(s.Data)
	s.ID = "statsd"
	err := persist.AddValue(string(config.KVBucketKey), []byte(s.ID), v)
	if err != nil {
		log.Error(err)
	}
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

func (module StatsStoreModule) Increment(category, key string) {
	module.IncrementBy(category, key, 1)
}

func (module StatsStoreModule) IncrementBy(category, key string, value int64) {
	initData(category, key)
	l.Lock()
	(*s.Data)[category][key] += value
	l.Unlock()
	runtime.Gosched()
}

func (module StatsStoreModule) Decrement(category, key string) {
	module.DecrementBy(category, key, 1)
}

func (module StatsStoreModule) DecrementBy(category, key string, value int64) {
	initData(category, key)
	l.Lock()
	(*s.Data)[category][key] -= value
	l.Unlock()
	runtime.Gosched()
}

func (module StatsStoreModule) Timing(category, key string, v int64) {

}

func (module StatsStoreModule) Gauge(category, key string, v int64) {

}

func (module StatsStoreModule) Stat(category, key string) int64 {
	initData(category, key)
	l.RLock()
	v := ((*s.Data)[category][key])
	l.RUnlock()
	return v
}

func (module StatsStoreModule) StatsAll() *[]byte {
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
		v := persist.GetValue(string(config.KVBucketKey), []byte(s.ID))
		d := map[string]map[string]int64{}
		err := json.Unmarshal(v, &d)
		if err != nil {
			log.Debug(err)
		}
		s.Data = &d
	}

	if s.Data == nil {
		s.Data = &map[string]map[string]int64{}
		log.Trace("inited stats map")
	}
	inited = true
}
