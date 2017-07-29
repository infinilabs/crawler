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

package database

import (
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/store"
	"github.com/infinitbyte/gopa/modules/database/mysql"
	"github.com/infinitbyte/gopa/modules/database/sqlite"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func (this DatabaseModule) Name() string {
	return "Database"
}

var (
	defaultConfig = DatabaseConfig{
		Dialect: "sqlite",
		SQLite:  &sqlite.SQLiteConfig{},
		MySQL:   &mysql.MySQLConfig{},
	}
)

func GetDefaultConfig() DatabaseConfig {
	return defaultConfig
}

var db *gorm.DB

type DatabaseConfig struct {
	Dialect string               `config:"dialect"` //only `mysql` and `sqlite` are available
	SQLite  *sqlite.SQLiteConfig `config:"sqlite"`
	MySQL   *mysql.MySQLConfig   `config:"mysql"`
}

func (this DatabaseModule) Start(cfg *Config) {

	//init config
	config := GetDefaultConfig()
	cfg.Unpack(&config)
	this.config = &config

	if config.Dialect == "sqlite" {
		db = sqlite.GetInstance(config.SQLite)
	} else if config.Dialect == "mysql" {
		db = mysql.GetInstance(config.MySQL)
	} else {
		panic(errors.New("database is not successful started, invalid type"))
	}

	// Migrate the schema
	db.AutoMigrate(&model.Domain{})
	db.AutoMigrate(&model.Task{})
	db.AutoMigrate(&model.Snapshot{})

	store.RegisterDBConnection(db)
}

func (this DatabaseModule) Stop() error {
	db.Close()
	return nil

}

type DatabaseModule struct {
	config *DatabaseConfig
}
