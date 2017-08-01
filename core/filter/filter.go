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

package filter

// Key is the key of filters
type Key string

// Filter is used to check if the object is in the filter or not
type Filter interface {
	Exists(bucket Key, key []byte) bool
	Add(bucket Key, key []byte) error
	Delete(bucket Key, key []byte) error
	CheckThenAdd(bucket Key, key []byte) (bool, error)
}

var handler Filter

// Exists checks if the key are already in filter bucket
func Exists(bucket Key, key []byte) bool {
	return handler.Exists(bucket, key)
}

// Add will add key to filter bucket
func Add(bucket Key, key []byte) error {
	return handler.Add(bucket, key)
}

// Remove will remove key from bucket
func Remove(bucket Key, key []byte) error {
	return handler.Delete(bucket, key)
}

// CheckThenAdd will check first and if the key is not in the filter bucket, then it will add it and return false, if the key is already in the bucket, it will just return true
func CheckThenAdd(bucket Key, key []byte) (bool, error) {
	return handler.CheckThenAdd(bucket, key)
}

// Register used to register filter handler to dealing with filter operations
func Register(h Filter) {
	handler = h
}
