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
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/store"
	"time"
)

type KV struct {
	Key   string   `json:"key,omitempty"`
	Value []string `storm:"inline" json:"value,omitempty"`
}

type LinkGroup struct {
	Internal []PageLink `json:"internal,omitempty"`
	External []PageLink `json:"external,omitempty"`
}

type Snapshot struct {
	ID      string `json:"id,omitempty" gorm:"not null;unique;primary_key"`
	Version int    `json:"version,omitempty"`
	Url     string `json:"url,omitempty"`
	TaskID  string `json:"task_id,omitempty"`
	Path    string `json:"path,omitempty"  gorm:"-"` //path of this file
	File    string `json:"file,omitempty"  gorm:"-"` //filename of this page

	StatusCode int    `json:"-" gorm:"-"`
	Payload    []byte `json:"-" gorm:"-"`
	Size       uint64 `json:"size,omitempty"`

	Headers    map[string][]string     `json:"-" gorm:"-"`
	Metadata   *map[string]interface{} `json:"-" gorm:"-"`
	Parameters []KV                    `json:"-" gorm:"-"`

	Language string `json:"lang,omitempty" gorm:"-"`

	Title       string `json:"title,omitempty"`
	Summary     string `json:"summary,omitempty" gorm:"-"`
	Text        string `json:"text,omitempty" gorm:"-"`
	ContentType string `json:"content_type,omitempty"`

	Tags []string `json:"tags,omitempty" gorm:"-"`

	Links LinkGroup `json:"links,omitempty" gorm:"-"`

	Images struct {
		Internal []PageLink `json:"internal,omitempty"`
		External []PageLink `json:"external,omitempty"`
	} `json:"images,omitempty" gorm:"-"`

	H1     []string `json:"h1,omitempty" gorm:"-"`
	H2     []string `json:"h2,omitempty" gorm:"-"`
	H3     []string `json:"h3,omitempty" gorm:"-"`
	H4     []string `json:"h4,omitempty" gorm:"-"`
	H5     []string `json:"h5,omitempty" gorm:"-"`
	Bold   []string `json:"bold,omitempty" gorm:"-"`
	Italic []string `json:"italic,omitempty" gorm:"-"`

	Classifications  []string                `json:"classifications,omitempty" gorm:"-"`
	EnrichedFeatures *map[string]interface{} `json:"enriched_features,omitempty" gorm:"-"`

	Hash    string `json:"hash,omitempty"`
	SimHash string `json:"sim_hash,omitempty"`

	CreateTime *time.Time `json:"created,omitempty"`
}

type PageLink struct {
	Url   string `json:"url"`
	Label string `json:"label"`
}

func CreateSnapshot(snapshot *Snapshot) error {
	err := store.Save(snapshot)
	return err
}

func DeleteSnapshot(snapshot *Snapshot) error {
	err := store.Delete(snapshot)
	return err
}

func GetSnapshotList(from, size int, taskId string) (int, []Snapshot, error) {
	var snapshots []Snapshot
	queryO := store.Query{Sort: "create_time desc", From: from, Size: size}
	if len(taskId) > 0 {
		queryO.Conds = store.And(store.Eq("task_id", taskId))
	}
	err, result := store.Search(&snapshots, &queryO)
	if err != nil {
		log.Error(err)
	}
	return result.Total, snapshots, err
}

func GetSnapshotAllList(taskId string) (int, []Snapshot, error) {
	var snapshots []Snapshot
	queryO := store.Query{Sort: "create_time desc"}
	if len(taskId) > 0 {
		queryO.Conds = store.And(store.Eq("task_id", taskId))
	}
	err, result := store.Search(&snapshots, &queryO)
	if err != nil {
		log.Error(err)
	}
	return result.Total, snapshots, err
}

func GetSnapshot(id string) (Snapshot, error) {
	snapshot := Snapshot{}
	err := store.GetBy("id", id, &snapshot)
	if err != nil {
		log.Error(id, ", ", err)
	}
	if len(snapshot.ID) == 0 || snapshot.CreateTime == nil {
		panic(errors.New("not found," + id))
	}

	return snapshot, err
}
