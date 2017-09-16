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

func (module SimpleStatsModule) Name() string {
	return "StatsStore"
}

const id = "simple-stats"

var data *Stats

func (module SimpleStatsModule) Start(cfg *Config) {

	data = &Stats{}
	data.initStats()
	stats.Register(data)
}

func (module SimpleStatsModule) Stop() error {
	data.l.Lock()
	defer data.l.Unlock()
	v, _ := json.Marshal(data.Data)
	data.ID = id
	err := persist.AddValue(string(config.KVBucketKey), []byte(data.ID), v)
	if err != nil {
		log.Error(err)
	}
	log.Trace("save stats db,", data.ID)
	return nil
}

type SimpleStatsModule struct {
}

type Stats struct {
	l    sync.RWMutex
	ID   string                       `storm:"id,unique" json:"id" gorm:"not null;unique;primary_key"`
	Data *map[string]map[string]int64 `storm:"inline" json:"data,omitempty"`
}

func (s *Stats) initData(category, key string) {

	s.l.Lock()
	_, ok := (*s.Data)[category]
	if !ok {
		(*s.Data)[category] = make(map[string]int64)
	}
	_, ok1 := (*s.Data)[category][key]
	if !ok1 {
		(*s.Data)[category][key] = 0
	}
	s.l.Unlock()
	runtime.Gosched()
}

func (s *Stats) Increment(category, key string) {
	s.IncrementBy(category, key, 1)
}

func (s *Stats) IncrementBy(category, key string, value int64) {
	s.initData(category, key)
	s.l.Lock()
	(*data.Data)[category][key] += value
	s.l.Unlock()
	runtime.Gosched()
}

func (s *Stats) Decrement(category, key string) {
	s.DecrementBy(category, key, 1)
}

func (s *Stats) DecrementBy(category, key string, value int64) {
	s.initData(category, key)
	s.l.Lock()
	(*data.Data)[category][key] -= value
	s.l.Unlock()
	runtime.Gosched()
}

func (s *Stats) Timing(category, key string, v int64) {

}

func (s *Stats) Gauge(category, key string, v int64) {

}

func (s *Stats) Stat(category, key string) int64 {
	s.initData(category, key)
	s.l.RLock()
	v := ((*data.Data)[category][key])
	s.l.RUnlock()
	return v
}

func (s *Stats) StatsAll() *[]byte {
	s.initStats()
	s.l.RLock()
	defer s.l.RUnlock()
	b, _ := json.MarshalIndent((*data.Data), "", " ")
	return &b
}

func (s *Stats) initStats() {
	s.ID = id
	v := persist.GetValue(string(config.KVBucketKey), []byte(s.ID))
	d := map[string]map[string]int64{}
	err := json.Unmarshal(v, &d)
	if err != nil {
		log.Debug(err)
	}
	s.Data = &d

	if s.Data == nil {
		s.Data = &map[string]map[string]int64{}
		log.Trace("inited stats map")
	}
}
