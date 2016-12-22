/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package boltdb

import (
	lz4 "github.com/bkaradzic/go-lz4"
	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/util"
	"time"
	"github.com/medcl/gopa/core/global"
)

type BoltdbStore struct {
	FileName string
}

var db       *bolt.DB

func (this BoltdbStore) Open() error {

	//loading or initializing boltdb
	if util.IsExist(this.FileName) {
		log.Debug("found boltdb file, start reload,", this.FileName)
	}

	var err error
	db, err = bolt.Open(this.FileName, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Error("error open boltdb:", this.FileName, err)
		return err
	}

	buckets := []string{TaskBucketKey, StatsBucketKey, SnapshotBucketKey}
	for _, bucket := range buckets {
		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				log.Error("create bucket: ", err, ",", bucket)
				panic(err)
			}
			return nil
		})
	}

	global.Register(global.REGISTER_BOLTDB, db)

	log.Debug("boltdb successfully started:", this.FileName)

	return nil
}

func (this BoltdbStore) Close() error {
	err := db.Close()
	if err != nil {
		log.Error("boltdb:", this.FileName, err)
	}
	return err
}

const TaskBucketKey string = "Task"
const StatsBucketKey string = "Stats"
const SnapshotBucketKey string = "Snapshot"

func (filter BoltdbStore) GetCompressedValue(bucket string, key []byte) []byte {

	data := filter.GetValue(bucket, key)
	data, err := lz4.Decode(nil, data)
	if err != nil {
		log.Error("Failed to decode:", err)
		return nil
	}
	return data
}

func (filter BoltdbStore) GetValue(bucket string, key []byte) []byte {
	var ret []byte = nil
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get(key)
		if v != nil {
			ret = v
		}
		return nil
	})
	return ret
}

func (filter BoltdbStore) AddValueCompress(bucket string, key []byte, value []byte) error {
	value, err := lz4.Encode(nil, value)
	if err != nil {
		log.Error("Failed to encode:", err)
		return err
	}
	return filter.AddValue(bucket, key, value)
}

func (filter BoltdbStore) AddValue(bucket string, key []byte, value []byte) error {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(key, value)
		return err
	})
	return nil
}

func (filter BoltdbStore) DeleteValue(bucket string, key []byte, value []byte) error {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Delete(key)
		return err
	})
	return nil
}

func (filter BoltdbStore) DeleteBucket(bucket string, key []byte, value []byte) error {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.DeleteBucket(key)
		return err
	})
	return nil
}
