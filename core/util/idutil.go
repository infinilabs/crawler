package util

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/rs/xid"
	"io/ioutil"
	"sync"
	"sync/atomic"
)

var lock sync.Mutex

// GetUUID return a generated uuid
func GetUUID() string {
	lock.Lock()
	defer lock.Unlock()
	return xid.New().String()
}

type incrementCounter struct {
	l  *sync.RWMutex
	ID map[string]*atomicID
}

var count = incrementCounter{l: &sync.RWMutex{}, ID: make(map[string]*atomicID)}

type atomicID struct {
	Sequence int64
}

func (id *atomicID) Increment() int64 {
	return atomic.AddInt64(&id.Sequence, 1)
}

var lock1 sync.Mutex
var persistedPath string

// GetIncrementID return incremented id in specify bucket
func GetIncrementID(bucket string) string {

	count.l.Lock()
	o := count.ID[bucket]
	if o == nil {
		o = &atomicID{}
		count.ID[bucket] = o
	}
	v := o.Increment()
	count.l.Unlock()
	return fmt.Sprintf("%d", v)
}

// SnapshotPersistID will make a snapshot and persist id stats to disk
func SnapshotPersistID() {
	lock1.Lock()
	defer lock1.Unlock()

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(count)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(persistedPath, buf.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

// RestorePersistID will take the snapshot and restore to id seeds
func RestorePersistID(path string) {
	lock1.Lock()
	defer lock1.Unlock()

	persistedPath = JoinPath(path, ".sequence")

	if !FileExists(persistedPath) {
		return
	}

	n, err := ioutil.ReadFile(persistedPath)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewReader(n)
	err = gob.NewDecoder(buf).Decode(&count)
	if err != nil {
		panic(err)
	}
}
