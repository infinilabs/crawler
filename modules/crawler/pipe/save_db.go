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

package pipe

import (
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/model"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/core/store"
	"github.com/medcl/gopa/core/util"
	"github.com/medcl/gopa/modules/config"
	"path"
	"strings"
)

const SaveSnapshotToDB JointKey = "save_snapshot_db"

type SaveSnapshotToDBJoint struct {
	CompressBody bool
	Bucket       string
}

func (this SaveSnapshotToDBJoint) Name() string {
	return string(SaveSnapshotToDB)
}

func (this SaveSnapshotToDBJoint) Process(c *Context) error {
	task := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := c.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	url := task.Url

	savePath := snapshot.Path
	saveFile := snapshot.File
	domain := task.Host

	saveKey := GetKey(path.Join(task.Host, savePath, saveFile))
	log.Debug("save url to db, url:", url, ",domain:", task.Host, ",path:", savePath, ",file:", saveFile, ",saveKey:", string(saveKey))

	if this.CompressBody {
		store.AddValueCompress(config.SnapshotBucketKey, saveKey, snapshot.Payload)

	} else {
		store.AddValue(config.SnapshotBucketKey, saveKey, snapshot.Payload)
	}

	stats.IncrementBy("domain.stats", domain+"."+config.STATS_STORAGE_FILE_SIZE, int64(len(snapshot.Payload)))
	stats.Increment("domain.stats", domain+"."+config.STATS_STORAGE_FILE_COUNT)

	return nil
}

const KeyDelimiter string = "/"

func GetKey(args ...string) []byte {
	key := config.SnapshotMappingBucketKey
	url := []byte(strings.Join(args, KeyDelimiter))
	v := store.GetValue(key, url)
	if v != nil {
		stats.Increment("save", "duplicated_url")
		log.Warnf("get snapshotId from db, maybe previous already saved, %s, %s", string(v), string(url))
		return v
	}
	snapshotId := util.GetIncrementID("snapshot")
	store.AddValue(key, url, []byte(snapshotId))
	return []byte(snapshotId)
}
