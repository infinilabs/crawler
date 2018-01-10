package impl

import (
	"github.com/infinitbyte/gopa/core/persist"
)

type PersistFilter struct {
	FilterBucket string
}

var v = []byte("true")

func (filter *PersistFilter) Open(fileName string) error {
	return nil
}

func (filter *PersistFilter) Close() error {
	return nil
}

func (filter *PersistFilter) Exists(key []byte) bool {
	b, _ := persist.GetValue(filter.FilterBucket, key)
	return b != nil
}

func (filter *PersistFilter) Add(key []byte) error {
	return persist.AddValue(filter.FilterBucket, key, v)
}

func (filter *PersistFilter) Delete(key []byte) error {
	return persist.DeleteKey(filter.FilterBucket, key)
}
