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

package model

import (
	"github.com/infinitbyte/framework/core/errors"
	"github.com/infinitbyte/framework/core/persist"
	"github.com/infinitbyte/framework/core/util"
	"regexp"
	"time"
)

type HostConfig struct {
	ID         string `json:"id,omitempty" index:"id"`
	Host       string `json:"host"`
	UrlPattern string `json:"url_pattern"`
	Runner     string `json:"runner"`
	SortOrder  int    `json:"sort_order"`

	PipelineID string `json:"pipeline_id"`
	Cookies    string `json:"cookies,omitempty"`

	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

func CreateHostConfig(config *HostConfig) error {
	time := time.Now().UTC()
	config.ID = util.GetUUID()
	config.Created = time
	config.Updated = time
	return persist.Save(config)
}

func UpdateHostConfig(config *HostConfig) error {
	time := time.Now().UTC()
	config.Updated = time
	return persist.Update(config)
}

func DeleteHostConfig(id string) error {
	config := HostConfig{ID: id}
	return persist.Delete(&config)
}

func GetHostConfigByID(id string) (HostConfig, error) {
	o := HostConfig{ID: id}
	err := persist.Get(&o)
	return o, err
}

func GetHostConfigList(from, size int, host string) (int, []HostConfig, error) {
	var configs []HostConfig

	query := persist.Query{From: from, Size: size}
	if len(host) > 0 {
		query.Conds = persist.And(persist.Eq("host", host))
	}

	err, result := persist.Search(HostConfig{}, &configs, &query)

	if result.Result != nil && configs == nil || len(configs) == 0 {
		convertHostConfig(result, &configs)
	}

	return result.Total, configs, err
}

func GetHostConfig(runner, host string) []HostConfig {
	var configs []HostConfig
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "sort_order", SortType: persist.ASC})
	queryO := persist.Query{Sort: &sort, From: 0, Size: 100}
	if len(host) > 0 {
		if runner != "" {
			queryO.Conds = persist.And(persist.Eq("host", host), persist.Eq("runner", runner))
		} else {
			queryO.Conds = persist.And(persist.Eq("host", host))
		}
	}
	err, result := persist.Search(HostConfig{}, &configs, &queryO)
	if err != nil {
		panic(err)
	}

	if result.Result != nil && configs == nil || len(configs) == 0 {
		convertHostConfig(result, &configs)
	}

	return configs
}

func GetHostConfigByHostAndUrl(runner, host, url string) (*HostConfig, error) {
	configs := GetHostConfig(runner, host)
	if len(configs) > 0 {
		for _, c := range configs {
			ok, err := regexp.Match(c.UrlPattern, []byte(url))
			if err != nil {
				return nil, err
			}

			if ok {
				return &c, nil
			}
		}
	}
	return nil, errors.New("not found")
}

func convertHostConfig(result persist.Result, configs *[]HostConfig) {
	if result.Result == nil {
		return
	}

	t, ok := result.Result.([]interface{})
	if ok {
		for _, i := range t {
			js := util.ToJson(i, false)
			t := HostConfig{}
			util.FromJson(js, &t)
			*configs = append(*configs, t)
		}
	}
}
