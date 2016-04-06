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

package fs

import (
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/config"
	util "github.com/medcl/gopa/core/util"
)

type FsStore struct {
	WalkBloomFilterFileName         string
	FetchBloomFilterFileName        string
	ParseBloomFilterFileName        string
	PendingFetchBloomFilterFileName string
	WalkBloomFilter                 util.DeduplicatePlugin
	FetchBloomFilter                util.DeduplicatePlugin
	ParseBloomFilter                util.DeduplicatePlugin
	PendingFetchBloomFilter         util.DeduplicatePlugin
}

func (this *FsStore) Store(url string, data []byte) {
	util.FilePutContentWithByte(url, data)
}

func (this *FsStore) Get(key string) []byte {
	file, error := util.FileGetContent(key)
	if error != nil {
		log.Error("get file:", key, error)
	}
	return file
}

func (this *FsStore) List(from int, size int) [][]byte {
	return nil
}

func (this *FsStore) TaskEnqueue(url []byte) {
	log.Info("task enqueue:", string(url))
}

func (this *FsStore) Open() error {

	var runtimeConfig = config.InitOrGetConfig()
	this.FetchBloomFilterFileName = runtimeConfig.FetchBloomFilterFileName
	this.WalkBloomFilterFileName = runtimeConfig.WalkBloomFilterFileName
	this.ParseBloomFilterFileName = runtimeConfig.ParseBloomFilterFileName
	this.PendingFetchBloomFilterFileName = runtimeConfig.PendingFetchBloomFilterFileName

	this.WalkBloomFilter = new(BloomFilter)
	this.WalkBloomFilter.Init(this.WalkBloomFilterFileName)

	this.FetchBloomFilter = new(BloomFilter)
	this.FetchBloomFilter.Init(this.FetchBloomFilterFileName)

	this.ParseBloomFilter = new(BloomFilter)
	this.ParseBloomFilter.Init(this.ParseBloomFilterFileName)

	this.PendingFetchBloomFilter = new(BloomFilter)
	this.PendingFetchBloomFilter.Init(this.PendingFetchBloomFilterFileName)

	return nil
}

func (this *FsStore) Close() error

func (this *FsStore) PersistBloomFilter() {
	this.WalkBloomFilter.Persist()
	this.FetchBloomFilter.Persist()
	this.ParseBloomFilter.Persist()
	this.PendingFetchBloomFilter.Persist()
}

func (this *FsStore) CheckWalkedUrl(url []byte) bool {
	return this.WalkBloomFilter.Lookup(url)
}

func (this *FsStore) CheckFetchedUrl(url []byte) bool {
	return this.FetchBloomFilter.Lookup(url)
}

func (this *FsStore) CheckParsedFile(url []byte) bool {
	return this.ParseBloomFilter.Lookup(url)
}

func (this *FsStore) CheckPendingFetchUrl(url []byte) bool {
	return this.PendingFetchBloomFilter.Lookup(url)
}

func (this *FsStore) AddWalkedUrl(url []byte) {
	this.WalkBloomFilter.Add(url)
}

func (this *FsStore) AddPendingFetchUrl(url []byte) {
	this.PendingFetchBloomFilter.Add(url)
}

func (this *FsStore) AddSavedUrl(url []byte) {
	this.WalkBloomFilter.Add(url)
	this.FetchBloomFilter.Add(url)
}

func (this *FsStore) LogSavedFile(path string, content string) {
	util.FileAppendNewLine(path, content)
}

func (this *FsStore) LogPendingFetchUrl(path string, content string) {
	util.FileAppendNewLine(path, content)
}

func (this *FsStore) LogFetchFailedUrl(path string, content string) {
	util.FileAppendNewLine(path, content)
}

func (this *FsStore) AddFetchedUrl(url []byte) {
	this.FetchBloomFilter.Add(url)
}

func (this *FsStore) saveFetchedUrlToLocalFile(path string, url string) {
	util.FileAppendNewLine(path, url)
}

func (this *FsStore) AddParsedFile(url []byte) {
	this.ParseBloomFilter.Add(url)
}

func (this *FsStore) AddFetchFailedUrl(url []byte) {
	//TODO
	log.Error("fetch failed url:", string(url))
}

func (this *FsStore) CheckSavedFile(file string) bool {
	log.Debug("start check file:", file)
	return util.CheckFileExists(file)
}

func (this *FsStore) LoadOffset(fileName string) int64 {
	log.Debug("start init offsets,", fileName)
	if util.CheckFileExists(fileName) {
		log.Debug("found offset file,start loading,", fileName)
		n, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Error("offset", fileName, ",", err)
			return 0
		}
		ret, err := strconv.ParseInt(string(n), 10, 64)
		if err != nil {
			log.Error("offset", fileName, ",", err)
			return 0
		}
		log.Info("init offsets successfully,", fileName, ":", ret)
		return int64(ret)
	}

	return 0
}

func (this *FsStore) PersistOffset(fileName string, offset int64) {
	//persist worker's offset
	path := fileName + ".tmp"
	fout, error := os.Create(path)
	if error != nil {
		log.Error(path, error)
		return
	}

	defer fout.Close()
	log.Debug("saved offset:", fileName, ":", offset)
	fout.Write([]byte(strconv.FormatInt(offset, 10)))
	util.CopyFile(path, fileName)
}

func (this *FsStore) InitPendingFetchBloomFilter(fileName string) {}
