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

package types


type Store interface {
	//Store(url string, data []byte)
	//Get(key string) []byte
	//List(from int, size int) [][]byte
	//TaskEnqueue([]byte)

	Open() error
	Close() error

	//UrlHasWalked(url []byte) bool
	//UrlHasFetched(url []byte) bool
	//FileHasParsed(url []byte) bool
	//AddWalkedUrl(url []byte )
	//AddFetchedUrl(url []byte)
	//AddSavedUrl(url []byte )   //the file already saved,but is missing in bloom filter,run this method
	//AddParsedFile(url []byte )

	//LogSavedFile(path,content string )

	//LogFetchFailedUrl(path,content string )

	//FileHasSaved(file string)  bool

	//InitPendingFetchBloomFilter(fileName string)
	//PendingFetchUrlHasAdded(url []byte) bool
	//AddPendingFetchUrl(url []byte )
	//LogPendingFetchUrl(path,content string )

	//LoadOffset(fileName string) int64
	//PersistOffset(fileName string,offset int64)


	GetValue(bucket string, key []byte) []byte

	AddValue(bucket string, key []byte, value []byte) error

	DeleteValue(bucket string, key []byte, value []byte) error

	DeleteBucket(bucket string, key []byte, value []byte) error

}

