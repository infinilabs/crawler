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
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/core/util"
	"regexp"
	"time"
)

// JointConfig configs for each joint
type JointConfig struct {
	JointName  string                 `json:"joint" config:"joint"`                     //the joint name
	Parameters map[string]interface{} `json:"parameters,omitempty" config:"parameters"` //kv parameters for this joint
	Enabled    bool                   `json:"enabled" config:"enabled"`
}

// PipelineConfig config for each pipeline, a pipeline may have more than one joints
type PipelineConfig struct {
	ID            string         `json:"id,omitempty" index:"id"`
	Name          string         `json:"name,omitempty" config:"name"`
	StartJoint    *JointConfig   `json:"start,omitempty" config:"start"`
	ProcessJoints []*JointConfig `json:"process,omitempty" config:"process"`
	EndJoint      *JointConfig   `json:"end,omitempty" config:"end"`
	Created       *time.Time     `json:"created,omitempty"`
	Updated       *time.Time     `json:"updated,omitempty"`
	Tags          []string       `json:"tags,omitempty" config:"tags"`
}

const PipelineConfigBucket = "PipelineConfig"

func GetPipelineConfig(id string) (*PipelineConfig, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	b, err := persist.GetValue(PipelineConfigBucket, []byte(id))
	if err != nil {
		panic(err)
	}
	if len(b) > 0 {
		v := PipelineConfig{}
		err = json.Unmarshal(b, &v)
		return &v, err
	}
	return nil, errors.Errorf("not found, %s", id)
}

func CreatePipelineConfig(cfg *PipelineConfig) error {
	time := time.Now().UTC()
	cfg.ID = util.GetUUID()
	cfg.Created = &time
	cfg.Updated = &time
	b, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	return persist.AddValue(PipelineConfigBucket, []byte(cfg.ID), b)
}

func UpdatePipelineConfig(id string, cfg *PipelineConfig) error {
	time := time.Now().UTC()
	cfg.ID = id
	cfg.Updated = &time
	b, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	return persist.AddValue(PipelineConfigBucket, []byte(cfg.ID), b)
}

func DeletePipelineConfig(id string) error {
	return persist.DeleteKey(PipelineConfigBucket, []byte(id))
}

type HostConfig struct {
	ID         string `json:"id,omitempty" gorm:"not null;unique;primary_key" index:"id"`
	Host       string `gorm:"index" json:"host"`
	UrlPattern string `gorm:"index" json:"url_pattern"`
	Runner     string `gorm:"index" json:"runner"`
	SortOrder  int    `gorm:"index" json:"sort_order"`

	PipelineID string `gorm:"index" json:"pipeline_id"`
	Cookies    string `json:"cookies,omitempty"`

	Created time.Time `gorm:"index" json:"created,omitempty"`
	Updated time.Time `gorm:"index" json:"updated,omitempty"`
}

func CreateHostConfig(config *HostConfig) error {
	time := time.Now().UTC()
	config.ID = util.GetUUID()
	config.Created = time
	config.Updated = time
	err := persist.Save(config)
	if err != nil {
		panic(err)
	}

	return err
}

func UpdateHostConfig(config *HostConfig) {
	time := time.Now().UTC()
	config.Updated = time
	err := persist.Update(config)
	if err != nil {
		panic(err)
	}
}

func DeleteHostConfig(id string) error {
	config := HostConfig{ID: id}
	err := persist.Delete(&config)
	if err != nil {
		panic(err)
	}
	return err
}

func GetHostConfig(runner, host string) []HostConfig {
	var configs []HostConfig
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "sort_order", SortType: persist.DESC})
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

func GetHostConfigByHostAndUrl(runner, host, url string) *HostConfig {
	configs := GetHostConfig(runner, host)
	if len(configs) > 0 {
		for _, c := range configs {
			ok, err := regexp.Match(c.UrlPattern, []byte(url))
			if err != nil {
				log.Error(err)
				return nil
			}

			log.Debugf("match url:%v %v %v %v", host, url, c.UrlPattern, ok)
			if ok {
				return &c
			}
		}
	}
	return nil
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
