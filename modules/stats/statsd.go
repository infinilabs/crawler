package statsd

import (
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/modules/stats/statsd"
	"sync"
	"time"
	"github.com/medcl/gopa/core/stats"
	."github.com/medcl/gopa/core/env"
)

func (this StatsDModule) Name() string {
	return "StatsD"
}

func (this StatsDModule) Start(env *Env) {
	initStatsd()
	stats.Register(this)
}

func (this StatsDModule) Stop() error {
	if(statsdbuffer!=nil){
		statsdbuffer.Close()
	}
	return nil
}

type StatsDModule struct {
}

var statsdInited bool
var statsdbuffer *statsd.StatsdBuffer
var l1 sync.RWMutex

func initStatsd() {
	if statsdInited {
		return
	}
	l1.Lock()
	prefix := "gopa."
	statsdclient := statsd.NewStatsdClient("statsdhost:8125", prefix) //TODO configable
	err := statsdclient.CreateSocket()
	if nil != err {
		log.Warn(err)
		return
	}

	interval := time.Second * 1 // aggregate stats and flush every 1 seconds
	statsdbuffer = statsd.NewStatsdBuffer(interval, statsdclient)

	statsdInited = true
	l1.Unlock()
}

func (this StatsDModule) Increment(category, key string) {
	this.IncrementBy(category, key, 1)
}

func (this StatsDModule) IncrementBy(category, key string, value int64) {
	initStatsd()
	statsdbuffer.Incr(category+"."+key, value)
}

func (this StatsDModule) Decrement(category, key string) {
	this.DecrementBy(category, key, 1)
}

func (this StatsDModule) DecrementBy(category, key string, value int64) {
	initStatsd()
	statsdbuffer.Decr(category+"."+key, value)
}

func (this StatsDModule) Timing(category, key string, v int64) {
	initStatsd()
	statsdbuffer.Timing(category+"."+key, v)

}

func (this StatsDModule) Gauge(category, key string, v int64) {
	initStatsd()
	statsdbuffer.Gauge(category+"."+key, v)
}

func (this StatsDModule) Stat(category, key string) int64 {
	return 0
}

func (this StatsDModule) StatsAll() *[]byte {	return nil
}
