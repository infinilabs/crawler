package stats

import (
	"fmt"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/stats"
	"sync"
	"github.com/quipo/statsd"
)

type StatsDConfig struct {
	Host              string        `config:"host"`
	Port              int           `config:"port"`
	Namespace         string        `config:"namespace"`
}
type StatsDModule struct {
}

var statsdInited bool
var statsdclient *statsd.StatsdClient
var l1 sync.RWMutex

var defaultStatsdConfig = StatsDConfig{
	Host:              "localhost",
	Port:              8125,
	Namespace:         "gopa.",
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

	log.Debug("statsd connec to, ", addr,",prefix:",config.Namespace)

	err := statsdclient.CreateSocket()
	if nil != err {
		log.Warn(err)
		return
	}

	statsdInited = true

	stats.Register(this)
}

func (this StatsDModule) Stop() error {
	if statsdclient != nil {
		statsdclient.Close()
	}
	return nil
}

func (this StatsDModule) Increment(category, key string) {

	this.IncrementBy(category, key, 1)
}

func (this StatsDModule) IncrementBy(category, key string, value int64) {
	statsdclient.Incr(category+"."+key, value)
}

func (this StatsDModule) Decrement(category, key string) {
	this.DecrementBy(category, key, 1)
}

func (this StatsDModule) DecrementBy(category, key string, value int64) {
	statsdclient.Decr(category+"."+key, value)
}

func (this StatsDModule) Timing(category, key string, v int64) {
	statsdclient.Timing(category+"."+key, v)

}

func (this StatsDModule) Gauge(category, key string, v int64) {
	statsdclient.Gauge(category+"."+key, v)
}

func (this StatsDModule) Stat(category, key string) int64 {
	return 0
}

func (this StatsDModule) StatsAll() *[]byte {
	return nil
}
