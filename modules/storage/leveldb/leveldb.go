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

package leveldb

import (
	"bytes"
	"strconv"

	log "github.com/cihub/seelog"
	util "github.com/infinitbyte/gopa/core/util"
	leveldb "github.com/syndtr/goleveldb/leveldb"
)

type LeveldbStore struct {
	WalkPrefix         string
	FetchPrefix        string
	ParsePrefix        string
	PendingFetchPrefix string
	OffsetPrefix       string

	PersistFileName string
	Leveldb         *leveldb.DB
}

//store webpage file
func (this *LeveldbStore) Store(url string, data []byte) {
	util.FilePutContentWithByte(url, data)
}

//get webpage file
func (this *LeveldbStore) Get(key string) []byte {
	file, error := util.FileGetContent(key)
	if error != nil {
		log.Error("get file:", key, error)
	}
	return file
}

func (this *LeveldbStore) List(from int, size int) [][]byte {
	return nil
}

func (this *LeveldbStore) TaskEnqueue(url []byte) {
	log.Info("task enqueue:", string(url))
}

func (this *LeveldbStore) Open() error {

	//var runtimeConfig= config.InitOrGetConfig()
	this.WalkPrefix = "walk"
	this.FetchPrefix = "fetch"
	this.ParsePrefix = "parse"
	this.PendingFetchPrefix = "pfetch"
	this.OffsetPrefix = "offset"

	this.PersistFileName = "leveldb"

	//loading or initializing leveldb
	log.Debug("found leveldb file, start reload,", this.PersistFileName)

	db, err := leveldb.OpenFile(this.PersistFileName, nil)
	this.Leveldb = db

	if err != nil {
		log.Error("leveldb:", this.PersistFileName, err)
		return err
	}

	log.Info("leveldb successfully reloaded:", this.PersistFileName)

	return nil
}

func (this *LeveldbStore) Close() error {
	err := this.Leveldb.Close()
	if err != nil {
		log.Error("leveldb:", this.PersistFileName, err)
	}
	return err
}

func (this *LeveldbStore) PersistBloomFilter() {

}

func (this *LeveldbStore) UrlHasWalked(url []byte) bool {

	c := [][]byte{[]byte(this.WalkPrefix), url}

	return this.Lookup(bytes.Join(c, []byte(":")))
}

func (this *LeveldbStore) UrlHasFetched(url []byte) bool {
	c := [][]byte{[]byte(this.FetchPrefix), url}

	return this.Lookup(bytes.Join(c, []byte(":")))
}

func (this *LeveldbStore) FileHasParsed(url []byte) bool {
	c := [][]byte{[]byte(this.ParsePrefix), url}

	return this.Lookup(bytes.Join(c, []byte(":")))
}

func (this *LeveldbStore) PendingFetchUrlHasAdded(url []byte) bool {
	c := [][]byte{[]byte(this.PendingFetchPrefix), url}

	key := bytes.Join(c, []byte(":"))
	result := this.Lookup(key)

	value := this.GetValue(key)

	log.Trace("check pending url error,", string(key), ",", result, ",value:", string(value))
	return result
}

func (this *LeveldbStore) AddWalkedUrl(url []byte) {
	c := [][]byte{[]byte(this.WalkPrefix), url}

	this.Add(bytes.Join(c, []byte(":")))
}

func (this *LeveldbStore) AddPendingFetchUrl(url []byte) {
	c := [][]byte{[]byte(this.PendingFetchPrefix), url}

	key := bytes.Join(c, []byte(":"))
	err := this.Add(key)

	if err != nil {
		log.Error("add pending url error,", url, ",", err)
	}
}

func (this *LeveldbStore) AddSavedUrl(url []byte) {
	this.AddWalkedUrl(url)
	this.AddFetchedUrl(url)
}

func (this *LeveldbStore) LogSavedFile(path string, content string) {
	util.FileAppendNewLine(path, content)
}

func (this *LeveldbStore) LogPendingFetchUrl(path string, content string) {
	util.FileAppendNewLine(path, content)
}

func (this *LeveldbStore) LogFetchFailedUrl(path string, content string) {
	this.AddFetchFailedUrl([]byte(path))
	util.FileAppendNewLine(path, content)
}

func (this *LeveldbStore) AddFetchedUrl(url []byte) {
	c := [][]byte{[]byte(this.FetchPrefix), url}

	this.Add(bytes.Join(c, []byte(":")))

}

func (this *LeveldbStore) saveFetchedUrlToLocalFile(path string, url string) {
	util.FileAppendNewLine(path, url)
}

func (this *LeveldbStore) AddParsedFile(url []byte) {
	c := [][]byte{[]byte(this.ParsePrefix), url}

	this.Add(bytes.Join(c, []byte(":")))
}

func (this *LeveldbStore) AddFetchFailedUrl(url []byte) {
	c := [][]byte{[]byte(this.WalkPrefix), url}

	this.Add(bytes.Join(c, []byte(":")))

	log.Debug("fetch failed url:", string(url))
}

func (this *LeveldbStore) FileHasSaved(file string) bool {
	log.Debug("start check file:", file)
	return util.FileExists(file)
}

func (this *LeveldbStore) LoadOffset(fileName string) int64 {
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

func (this *LeveldbStore) PersistOffset(fileName string, offset int64) {
	//persist worker's offset

	c := [][]byte{[]byte(this.OffsetPrefix), []byte(fileName)}

	error := this.AddValue(bytes.Join(c, []byte(":")), []byte(strconv.FormatInt(offset, 10)))

	if error != nil {
		log.Error(fileName, error)
		return
	}
}

func (this *LeveldbStore) InitPendingFetchBloomFilter(fileName string) {}

//TODO REMOVE
func (filter *LeveldbStore) Persist() error {

	log.Debug("leveldb start persist,file:", filter.PersistFileName)

	log.Info("leveldb safety persisted.")

	return nil
}

func (filter *LeveldbStore) Lookup(key []byte) bool {
	value := filter.GetValue(key)

	if value != nil {
		log.Trace("return true,hit key, ", string(key), " : ", value)
		return true
	}
	log.Trace("return false, hit key", string(key), " : ", value)

	return false
}

func (filter *LeveldbStore) Add(key []byte) error {
	log.Trace("add key,", string(key))

	return filter.Leveldb.Put(key, []byte("true"), nil)
}

func (filter *LeveldbStore) GetValue(key []byte) []byte {

	value, err := filter.Leveldb.Get(key, nil)
	if err != nil {
		log.Trace("leveldb getValue error, ", err, " , ", string(key), " : ", value)
		return value
	}

	log.Trace("get key, ", err, " , ", string(key), " : ", value)

	return value
}

func (filter *LeveldbStore) AddValue(key []byte, value []byte) error {
	log.Trace("add value key,", string(key), " : ", value)

	return filter.Leveldb.Put(key, value, nil)
}
