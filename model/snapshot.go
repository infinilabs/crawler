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
	"github.com/infinitbyte/framework/core/errors"
	"github.com/infinitbyte/framework/core/persist"
	"github.com/infinitbyte/framework/core/util"
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
	ID      string `json:"id,omitempty" index:"id"`
	Version int    `json:"version,omitempty"`
	Url     string `json:"url,omitempty"`
	TaskID  string `json:"task_id,omitempty"`
	Path    string `json:"path,omitempty"` //path of this file
	File    string `json:"file,omitempty"` //filename of this page
	Ext     string `json:"ext,omitempty"`  //extension of filename

	StatusCode int    `json:"-"`
	Payload    []byte `json:"-"`
	Size       uint64 `json:"size,omitempty"`

	ScreenshotID string `json:"screenshot_id,omitempty"`

	Headers    map[string][]string     `json:"-"`
	Metadata   *map[string]interface{} `json:"-"`
	Parameters []KV                    `json:"-"`

	Language string `json:"lang,omitempty"`

	Title       string `json:"title,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Text        string `json:"text,omitempty"`
	ContentType string `json:"content_type,omitempty"`

	Tags []string `json:"tags,omitempty"`

	Links LinkGroup `json:"links,omitempty"`

	Images struct {
		Internal []PageLink `json:"internal,omitempty"`
		External []PageLink `json:"external,omitempty"`
	} `json:"images,omitempty"`

	H1     []string `json:"h1,omitempty"`
	H2     []string `json:"h2,omitempty"`
	H3     []string `json:"h3,omitempty"`
	H4     []string `json:"h4,omitempty"`
	H5     []string `json:"h5,omitempty"`
	Bold   []string `json:"bold,omitempty"`
	Italic []string `json:"italic,omitempty"`

	Classifications  []string                `json:"classifications,omitempty"`
	EnrichedFeatures *map[string]interface{} `json:"enriched_features,omitempty"`

	Hash    string `json:"hash,omitempty"`
	SimHash string `json:"sim_hash,omitempty"`

	Created time.Time `json:"created,omitempty"`
}

type PageLink struct {
	Url   string `json:"url,omitempty"`
	Label string `json:"label,omitempty"`
}

func CreateSnapshot(snapshot *Snapshot) error {
	return persist.Save(snapshot)
}

func DeleteSnapshot(snapshot *Snapshot) error {
	return persist.Delete(snapshot)
}

func GetSnapshotList(from, size int, taskId string) (int, []Snapshot, error) {
	var snapshots []Snapshot
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.ASC})
	query := persist.Query{Sort: &sort, From: from, Size: size}
	if len(taskId) > 0 {
		query.Conds = persist.And(persist.Eq("task_id", taskId))
	}
	err, result := persist.Search(Snapshot{}, &snapshots, &query)
	if err != nil {
		log.Error(err)
		return 0, snapshots, err
	}
	if result.Result != nil && snapshots == nil || len(snapshots) == 0 {
		convertSnapshot(result, &snapshots)
	}

	return result.Total, snapshots, err
}

func GetSnapshotByField(k, v string) ([]Snapshot, error) {
	log.Trace("start get snapshot: ", k, ", ", v)
	snapshot := Snapshot{}
	snapshots := []Snapshot{}
	err, result := persist.GetBy(k, v, snapshot, &snapshots)

	if err != nil {
		log.Error(k, ", ", err)
		return snapshots, err
	}
	if result.Result != nil && snapshots == nil || len(snapshots) == 0 {
		convertSnapshot(result, &snapshots)
	}

	return snapshots, err
}

func GetSnapshot(id string) (Snapshot, error) {
	snapshot := Snapshot{}
	snapshot.ID = id
	err := persist.Get(&snapshot)
	if err != nil {
		log.Error(err)
		return snapshot, err
	}
	if len(snapshot.ID) == 0 || snapshot.Created.IsZero() {
		panic(errors.New("not found," + id))
	}

	return snapshot, err
}

func convertSnapshot(result persist.Result, snapshots *[]Snapshot) {
	if result.Result == nil {
		return
	}

	t, ok := result.Result.([]interface{})
	if ok {
		for _, i := range t {
			js := util.ToJson(i, false)
			t := Snapshot{}
			util.FromJson(js, &t)
			*snapshots = append(*snapshots, t)
		}
	}
}
