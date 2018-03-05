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

package impl

import (
	log "github.com/cihub/seelog"
	core "github.com/infinitbyte/gopa/core/filter"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/util"
	. "github.com/zeebo/sbloom"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type BloomFilter struct {
	PersistFileName string
	filter          *Filter
	ProbItems       int
}

var filters map[string]core.Filter

func (filter *BloomFilter) Open() error {

	//loading or initializing bloom filter
	if util.FileExists(filter.PersistFileName) {
		log.Debug("found bloomFilter,start reload,", filter.PersistFileName)
		n, err := ioutil.ReadFile(filter.PersistFileName)
		if err != nil {
			log.Error("bloomFilter:", filter.PersistFileName, err)
		}
		filter.filter = new(Filter)
		if err := filter.filter.GobDecode(n); err != nil {
			log.Error("bloomFilter:", filter.PersistFileName, err)
		}
		log.Info("bloomFilter successfully reloaded:", filter.PersistFileName)
	} else {

		probItems := 1000000 //config.GetIntConfig("BloomFilter", "ItemSize", 100000)
		log.Debug("initializing bloom-filter", filter.PersistFileName, ",virual size is,", probItems)
		filter.filter = NewFilter(fnv.New64(), probItems)
		log.Info("bloomFilter successfully initialized:", filter.PersistFileName)
	}

	return nil
}

func (filter *BloomFilter) Close() error {

	log.Debug("bloomFilter start persist,file:", filter.PersistFileName)

	//save bloom-filter
	m, err := filter.filter.GobEncode()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filter.PersistFileName, m, 0600)
	if err != nil {
		panic(err)
	}
	log.Info("bloomFilter safety persisted.")

	return nil
}

func (filter *BloomFilter) Exists(bucket string, key []byte) bool {
	return filter.filter.Lookup(key)
}

func (filter *BloomFilter) Add(bucket string, key []byte) error {
	filter.filter.Add(key)
	return nil
}

func (filter *BloomFilter) Delete(bucket string, key []byte) error {

	return nil
}

var l sync.RWMutex

func (filter *BloomFilter) CheckThenAdd(bucket string, key []byte) (b bool, err error) {
	f := filters[bucket]
	l.Lock()
	defer l.Unlock()
	b = f.Exists(bucket, key)
	if !b {
		err = f.Add(bucket, key)
	}
	return b, err
}

func initBloomFilter(key string) {
	//f := impl.PersistFilter{FilterBucket: string(key)}
	dir := path.Join(global.Env().SystemConfig.GetWorkingDir(), "filters")
	os.MkdirAll(dir, 0777)
	file := path.Join(dir, string(key))
	f := BloomFilter{PersistFileName: file}
	f.Open()
	filters[key] = f
}
