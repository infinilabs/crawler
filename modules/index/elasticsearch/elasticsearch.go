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

package util

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/util"
)

// ElasticsearchConfig contains common settings for elasticsearch
type ElasticsearchConfig struct {
	Endpoint string `config:"elasticsearch.endpoint"`
	Index    string `config:"elasticsearch.index"`
}

// ElasticsearchClient elasticsearch client api
type ElasticsearchClient struct {
	Endpoint string
	Index    string
}

// InsertResponse is a response object
type InsertResponse struct {
	Created bool   `json:"created"`
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	ID      string `json:"_id"`
	Version int    `json:"_version"`
}

// IndexDoc index a document into elasticsearch
func (c *ElasticsearchClient) IndexDoc(typeName, id string, data map[string]interface{}) (*InsertResponse, error) {

	url := c.Endpoint + "/" + c.Index + "/" + typeName + "/" + id

	js, err := json.Marshal(data)

	log.Debug("indexing doc: ", url, ",", string(js))

	if err != nil {
		return nil, err
	}
	response := util.HttpPostJSON(url, "", string(js))
	if err != nil {
		return nil, err
	}

	log.Debug("indexing response: ", string(response))

	esResp := &InsertResponse{}
	err = json.Unmarshal(response, esResp)
	if err != nil {
		return &InsertResponse{}, err
	}

	return esResp, nil
}
