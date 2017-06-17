package util

import (
	"fmt"
	"github.com/rs/xid"
	"sync"
)

var lock sync.Mutex

func GetUUID() string {
	lock.Lock()
	defer lock.Unlock()
	return xid.New().String()
}

var idseed map[string]int64

var lock1 sync.Mutex

func GetIncrementID(bucket string) string {
	lock1.Lock()
	defer lock1.Unlock()
	if idseed == nil {
		idseed = map[string]int64{}
	}

	if _, ok := idseed[bucket]; !ok {
		idseed[bucket] = int64(0)
	}

	v := idseed[bucket]
	v++
	idseed[bucket] = v
	id := fmt.Sprintf("%d", v)
	return id
}
