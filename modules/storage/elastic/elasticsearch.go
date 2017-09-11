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

package elastic

import (
	"encoding/base64"
	lz4 "github.com/bkaradzic/go-lz4"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/index"
)

type ElasticsearchStore struct {
	Client *index.ElasticsearchClient
}

func (store ElasticsearchStore) Open() error {
	//TODO check index and mapping
	return nil
}

func (store ElasticsearchStore) Close() error {
	return nil
}

func (store ElasticsearchStore) GetCompressedValue(bucket string, key []byte) ([]byte, error) {

	data, err := store.GetValue(bucket, key)
	if err != nil {
		return nil, err
	}
	data, err = lz4.Decode(nil, data)
	if err != nil {
		log.Error("Failed to decode:", err)
		return nil, err
	}
	return data, nil
}

func (store ElasticsearchStore) GetValue(bucket string, key []byte) ([]byte, error) {
	response, err := store.Client.Get(indexName, string(key))
	if err != nil {
		return nil, err
	}
	if response.Found {
		content := response.Source["content"]
		uDec, err := base64.URLEncoding.DecodeString(content.(string))
		if err != nil {
			return nil, err
		}
		return uDec, nil
	}
	return nil, errors.New("not found")
}

var indexName = "blob"

type Blob struct {
	Content string `json:"content,omitempty"`
}

func (store ElasticsearchStore) AddValueCompress(bucket string, key []byte, value []byte) error {
	value, err := lz4.Encode(nil, value)
	if err != nil {
		log.Error("Failed to encode:", err)
		return err
	}
	return store.AddValue(bucket, key, value)
}

func (store ElasticsearchStore) AddValue(bucket string, key []byte, value []byte) error {
	file := Blob{}
	file.Content = base64.URLEncoding.EncodeToString(value)
	_, err := store.Client.Index(indexName, string(key), file)
	return err
}

func (store ElasticsearchStore) DeleteValue(bucket string, key []byte, value []byte) error {
	_, err := store.Client.Delete(indexName, string(key))
	return err
}

func (store ElasticsearchStore) DeleteBucket(bucket string, key []byte) error {
	_, err := store.Client.Delete(indexName, string(key))
	return err
}
