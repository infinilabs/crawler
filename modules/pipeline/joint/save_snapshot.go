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

package joint

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/modules/config"
)

type SaveSnapshotToDBJoint struct {
	model.Parameters
}

const compressEnabled model.ParaKey = "compress_enabled"
const bucket model.ParaKey = "bucket"
const maxRevision model.ParaKey = "max_revision"

func (this SaveSnapshotToDBJoint) Name() string {
	return "save_snapshot_db"
}

func (this SaveSnapshotToDBJoint) Process(c *model.Context) error {

	taskID := c.MustGetString(model.CONTEXT_TASK_ID)
	taskUrl := c.MustGetString(model.CONTEXT_TASK_URL)
	taskHost := c.MustGetString(model.CONTEXT_TASK_Host)
	previousSnapshotHash := c.MustGetString(model.CONTEXT_TASK_SnapshotHash)
	previousSnapshotVersion := c.GetIntOrDefault(model.CONTEXT_TASK_SnapshotVersion, 0)

	snapshot := c.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	if snapshot == nil {
		return errors.Errorf("snapshot is nil, %s , %s", taskID, taskUrl)
	}

	//update task's snapshot, detect duplicated snapshot
	if snapshot.Hash != "" && snapshot.Hash == previousSnapshotHash {
		msg := fmt.Sprintf("content unchanged, snapshot with same hash: %s, %s, prev hash: %s,prev version: %s", snapshot.Hash, taskUrl, previousSnapshotHash, previousSnapshotVersion)
		c.End(msg)
		return nil
	}

	previousSnapshotVersion = previousSnapshotVersion + 1

	snapshot.Version = previousSnapshotVersion
	snapshot.Url = taskUrl
	snapshot.TaskID = taskID

	savePath := snapshot.Path
	saveFile := snapshot.File

	saveKey := []byte(snapshot.ID)

	log.Debug("save snapshot to db, url:", taskUrl, ",path:", savePath, ",file:", saveFile, ",saveKey:", string(saveKey))

	bucketName := this.GetStringOrDefault(bucket, "Snapshot")

	var err error
	if this.GetBool(compressEnabled, true) {
		err = persist.AddValueCompress(bucketName, saveKey, snapshot.Payload)
	} else {
		err = persist.AddValue(bucketName, saveKey, snapshot.Payload)
	}
	if err != nil {
		return err
	}

	model.CreateSnapshot(snapshot)

	c.Set(model.CONTEXT_TASK_SnapshotID, snapshot.ID)
	c.Set(model.CONTEXT_TASK_SnapshotVersion, previousSnapshotVersion)
	c.Set(model.CONTEXT_TASK_SnapshotHash, snapshot.Hash)
	c.Set(model.CONTEXT_TASK_SnapshotSimHash, snapshot.SimHash)
	c.Set(model.CONTEXT_TASK_SnapshotCreated, snapshot.Created)

	deleteRedundantSnapShot(int(this.GetInt64OrDefault(maxRevision, 5)), bucketName, taskID)

	stats.IncrementBy("host.stats", taskHost+"."+config.STATS_STORAGE_FILE_SIZE, int64(len(snapshot.Payload)))
	stats.Increment("host.stats", taskHost+"."+config.STATS_STORAGE_FILE_COUNT)

	return nil
}

//TODO optimization algorithm
func deleteRedundantSnapShot(maxRevisionNum int, bucketStr string, taskId string) {
	//get current snapshot list and total num
	snapshotTotal, _, err := model.GetSnapshotList(0, 1, taskId)
	if err == nil {
		//get max snapshot num
		maxSnapshotNum := maxRevisionNum
		//if more than max snapshot num,delete old snapshot
		if snapshotTotal > maxSnapshotNum {
			mustDeleteNum := snapshotTotal - maxSnapshotNum
			_, snapshotsList, errReadList := model.GetSnapshotList(1, mustDeleteNum, taskId)
			if errReadList == nil {
				for i := 0; i < len(snapshotsList); i++ {
					model.DeleteSnapshot(&snapshotsList[i])
					persist.DeleteKey(bucketStr, []byte(snapshotsList[i].ID))
				}
			}
		}
	}
}
