package impl

import (
	"github.com/infinitbyte/gopa/core/persist"
)

type BoltdbImpl struct {
	FilterBucket string
}

var v = []byte("")

func (filter *BoltdbImpl) Open(fileName string) error {
	return nil
}

func (filter *BoltdbImpl) Close() error {
	return nil
}

func (filter *BoltdbImpl) Exists(key []byte) bool {
	b, _ := persist.GetValue(filter.FilterBucket, key)
	return b != nil
}

func (filter *BoltdbImpl) Add(key []byte) error {
	return persist.AddValue(filter.FilterBucket, key, v)
}

func (filter *BoltdbImpl) Delete(key []byte) error {
	return persist.DeleteKey(filter.FilterBucket, key)
}
