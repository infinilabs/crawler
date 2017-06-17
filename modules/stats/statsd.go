package stats

import (
	"fmt"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/stats"
	"github.com/quipo/statsd"
	"sync"
	"time"
)

type StatsDConfig struct {
	Host      string `config:"host"`
	Port      int    `config:"port"`
	Namespace string `config:"namespace"`
	IntervalInSeconds string `config:"interval_in_seconds"`
}
type StatsDModule struct {
}

var statsdInited bool
var statsdclient *statsd.StatsdClient
var buffer *statsd.StatsdBuffer
var l1 sync.RWMutex

var defaultStatsdConfig = StatsDConfig{
	Host:      "localhost",
	Port:      8125,
	Namespace: "gopa.",
	IntervalInSeconds: 2,
}

func (this StatsDModule) Name() string {
	return "StatsD"
}

func (this StatsDModule) Start(cfg *Config) {
	if statsdInited {
		return
	}

	config := defaultStatsdConfig
	cfg.Unpack(&config)

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	l1.Lock()
	defer l1.Unlock()
	statsdclient = statsd.NewStatsdClient(addr, config.Namespace)

	log.Debug("statsd connec to, ", addr, ",prefix:", config.Namespace)

	err := statsdclient.CreateSocket()
	if nil != err {
		log.Warn(err)
		return
	}

	interval := time.Second * config.IntervalInSeconds // aggregate stats and flush every 2 seconds
	buffer = statsd.NewStatsdBuffer(interval, statsdclient)

	statsdInited = true

	stats.Register(this)
}

func (this StatsDModule) Stop() error {
	if statsdclient != nil {
		buffer.Close()
		statsdclient.Close()
	}
	return nil
}

func (this StatsDModule) Increment(category, key string) {

	this.IncrementBy(category, key, 1)
}

func (this StatsDModule) IncrementBy(category, key string, value int64) {
	if !statsdInited {
		return
	}
	buffer.Incr(category+"."+key, value)
}

func (this StatsDModule) Decrement(category, key string) {
	this.DecrementBy(category, key, 1)
}

func (this StatsDModule) DecrementBy(category, key string, value int64) {
	if !statsdInited {
		return
	}
	buffer.Decr(category+"."+key, value)
}

func (this StatsDModule) Timing(category, key string, v int64) {
	if !statsdInited {
		return
	}
	buffer.Timing(category+"."+key, v)

}

func (this StatsDModule) Gauge(category, key string, v int64) {
	if !statsdInited {
		return
	}
	buffer.Gauge(category+"."+key, v)
}

func (this StatsDModule) Stat(category, key string) int64 {
	return 0
}

func (this StatsDModule) StatsAll() *[]byte {
	return nil
}
