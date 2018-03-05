package kv

import (
	"github.com/infinitbyte/gopa/core/persist"
	"sync"
)

type KVFilter struct {
}

var v = []byte("true")
var l sync.RWMutex

func (filter KVFilter) Open() error {
	return nil
}

func (filter KVFilter) Close() error {
	return nil
}

func (filter KVFilter) Exists(bucket string, key []byte) bool {
	b, _ := persist.GetValue(bucket, key)
	return b != nil
}

func (filter KVFilter) Add(bucket string, key []byte) error {
	return persist.AddValue(bucket, key, v)
}

func (filter KVFilter) Delete(bucket string, key []byte) error {
	return persist.DeleteKey(bucket, key)
}

func (filter KVFilter) CheckThenAdd(bucket string, key []byte) (b bool, err error) {
	l.Lock()
	defer l.Unlock()
	b = filter.Exists(bucket, key)
	if !b {
		err = filter.Add(bucket, key)
	}
	return b, err
}
