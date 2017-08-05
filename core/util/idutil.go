package util

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/rs/xid"
	"io/ioutil"
	"sync"
)

var lock sync.Mutex

// GetUUID return a generated uuid
func GetUUID() string {
	lock.Lock()
	defer lock.Unlock()
	return xid.New().String()
}

var idseed map[string]int64

var lock1 sync.Mutex
var persistedPath string

// GetIncrementID return incremented id in specify bucket
func GetIncrementID(bucket string) string {
	lock1.Lock()
	defer lock1.Unlock()
	if idseed == nil {
		if persistedPath == "" {
			panic(errors.New("persistence path is not set"))
		}
		restorePersistID()
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

// SnapshotPersistID will make a snapshot and persist id stats to disk
func SnapshotPersistID() {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(idseed)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(persistedPath, buf.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

// RestorePersistID will take the snapshot and restore to id seeds
func restorePersistID() {
	idseed = map[string]int64{}

	if !FileExists(persistedPath) {
		return
	}

	n, err := ioutil.ReadFile(persistedPath)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewReader(n)
	err = gob.NewDecoder(buf).Decode(&idseed)
	if err != nil {
		panic(err)
	}
}

// SetIDPersistencePath set the persist path
func SetIDPersistencePath(path string) {
	persistedPath = JoinPath(path, ".idseeds")
}
