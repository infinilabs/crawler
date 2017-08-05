package impl

import (
	log "github.com/cihub/seelog"
	f "github.com/seiflotfy/cuckoofilter"
	"sync"
)

type CuckooFilterImpl struct {
	persistFileName string
	l               sync.Mutex
	cf              *f.CuckooFilter
}

func (filter *CuckooFilterImpl) Open(fileName string) error {

	filter.cf = f.NewDefaultCuckooFilter()
	return nil
}

func (filter *CuckooFilterImpl) Close() error {
	filter.l.Lock()
	defer filter.l.Unlock()
	log.Debug("start persist leveldb, file:", filter.persistFileName)
	return nil

}

func (filter *CuckooFilterImpl) Exists(key []byte) bool {
	filter.l.Lock()
	defer filter.l.Unlock()
	return filter.cf.Lookup(key)
}

func (filter *CuckooFilterImpl) Add(key []byte) error {
	filter.l.Lock()
	defer filter.l.Unlock()
	filter.cf.Insert(key)
	return nil
}

func (filter *CuckooFilterImpl) Delete(key []byte) error {
	filter.l.Lock()
	defer filter.l.Unlock()
	filter.cf.Delete(key)
	return nil
}
