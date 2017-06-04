package filter

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/filter"
	. "github.com/medcl/gopa/core/filter"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/modules/config"
	"github.com/medcl/gopa/modules/filter/impl"
	"path"
	"sync"
)

type FilterModule struct {
}

func (this FilterModule) Name() string {
	return "Filter"
}

func (this FilterModule) Exists(bucket FilterKey, key []byte) bool {
	f := filters[bucket]
	return f.Exists(key)
}

func (this FilterModule) Add(bucket FilterKey, key []byte) error {
	f := filters[bucket]
	return f.Add(key)
}

func (this FilterModule) Delete(bucket FilterKey, key []byte) error {
	f := filters[bucket]
	return f.Delete(key)
}

var l sync.RWMutex

func (this FilterModule) CheckThenAdd(bucket FilterKey, key []byte) (b bool, err error) {
	f := filters[bucket]
	l.Lock()
	defer l.Unlock()
	b = f.Exists(key)
	if !b {
		err = f.Add(key)
	}
	return b, err
}

func initFilter(key FilterKey) {
	//f := impl.EmptyFilter{}
	f := impl.LeveldbFilter{}
	file := path.Join(global.Env().SystemConfig.GetDataDir(), "filters", string(key))
	err := f.Open(file)
	if err != nil {
		panic(err)
	}

	filters[key] = &f
}

var filters map[FilterKey]*impl.LeveldbFilter

func (this FilterModule) Start(cfg *Config) {

	filters = map[FilterKey]*impl.LeveldbFilter{}

	//TODO dynamic config
	initFilter(config.DispatchFilter)
	initFilter(config.FetchFilter)
	initFilter(config.CheckFilter)

	filter.Regsiter(this)
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
