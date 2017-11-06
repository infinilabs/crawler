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

package persist

import (
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/index"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/modules/persist/elastic"
	"github.com/infinitbyte/gopa/modules/persist/mysql"
	"github.com/infinitbyte/gopa/modules/persist/sqlite"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func (module DatabaseModule) Name() string {
	return "Persist"
}

var (
	defaultConfig = PersistConfig{
		Driver: "sqlite",
		SQLite: &sqlite.SQLiteConfig{},
		MySQL:  &mysql.MySQLConfig{},
		Elastic: &index.ElasticsearchConfig{
			Endpoint:    "http://localhost:9200",
			IndexPrefix: "gopa-",
		},
	}
)

func getDefaultConfig() PersistConfig {
	return defaultConfig
}

var db *gorm.DB

type PersistConfig struct {
	//Driver only `mysql` and `sqlite` are available
	Driver  string                     `config:"driver"`
	SQLite  *sqlite.SQLiteConfig       `config:"sqlite"`
	MySQL   *mysql.MySQLConfig         `config:"mysql"`
	Elastic *index.ElasticsearchConfig `config:"elasticsearch"`
}

func (module DatabaseModule) Start(cfg *Config) {

	//init config
	config := getDefaultConfig()
	cfg.Unpack(&config)
	module.config = &config

	if config.Driver == "elasticsearch" {
		client := index.ElasticsearchClient{Config: config.Elastic}
		handler := elastic.ElasticORM{&client}
		persist.Register(handler)
		return
	}

	//whether use lock, only sqlite need lock
	userLock := false
	if config.Driver == "sqlite" {
		db = sqlite.GetInstance(config.SQLite)
		userLock = true
	} else if config.Driver == "mysql" {
		db = mysql.GetInstance(config.MySQL)
	} else {
		panic(errors.Errorf("invalid driver, %s", config.Driver))
	}

	// Migrate the schema
	db.AutoMigrate(&model.Host{})
	db.AutoMigrate(&model.Task{})
	db.AutoMigrate(&model.Snapshot{})

	handler := SQLORM{conn: db, useLock: userLock}

	persist.Register(handler)
}

func (module DatabaseModule) Stop() error {
	if db != nil {
		db.Close()
	}
	return nil

}

type DatabaseModule struct {
	config *PersistConfig
}
