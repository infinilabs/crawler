package filter

import (
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/filter"
	. "github.com/infinitbyte/gopa/core/filter"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/modules/config"
	"github.com/infinitbyte/gopa/modules/filter/impl"
	"path"
	"sync"
)

type FilterModule struct {
}

func (this FilterModule) Name() string {
	return "Filter"
}

func (this FilterModule) Exists(bucket Key, key []byte) bool {
	f := filters[bucket]
	return f.Exists(key)
}

func (this FilterModule) Add(bucket Key, key []byte) error {
	f := filters[bucket]
	return f.Add(key)
}

func (this FilterModule) Delete(bucket Key, key []byte) error {
	f := filters[bucket]
	return f.Delete(key)
}

var l sync.RWMutex

func (this FilterModule) CheckThenAdd(bucket Key, key []byte) (b bool, err error) {
	f := filters[bucket]
	l.Lock()
	defer l.Unlock()
	b = f.Exists(key)
	if !b {
		err = f.Add(key)
	}
	return b, err
}

func initFilter(key Key) {
	//f := impl.EmptyFilter{}
	f := impl.LeveldbFilter{}
	file := path.Join(global.Env().SystemConfig.GetWorkingDir(), "filters", string(key))
	err := f.Open(file)
	if err != nil {
		panic(err)
	}

	filters[key] = &f
}

var filters map[Key]*impl.LeveldbFilter

func (this FilterModule) Start(cfg *Config) {

	filters = map[Key]*impl.LeveldbFilter{}

	//TODO dynamic config
	initFilter(config.DispatchFilter)
	initFilter(config.FetchFilter)
	initFilter(config.CheckFilter)

	filter.Register(this)
}

func (this FilterModule) Stop() error {
	for _, v := range filters {
		err := (*v).Close()
		if err != nil {
			log.Error(err)
		}
	}
	return nil

}
