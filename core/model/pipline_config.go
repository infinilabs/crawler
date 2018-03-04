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
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/core/util"
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
	StartJoint    *JointConfig   `gorm:"-" json:"start,omitempty" config:"start"`
	ProcessJoints []*JointConfig `gorm:"-" json:"process,omitempty" config:"process"`
	EndJoint      *JointConfig   `gorm:"-" json:"end,omitempty" config:"end"`
	Created       int64          `json:"created,omitempty"`
	Updated       int64          `json:"updated,omitempty"`
	Tags          []string       `gorm:"-" json:"tags,omitempty" config:"tags"`
}

const PipelineConfigBucket = "PipelineConfig"

func GetPipelineConfig(id string) (*PipelineConfig, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	b, err := persist.GetValue(PipelineConfigBucket, []byte(id))
	if err != nil {
		return nil, err
	}
	if len(b) > 0 {
		v := PipelineConfig{}
		err = json.Unmarshal(b, &v)
		return &v, err
	}
	return nil, errors.Errorf("not found, %s", id)
}

func GetPipelineList(from, size int) (int, []PipelineConfig, error) {
	var configs []PipelineConfig

	query := persist.Query{From: from, Size: size}

	err, r := persist.Search(PipelineConfig{}, &configs, &query)
	if r.Result != nil && configs == nil || len(configs) == 0 {
		convertPipeline(r, &configs)
	}
	return r.Total, configs, err
}

func CreatePipelineConfig(cfg *PipelineConfig) error {
	time := time.Now().UTC().Unix()
	cfg.ID = util.GetUUID()
	cfg.Created = time
	cfg.Updated = time
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = persist.AddValue(PipelineConfigBucket, []byte(cfg.ID), b)
	if err != nil {
		return err
	}
	return persist.Save(cfg)
}

func UpdatePipelineConfig(id string, cfg *PipelineConfig) error {
	time := time.Now().UTC().Unix()
	cfg.ID = id
	cfg.Updated = time
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = persist.AddValue(PipelineConfigBucket, []byte(cfg.ID), b)
	if err != nil {
		return err
	}
	return persist.Update(cfg)
}

func DeletePipelineConfig(id string) error {
	err := persist.DeleteKey(PipelineConfigBucket, []byte(id))
	if err != nil {
		return err
	}
	o := PipelineConfig{ID: id}
	return persist.Delete(&o)
}

func convertPipeline(result persist.Result, pipelines *[]PipelineConfig) {
	if result.Result == nil {
		return
	}

	t, ok := result.Result.([]interface{})
	if ok {
		for _, i := range t {
			js := util.ToJson(i, false)
			t := PipelineConfig{}
			util.FromJson(js, &t)
			*pipelines = append(*pipelines, t)
		}
	}
}
