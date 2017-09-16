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
func (store *LeveldbStore) Store(url string, data []byte) {
	util.FilePutContentWithByte(url, data)
}

//get webpage file
func (store *LeveldbStore) Get(key string) []byte {
	file, error := util.FileGetContent(key)
	if error != nil {
		log.Error("get file:", key, error)
	}
	return file
}

func (store *LeveldbStore) List(from int, size int) [][]byte {
	return nil
}

func (store *LeveldbStore) TaskEnqueue(url []byte) {
	log.Info("task enqueue:", string(url))
}

func (store *LeveldbStore) Open() error {

	//var runtimeConfig= config.InitOrGetConfig()
	store.WalkPrefix = "walk"
	store.FetchPrefix = "fetch"
	store.ParsePrefix = "parse"
	store.PendingFetchPrefix = "pfetch"
	store.OffsetPrefix = "offset"

	store.PersistFileName = "leveldb"

	//loading or initializing leveldb
	log.Debug("found leveldb file, start reload,", store.PersistFileName)

	db, err := leveldb.OpenFile(store.PersistFileName, nil)
	store.Leveldb = db

	if err != nil {
		log.Error("leveldb:", store.PersistFileName, err)
		return err
	}

	log.Info("leveldb successfully reloaded:", store.PersistFileName)

	return nil
}

func (store *LeveldbStore) Close() error {
	err := store.Leveldb.Close()
	if err != nil {
		log.Error("leveldb:", store.PersistFileName, err)
	}
	return err
}

func (store *LeveldbStore) PersistBloomFilter() {

}

func (store *LeveldbStore) UrlHasWalked(url []byte) bool {

	c := [][]byte{[]byte(store.WalkPrefix), url}

	return store.Lookup(bytes.Join(c, []byte(":")))
}

func (store *LeveldbStore) UrlHasFetched(url []byte) bool {
	c := [][]byte{[]byte(store.FetchPrefix), url}

	return store.Lookup(bytes.Join(c, []byte(":")))
}

func (store *LeveldbStore) FileHasParsed(url []byte) bool {
	c := [][]byte{[]byte(store.ParsePrefix), url}

	return store.Lookup(bytes.Join(c, []byte(":")))
}

func (store *LeveldbStore) PendingFetchUrlHasAdded(url []byte) bool {
	c := [][]byte{[]byte(store.PendingFetchPrefix), url}

	key := bytes.Join(c, []byte(":"))
	result := store.Lookup(key)

	value := store.GetValue(key)

	log.Trace("check pending url error,", string(key), ",", result, ",value:", string(value))
	return result
}

func (store *LeveldbStore) AddWalkedUrl(url []byte) {
	c := [][]byte{[]byte(store.WalkPrefix), url}

	store.Add(bytes.Join(c, []byte(":")))
}

func (store *LeveldbStore) AddPendingFetchUrl(url []byte) {
	c := [][]byte{[]byte(store.PendingFetchPrefix), url}

	key := bytes.Join(c, []byte(":"))
	err := store.Add(key)

	if err != nil {
		log.Error("add pending url error,", url, ",", err)
	}
}

func (store *LeveldbStore) AddSavedUrl(url []byte) {
	store.AddWalkedUrl(url)
	store.AddFetchedUrl(url)
}

func (store *LeveldbStore) LogSavedFile(path string, content string) {
	util.FileAppendNewLine(path, content)
}

func (store *LeveldbStore) LogPendingFetchUrl(path string, content string) {
	util.FileAppendNewLine(path, content)
}

func (store *LeveldbStore) LogFetchFailedUrl(path string, content string) {
	store.AddFetchFailedUrl([]byte(path))
	util.FileAppendNewLine(path, content)
}

func (store *LeveldbStore) AddFetchedUrl(url []byte) {
	c := [][]byte{[]byte(store.FetchPrefix), url}

	store.Add(bytes.Join(c, []byte(":")))

}

func (store *LeveldbStore) saveFetchedUrlToLocalFile(path string, url string) {
	util.FileAppendNewLine(path, url)
}

func (store *LeveldbStore) AddParsedFile(url []byte) {
	c := [][]byte{[]byte(store.ParsePrefix), url}

	store.Add(bytes.Join(c, []byte(":")))
}

func (store *LeveldbStore) AddFetchFailedUrl(url []byte) {
	c := [][]byte{[]byte(store.WalkPrefix), url}

	store.Add(bytes.Join(c, []byte(":")))

	log.Debug("fetch failed url:", string(url))
}

func (store *LeveldbStore) FileHasSaved(file string) bool {
	log.Debug("start check file:", file)
	return util.FileExists(file)
}

func (store *LeveldbStore) LoadOffset(fileName string) int64 {
	log.Debug("start load offsets,", fileName)

	c := [][]byte{([]byte)(store.OffsetPrefix), []byte(fileName)}

	n := store.GetValue(bytes.Join(c, []byte(":")))

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

func (store *LeveldbStore) PersistOffset(fileName string, offset int64) {
	//persist worker's offset

	c := [][]byte{[]byte(store.OffsetPrefix), []byte(fileName)}

	error := store.AddValue(bytes.Join(c, []byte(":")), []byte(strconv.FormatInt(offset, 10)))

	if error != nil {
		log.Error(fileName, error)
		return
	}
}

func (store *LeveldbStore) InitPendingFetchBloomFilter(fileName string) {}

//TODO REMOVE
func (store *LeveldbStore) Persist() error {

	log.Debug("leveldb start persist,file:", store.PersistFileName)

	log.Info("leveldb safety persisted.")

	return nil
}

func (store *LeveldbStore) Lookup(key []byte) bool {
	value := store.GetValue(key)

	if value != nil {
		log.Trace("return true,hit key, ", string(key), " : ", value)
		return true
	}
	log.Trace("return false, hit key", string(key), " : ", value)

	return false
}

func (store *LeveldbStore) Add(key []byte) error {
	log.Trace("add key,", string(key))

	return store.Leveldb.Put(key, []byte("true"), nil)
}

func (store *LeveldbStore) GetValue(key []byte) []byte {

	value, err := store.Leveldb.Get(key, nil)
	if err != nil {
		log.Trace("leveldb getValue error, ", err, " , ", string(key), " : ", value)
		return value
	}

	log.Trace("get key, ", err, " , ", string(key), " : ", value)

	return value
}

func (store *LeveldbStore) AddValue(key []byte, value []byte) error {
	log.Trace("add value key,", string(key), " : ", value)

	return store.Leveldb.Put(key, value, nil)
}
