package impl

import (
	f "github.com/seiflotfy/cuckoofilter"
"sync"
	log "github.com/cihub/seelog"
)

type CuckooFilterImpl struct{
	persistFileName string
	l sync.Mutex
	cf *f.CuckooFilter
}


func (this *CuckooFilterImpl) Open(fileName string) error{

	this.cf = f.NewDefaultCuckooFilter()
	return nil
}

func (this *CuckooFilterImpl) Close() error{
	this.l.Lock()
	defer this.l.Unlock()
	log.Debug("start persist leveldb, file:",this.persistFileName)
	return nil

}

func (filter *CuckooFilterImpl) Exists(key []byte) bool{
	filter.l.Lock()
	defer filter.l.Unlock()
	return filter.cf.Lookup(key)
}

func (filter *CuckooFilterImpl) Add(key []byte) error{
	filter.l.Lock()
	defer filter.l.Unlock()
	filter.cf.Insert(key)
	return nil
}

func (filter *CuckooFilterImpl) Delete(key []byte) error{
	filter.l.Lock()
	defer filter.l.Unlock()
	filter.cf.Delete(key)
	return nil
}

