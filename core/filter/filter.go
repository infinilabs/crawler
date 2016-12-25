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

type FilterKey string


type Filter interface {
	Exists(bucket FilterKey, key []byte) bool
	Add(bucket FilterKey, key []byte) error
	CheckThenAdd(bucket FilterKey,key[]byte)(bool,error)
}

var handler Filter

func Exists(bucket FilterKey, key []byte) bool {
	return handler.Exists(bucket, key)
}

func Add(bucket FilterKey, key []byte) error {
	return handler.Add(bucket, key)
}

func CheckThenAdd(bucket FilterKey,key[]byte)(bool,error){
	return handler.CheckThenAdd(bucket,key)
}

func Regsiter(h Filter) {
	handler = h
}
