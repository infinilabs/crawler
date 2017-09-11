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

package storage

import (
	"github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/index"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/modules/storage/boltdb"
	"github.com/infinitbyte/gopa/modules/storage/elastic"
	"os"
	"path"
)

var impl boltdb.BoltdbStore

func (this StorageModule) Name() string {
	return "Storage"
}

var storeConfig *StorageConfig

type BlotdbConfig struct {
}
type LeveldbConfig struct {
}

type StorageConfig struct {
	//Driver only `boltdb` and `elasticsearch` are available
	Driver  string                     `config:"driver"`
	Blotdb  *BlotdbConfig              `config:"blotdb"`
	Leveldb *LeveldbConfig             `config:"leveldb"`
	Elastic *index.ElasticsearchConfig `config:"elasticsearch"`
}

var (
	defaultConfig = StorageConfig{
		Driver:  "blotdb",
		Blotdb:  &BlotdbConfig{},
		Leveldb: &LeveldbConfig{},
		Elastic: &index.ElasticsearchConfig{
			Endpoint:    "http://localhost:9200",
			IndexPrefix: "gopa-",
		},
	}
)

func getDefaultConfig() StorageConfig {
	return defaultConfig
}

func (module StorageModule) Start(cfg *config.Config) {

	//init config
	config := getDefaultConfig()
	cfg.Unpack(&config)
	storeConfig = &config

	if config.Driver == "elasticsearch" {
		client := index.ElasticsearchClient{Config: config.Elastic}
		handler := elastic.ElasticsearchStore{&client}
		persist.RegisterKVHandler(handler)
	} else if config.Driver == "boltdb" {

		folder := path.Join(global.Env().SystemConfig.GetWorkingDir(), "blob")
		os.MkdirAll(folder, 0777)
		impl = boltdb.BoltdbStore{FileName: path.Join(folder, "/bolt.db")}
		err := impl.Open()
		if err != nil {
			panic(err)
		}
		persist.RegisterKVHandler(impl)
	} else {
		panic(errors.Errorf("invalid driver, %s", config.Driver))
	}
}

func (module StorageModule) Stop() error {
	if storeConfig.Driver == "blotdb" {
		return impl.Close()
	}
	return nil
}

type StorageModule struct {
}
