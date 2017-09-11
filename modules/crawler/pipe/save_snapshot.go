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
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/persist"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/modules/config"
	"strings"
	"time"
)

type SaveSnapshotToDBJoint struct {
	Parameters
}

const compressEnabled ParaKey = "compress_enabled"
const bucket ParaKey = "bucket"
const decelerateSteps ParaKey = "decelerate_steps"
const accelerateSteps ParaKey = "accelerate_steps"
const maxRevision ParaKey = "max_revision"

func (this SaveSnapshotToDBJoint) Name() string {
	return "save_snapshot_db"
}

var oneSecond, _ = time.ParseDuration("1s")

func (this SaveSnapshotToDBJoint) Process(c *Context) error {
	task := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := c.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	if snapshot == nil {
		return errors.Errorf("snapshot is nil, %s , %s", task.ID, task.Url)
	}

	//init decelerate steps, unit is the second
	decelerateSteps := initFetchRateArr(this.GetStringOrDefault(decelerateSteps, "24h,48h,72h,168h,360h,720h"))
	//init accelerate steps, unit is the second
	accelerateSteps := initFetchRateArr(this.GetStringOrDefault(accelerateSteps, "24h,12h,6h,3h,1h30m,45m,30m,20m,10m"))

	current := time.Now().UTC()

	//update task's snapshot, detect duplicated snapshot
	if snapshot.Hash == task.SnapshotHash {

		//increase next check time
		updateNextCheckTime(task, current, decelerateSteps, false)

		msg := fmt.Sprintf("content unchanged, snapshot with same hash: %s, %s", snapshot.Hash, task.Url)

		c.End(msg)

		return errors.New(msg)
	}

	task.SnapshotVersion = task.SnapshotVersion + 1
	task.SnapshotID = snapshot.ID
	task.SnapshotHash = snapshot.Hash
	task.SnapshotSimHash = snapshot.SimHash
	task.SnapshotCreated = snapshot.Created

	snapshot.Version = task.SnapshotVersion
	snapshot.Url = task.Url
	snapshot.TaskID = task.ID

	savePath := snapshot.Path
	saveFile := snapshot.File

	saveKey := []byte(snapshot.ID)

	log.Debug("save snapshot to db, url:", task.Url, ",domain:", task.Host, ",path:", savePath, ",file:", saveFile, ",saveKey:", string(saveKey))

	var err error
	if this.GetBool(compressEnabled, true) {
		err = persist.AddValueCompress(this.MustGetString(bucket), saveKey, snapshot.Payload)
	} else {
		err = persist.AddValue(this.MustGetString(bucket), saveKey, snapshot.Payload)
	}
	if err != nil {
		return err
	}

	model.CreateSnapshot(snapshot)

	updateNextCheckTime(task, current, accelerateSteps, true)

	deleteRedundantSnapShot(int(this.MustGetInt64(maxRevision)), this.MustGetString(bucket), task.ID)

	stats.IncrementBy("domain.stats", task.Host+"."+config.STATS_STORAGE_FILE_SIZE, int64(len(snapshot.Payload)))
	stats.Increment("domain.stats", task.Host+"."+config.STATS_STORAGE_FILE_COUNT)

	return nil
}

//init the fetch rate array by cfg parameters
func initFetchRateArr(velocityStr string) []int {
	arrVelocityStr := strings.Split(velocityStr, ",")
	var velocityArr = make([]int, len(arrVelocityStr), len(arrVelocityStr))
	for i := 0; i < len(arrVelocityStr); i++ {
		m, err := time.ParseDuration(arrVelocityStr[i])
		if err == nil {
			velocityArr[i] = int(m.Seconds())
		} else {
			panic(fmt.Errorf("%s invalid config,only supports h, m, s", velocityStr))
		}
	}
	return velocityArr
}

//update the snapshot's next check time
func updateNextCheckTime(task *model.Task, current time.Time, steps []int, changed bool) {

	if len(steps) < 1 {
		panic(errors.New("invalid steps"))
	}

	//set one day as default next check time, unit is the second
	var timeIntervalNext = 24 * 60 * 60

	if task.SnapshotVersion <= 1 && task.LastCheck == nil && task.NextCheck == nil {

		timeIntervalNext = steps[0]

	} else {
		timeIntervalLast := getTimeInterval(*task.LastCheck, *task.NextCheck)

		if changed {
			arrTimeLength := len(steps)
			for i := 1; i < arrTimeLength; i++ {
				if timeIntervalLast > steps[0] {
					timeIntervalNext = steps[0]
					break
				}
				if timeIntervalLast < steps[arrTimeLength-2] {
					timeIntervalNext = steps[arrTimeLength-1]
					break
				}
				if i+1 >= arrTimeLength {
					timeIntervalNext = steps[arrTimeLength-1]
					break
				}
				if timeIntervalLast <= steps[i-1] && timeIntervalLast > steps[i] {
					timeIntervalNext = steps[i]
					break
				}
			}

		} else {
			arrTimeLength := len(steps)
			for i := 1; i < arrTimeLength; i++ {
				if timeIntervalLast < steps[0] {
					timeIntervalNext = steps[0]
					break
				}
				if timeIntervalLast > steps[arrTimeLength-2] {
					timeIntervalNext = steps[arrTimeLength-1]
					break
				}
				if i+1 >= arrTimeLength {
					timeIntervalNext = steps[arrTimeLength-1]
					break
				}
				if timeIntervalLast >= steps[i-1] && timeIntervalLast < steps[i] {
					timeIntervalNext = steps[i]
					break
				}
			}
		}
	}

	task.LastCheck = &current
	nextT := current.Add(oneSecond * time.Duration(timeIntervalNext))
	task.NextCheck = &nextT
}

func getTimeInterval(timeStart time.Time, timeEnd time.Time) int {
	ts := timeStart.Sub(timeEnd).Seconds()
	if ts < 0 {
		ts = -ts
	}
	return int(ts)
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
					persist.DeleteValue(bucketStr, []byte(snapshotsList[i].ID), snapshotsList[i].Payload)
				}
			}
		}
	}
}
