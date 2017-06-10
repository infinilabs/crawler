package util

import (
	"fmt"
	"sync"
	"github.com/rs/xid"
)


var lock sync.Mutex

func GetUUID() string {
	lock.Lock()
	defer lock.Unlock()
	return xid.New().String()
}

var idseed map[string]int64

var lock1 sync.Mutex

func GetIncrementID(prefix string) string {
	lock1.Lock()
	defer lock1.Unlock()
	if idseed == nil {
		idseed = map[string]int64{}
	}

	if _, ok := idseed[prefix]; !ok {
		idseed[prefix] = int64(0)
	}

	v:=idseed[prefix]
	v++
	idseed[prefix]=v
	id := fmt.Sprintf("%s-%d", prefix,v)
	return id
}
