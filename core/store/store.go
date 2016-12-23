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

import "errors"

type Store interface {
	Open() error

	Close() error

	GetValue(bucket string, key []byte) []byte

	GetCompressedValue(bucket string, key []byte) []byte

	AddValueCompress(bucket string, key []byte, value []byte) error

	AddValue(bucket string, key []byte, value []byte) error

	DeleteValue(bucket string, key []byte, value []byte) error

	DeleteBucket(bucket string, key []byte, value []byte) error
}

type ORM interface {
	Save(o interface{}) error

	Update(o interface{}) error

	Delete(o interface{}) error

	Search(t1,t2 interface{}, q *Query) (error, Result)

	Get(key string, value interface{}, to interface{}) error

	Count(o interface{}) (int, error)
}

type Query struct {
	Sort string
	From int
	Size int
}

type Result struct {
	Total  int
	Result interface{}
}

var handler Store
var theORMHandler ORM

func GetValue(bucket string, key []byte) []byte {
	return getHandler().GetValue(bucket, key)
}

func GetCompressedValue(bucket string, key []byte) []byte {
	return getHandler().GetCompressedValue(bucket, key)
}

func AddValueCompress(bucket string, key []byte, value []byte) error {
	return getHandler().AddValueCompress(bucket, key, value)
}

func AddValue(bucket string, key []byte, value []byte) error {
	return getHandler().AddValue(bucket, key, value)
}

func DeleteValue(bucket string, key []byte, value []byte) error {
	return getHandler().DeleteValue(bucket, key, value)
}

func DeleteBucket(bucket string, key []byte, value []byte) error {
	return getHandler().DeleteBucket(bucket, key, value)
}

func Get(key string, value interface{}, to interface{}) error {
	return getORMHandler().Get(key, value, to)
}

func Save(o interface{}) error {
	return getORMHandler().Save(o)
}

func Update(o interface{}) error {
	return getORMHandler().Update(o)
}

func Delete(o interface{}) error {
	return getORMHandler().Delete(o)
}

func Count(o interface{}) (int, error) {
	return getORMHandler().Count(o)
}

func Search(t1,t2 interface{}, q *Query) (error, Result) {
	return getORMHandler().Search(t1,t2, q)
}

func getHandler() Store {
	if handler == nil {
		panic(errors.New("store handler is not registered"))
	}
	return handler
}

func getORMHandler() ORM {
	if theORMHandler == nil {
		panic(errors.New("ORM handler is not registered"))
	}
	return theORMHandler
}

func RegisterStoreHandler(h Store) {
	handler = h
}

func RegisterORMHandler(h ORM) {
	theORMHandler = h
}
