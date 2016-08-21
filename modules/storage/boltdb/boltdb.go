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
	"bytes"
	"strconv"

	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
	util "github.com/medcl/gopa/core/util"
	"time"
)

type BoltdbStore struct {
	WalkPrefix         string
	FetchPrefix        string
	ParsePrefix        string
	PendingFetchPrefix string
	OffsetPrefix       string

	PersistFileName string
	DB              *bolt.DB
}

//store webpage file
func (this *BoltdbStore) Store(url string, data []byte) {
	util.FilePutContentWithByte(url, data)
}

//get webpage file
func (this *BoltdbStore) Get(key string) []byte {
	file, error := util.FileGetContent(key)
	if error != nil {
		log.Error("get file:", key, error)
	}
	return file
}

func (this *BoltdbStore) List(from int, size int) [][]byte {
	return nil
}

func (this *BoltdbStore) TaskEnqueue(url []byte) {
	log.Info("task enqueue:", string(url))
}

func (this *BoltdbStore) Open() error {

	//var runtimeConfig= config.InitOrGetConfig()
	this.WalkPrefix = "walk"
	this.FetchPrefix = "fetch"
	this.ParsePrefix = "parse"
	this.PendingFetchPrefix = "pfetch"
	this.OffsetPrefix = "offset"

	this.PersistFileName = "boltdb"

	//loading or initializing boltdb
	log.Debug("found boltdb file, start reload,", this.PersistFileName)

	db, err := bolt.Open(this.PersistFileName, 0600, &bolt.Options{Timeout: 5 * time.Second})
	this.DB = db

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(FilterBucketKey))
		if err != nil {
			return log.Error("create bucket: ", err)
		}
		return nil
	})

	if err != nil {
		log.Error("boltdb:", this.PersistFileName, err)
		return err
	}

	log.Info("boltdb successfully reloaded:", this.PersistFileName)

	return nil
}

func (this *BoltdbStore) Close() error {
	err := this.DB.Close()
	if err != nil {
		log.Error("boltdb:", this.PersistFileName, err)
	}
	return err
}

func (this *BoltdbStore) PersistBloomFilter() {

}

func (this *BoltdbStore) UrlHasWalked(url []byte) bool {

	c := [][]byte{[]byte(this.WalkPrefix), url}

	return this.Lookup(bytes.Join(c, []byte(":")))
}

func (this *BoltdbStore) UrlHasFetched(url []byte) bool {
	c := [][]byte{[]byte(this.FetchPrefix), url}

	return this.Lookup(bytes.Join(c, []byte(":")))
}

func (this *BoltdbStore) FileHasParsed(url []byte) bool {
	c := [][]byte{[]byte(this.ParsePrefix), url}

	return this.Lookup(bytes.Join(c, []byte(":")))
}

func (this *BoltdbStore) PendingFetchUrlHasAdded(url []byte) bool {
	c := [][]byte{[]byte(this.PendingFetchPrefix), url}

	key := bytes.Join(c, []byte(":"))
	result := this.Lookup(key)

	value := this.GetValue(key)

	log.Trace("check pending url error,", string(key), ",", result, ",value:", string(value))
	return result
}

func (this *BoltdbStore) AddWalkedUrl(url []byte) {
	c := [][]byte{[]byte(this.WalkPrefix), url}

	this.Add(bytes.Join(c, []byte(":")))
}

func (this *BoltdbStore) AddPendingFetchUrl(url []byte) {
	c := [][]byte{[]byte(this.PendingFetchPrefix), url}

	key := bytes.Join(c, []byte(":"))
	err := this.Add(key)

	if err != nil {
		log.Error("add pending url error,", url, ",", err)
	}
}

func (this *BoltdbStore) AddSavedUrl(url []byte) {
	this.AddWalkedUrl(url)
	this.AddFetchedUrl(url)
}

func (this *BoltdbStore) LogSavedFile(path string, content string) {
	util.FileAppendNewLine(path, content)
}

func (this *BoltdbStore) LogPendingFetchUrl(path string, content string) {
	util.FileAppendNewLine(path, content)
}

func (this *BoltdbStore) LogFetchFailedUrl(path string, content string) {
	this.AddFetchFailedUrl([]byte(path))
	util.FileAppendNewLine(path, content)
}

func (this *BoltdbStore) AddFetchedUrl(url []byte) {
	c := [][]byte{[]byte(this.FetchPrefix), url}

	this.Add(bytes.Join(c, []byte(":")))

}

func (this *BoltdbStore) saveFetchedUrlToLocalFile(path string, url string) {
	util.FileAppendNewLine(path, url)
}

func (this *BoltdbStore) AddParsedFile(url []byte) {
	c := [][]byte{[]byte(this.ParsePrefix), url}

	this.Add(bytes.Join(c, []byte(":")))
}

func (this *BoltdbStore) AddFetchFailedUrl(url []byte) {
	c := [][]byte{[]byte(this.WalkPrefix), url}

	this.Add(bytes.Join(c, []byte(":")))

	log.Debug("fetch failed url:", string(url))
}

func (this *BoltdbStore) FileHasSaved(file string) bool {
	log.Debug("start check file:", file)
	return util.CheckFileExists(file)
}

func (this *BoltdbStore) LoadOffset(fileName string) int64 {
	log.Debug("start load offsets,", fileName)

	c := [][]byte{([]byte)(this.OffsetPrefix), []byte(fileName)}

	n := this.GetValue(bytes.Join(c, []byte(":")))

	if n != nil {
		ret, err := strconv.ParseInt(string(n), 10, 64)
		if err != nil {
			log.Error("offset parse error, ", fileName, " , ", err)
			return 0
		}
		log.Debug("offset load successfully, ", fileName, " : ", ret)
		return int64(ret)
	}
	log.Debug("hit default offsets,", fileName)
	return 0
}

func (this *BoltdbStore) PersistOffset(fileName string, offset int64) {
	//persist worker's offset

	c := [][]byte{[]byte(this.OffsetPrefix), []byte(fileName)}

	error := this.AddValue(bytes.Join(c, []byte(":")), []byte(strconv.FormatInt(offset, 10)))

	if error != nil {
		log.Error(fileName, error)
		return
	}
}

func (this *BoltdbStore) InitPendingFetchBloomFilter(fileName string) {}

//TODO REMOVE
func (filter *BoltdbStore) Persist() error {

	log.Debug("boltdb start persist,file:", filter.PersistFileName)

	log.Info("boltdb safety persisted.")

	return nil
}

func (filter *BoltdbStore) Lookup(key []byte) bool {
	value := filter.GetValue(key)

	if value != nil {
		log.Trace("return true,hit key, ", string(key), " : ", value)
		return true
	}
	log.Trace("return false, hit key", string(key), " : ", value)

	return false
}

const FilterBucketKey string = "UrlFilter"

func (filter *BoltdbStore) Add(key []byte) error {
	log.Trace("add key,", string(key))
	filter.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilterBucketKey))
		err := b.Put(key, []byte("true"))
		return err
	})
	return nil
}

func (filter *BoltdbStore) GetValue(key []byte) []byte {
	var ret []byte = nil
	filter.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilterBucketKey))
		v := b.Get(key)
		if v != nil {
			ret = v
		}
		return nil
	})
	return ret

}

func (filter *BoltdbStore) AddValue(key []byte, value []byte) error {
	filter.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(FilterBucketKey))
		err := b.Put(key, value)
		return err
	})
	return nil
}
