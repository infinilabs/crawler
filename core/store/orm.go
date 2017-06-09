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

package store

import (
	"github.com/jinzhu/gorm"
	"sync"
)

var dbLock sync.RWMutex

type ORM interface {
	Save(o interface{}) error

	Update(o interface{}) error

	Delete(o interface{}) error

	Search(o interface{}, q *Query) (error, Result)

	Get(key string, value interface{}, to interface{}) error

	Count(o interface{}) (int, error)
}

type Query struct {
	Sort   string
	From   int
	Size   int
	Filter *Cond
}

type Cond struct {
	Name  string
	Value interface{}
}

type Result struct {
	Total  int
	Result interface{}
}

func GetBy(field string, value interface{}, to interface{}) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	return getDBConnection().Where(field+" = ?", value).First(to).Error
}

func Get(o interface{}) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	return getDBConnection().First(o).Error
}

func Save(o interface{}) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	return getDBConnection().Save(o).Error
}

func Create(o interface{}) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	return getDBConnection().Create(o).Error
}

func Update(o interface{}) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	return getDBConnection().Save(o).Error
}

func Delete(o interface{}) error {
	dbLock.Lock()
	defer dbLock.Unlock()

	return getDBConnection().Delete(o).Error
}

func Count(o interface{}) (int, error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	var count int
	return count, getDBConnection().Model(&o).Count(count).Error
}

func Search(o interface{}, q *Query) (error, Result) {
	dbLock.Lock()
	defer dbLock.Unlock()

	if q.From < 0 {
		q.From = 0
	}
	if q.Size < 0 {
		q.Size = 10
	}

	var c int
	var err error
	db1 := getDBConnection().Model(o)
	if len(q.Sort) > 0 {
		db1 = db1.Order(q.Sort)
	}

	if q.Filter != nil {
		err = db1.Limit(q.Size).Offset(q.From).Where(q.Filter.Name+" = ?", q.Filter.Value).Find(o).Error
		db1.Where(q.Filter.Name+" = ?", q.Filter.Value).Count(&c)
	} else {
		err = db1.Limit(q.Size).Offset(q.From).Find(o).Error
		db1.Count(&c)
	}

	resut := Result{}
	resut.Result = o
	resut.Total = c
	return err, resut
}

var conn *gorm.DB

// create a session for each business unit of execution (e.g. a web request or goworkers job)
func getDBConnection() *gorm.DB {
	return conn
}

func RegisterDBConnection(h *gorm.DB) {
	conn = h
}
