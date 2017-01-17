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
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/asdine/storm/q"
	lz4 "github.com/bkaradzic/go-lz4"
	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/store"
	"github.com/medcl/gopa/core/util"
	"github.com/medcl/gopa/modules/config"
	"path"
	"time"
)

type BoltdbStore struct {
	FileName string
}

var db *storm.DB

func (this BoltdbStore) Open() error {

	//loading or initializing boltdb
	if util.IsExist(this.FileName) {
		log.Debug("found boltdb file, start reload,", this.FileName)
	}

	var err error
	v := global.Lookup(config.REGISTER_BOLTDB)
	if v != nil {
		boltDb := v.(*bolt.DB)
		db, err = storm.Open("boltdb", storm.UseDB(boltDb), storm.Codec(protobuf.Codec))
	} else {
		file := path.Join(global.Env().SystemConfig.Data, "boltdb")
		db, err = storm.Open(file, storm.BoltOptions(0600, &bolt.Options{Timeout: 5 * time.Second}), storm.Codec(protobuf.Codec))
	}
	if err != nil {
		log.Errorf("error open boltdb: %s, %s", this.FileName, err)
		return err
	}

	buckets := []string{config.TaskBucketKey, config.StatsBucketKey, config.SnapshotBucketKey, config.SnapshotMappingBucketKey}
	for _, bucket := range buckets {
		db.Bolt.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				log.Error("create bucket: ", err, ",", bucket)
				panic(err)
			}
			return nil
		})
	}

	global.Register(config.REGISTER_BOLTDB, db)

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
	db.Bolt.View(func(tx *bolt.Tx) error {
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
	db.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(key, value)
		return err
	})
	return nil
}

func (filter BoltdbStore) DeleteValue(bucket string, key []byte, value []byte) error {
	db.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Delete(key)
		return err
	})
	return nil
}

func (filter BoltdbStore) DeleteBucket(bucket string, key []byte, value []byte) error {
	db.Bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.DeleteBucket(key)
		return err
	})
	return nil
}

func (filter BoltdbStore) Get(key string, value interface{}, to interface{}) error {
	return db.One(key, value, to)
}

func (filter BoltdbStore) Save(o interface{}) error {
	return db.Save(o)
}

func (filter BoltdbStore) Update(o interface{}) error {
	return db.Update(o)
}

func (filter BoltdbStore) Delete(o interface{}) error {
	return db.DeleteStruct(o)
}

func (filter BoltdbStore) Count(o interface{}) (int, error) {
	return db.Count(o)
}

func (filter BoltdbStore) Search(t1, t2 interface{}, q1 *store.Query) (error, store.Result) {
	result := store.Result{}
	total, err := store.Count(t1)
	if err != nil {
		log.Debug(err)
		total = -1
	}
	result.Total = total
	result.Result = t2

	if q1.From < 0 {
		q1.From = 0
	}
	if q1.Size < 0 {
		q1.Size = 10
	}

	var q2 storm.Query
	if q1.Filter != nil {
		q2 = db.Select(q.Eq(q1.Filter.Name, q1.Filter.Value)) //can't limit here, bug .Limit(q1.Size).Skip(q1.From)

	} else {
		q2 = db.Select(q.True()).Limit(q1.Size).Skip(q1.From)
	}

	if q1.Sort != "" {
		q2 = q2.OrderBy(q1.Sort).Reverse()
	}

	//t, _ := time.Parse(layout, skipDate)
	//query := db.Select(q.Gt("CreateTime", t)).Limit(size).Skip(from).Reverse().OrderBy("CreateTime")
	err = q2.Find(t2)
	if err != nil {
		log.Trace(err)
	}
	return err, result
}

