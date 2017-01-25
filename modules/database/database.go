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
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/core/store"
	"os"
	"path"
)

func (this DatabaseModule) Name() string {
	return "Database"
}

var db *gorm.DB

func (this DatabaseModule) Start(env *Env) {
	os.MkdirAll(path.Join(global.Env().SystemConfig.GetDataDir(), "database/"), 0777)
	fileName := fmt.Sprintf("file:%s?cache=shared&mode=rwc&_busy_timeout=50000000", path.Join(global.Env().SystemConfig.GetDataDir(), "database/db.sqlite"))

	var err error
	db, err = gorm.Open("sqlite3", fileName)
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&model.Domain{})
	db.AutoMigrate(&model.Task{})

	store.RegisterConnection(db)
}

func (this DatabaseModule) Stop() error {
	db.Close()
	return nil

}

type DatabaseModule struct {
}
