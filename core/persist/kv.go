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

import "errors"

type KVStore interface {
	Open() error

	Close() error

	GetValue(bucket string, key []byte) []byte

	GetCompressedValue(bucket string, key []byte) []byte

	AddValueCompress(bucket string, key []byte, value []byte) error

	AddValue(bucket string, key []byte, value []byte) error

	DeleteValue(bucket string, key []byte, value []byte) error

	DeleteBucket(bucket string, key []byte, value []byte) error
}

var kvHandler KVStore

func getKVHandler() KVStore {

	if kvHandler == nil {
		panic(errors.New("kv store handler is not registered"))
	}
	return kvHandler
}

func GetValue(bucket string, key []byte) []byte {
	return getKVHandler().GetValue(bucket, key)
}

func GetCompressedValue(bucket string, key []byte) []byte {
	return getKVHandler().GetCompressedValue(bucket, key)
}

func AddValueCompress(bucket string, key []byte, value []byte) error {
	return getKVHandler().AddValueCompress(bucket, key, value)
}

func AddValue(bucket string, key []byte, value []byte) error {
	return getKVHandler().AddValue(bucket, key, value)
}

func DeleteValue(bucket string, key []byte, value []byte) error {
	return getKVHandler().DeleteValue(bucket, key, value)
}

func DeleteBucket(bucket string, key []byte, value []byte) error {
	return getKVHandler().DeleteBucket(bucket, key, value)
}

func RegisterKVHandler(h KVStore) {
	kvHandler = h
}
