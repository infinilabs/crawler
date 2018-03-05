package filter

import (
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/filter"
	"github.com/infinitbyte/gopa/modules/filter/kv"
)

type FilterModule struct {
}

type FilterConfig struct {
	Driver string `config:"driver"`
	Bloom  *BloomFilterConfig
	Cuckoo *CuckooFilterConfig
	KV     *KVFilterConfig
}

type BloomFilterConfig struct{}
type CuckooFilterConfig struct{}
type KVFilterConfig struct{}

var (
	defaultConfig = FilterConfig{
		Driver: "kv",
		Bloom:  &BloomFilterConfig{},
		Cuckoo: &CuckooFilterConfig{},
		KV:     &KVFilterConfig{},
	}
)

func (module FilterModule) Name() string {
	return "Filter"
}

var handler filter.Filter

func (module FilterModule) Start(cfg *Config) {

	//init config
	cfg.Unpack(&defaultConfig)

	if defaultConfig.Driver == "kv" {
		handler = kv.KVFilter{}
		filter.Register(handler)
		return

	} else if defaultConfig.Driver == "bloom" {
		//TODO
	} else if defaultConfig.Driver == "cuckoo" {
		//TODO
	} else {
		panic(errors.Errorf("invalid driver, %s", defaultConfig.Driver))
	}
}

func (module FilterModule) Stop() error {
	if handler != nil {
		err := handler.Close()
		if err != nil {
			log.Error(err)
		}
	}
	return nil

}
