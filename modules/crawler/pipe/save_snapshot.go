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
	"github.com/infinitbyte/gopa/core/model"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/store"
	"github.com/infinitbyte/gopa/modules/config"
	"time"
	"strings"
)

const SaveSnapshotToDB JointKey = "save_snapshot_db"

type SaveSnapshotToDBJoint struct {
	Parameters
}

const compressEnabled ParaKey = "compress_enabled"
const bucket ParaKey = "bucket"
const decelerateSteps ParaKey = "decelerate_steps"
const accelerateSteps ParaKey = "accelerate_steps"
const maxRevision ParaKey = "max_revision"

//minutes
var arrDecelerateSteps []int
var arrAccelerateSteps []int

func (this SaveSnapshotToDBJoint) Name() string {
	return string(SaveSnapshotToDB)
}

func (this SaveSnapshotToDBJoint) Process(c *Context) error {
	task := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := c.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	//init decelerateSteps
	arrDecelerateStepsStr := strings.Split(this.MustGetString(decelerateSteps),",")
	arrDecelerateSteps = make([]int,len(arrDecelerateStepsStr),len(arrDecelerateStepsStr))
	for i := 0; i < len(arrDecelerateStepsStr); i++ {
		m,error := time.ParseDuration(arrDecelerateStepsStr[i])
		if error == nil {
			arrDecelerateSteps[i] = int(m.Minutes())
		}
	}
	//init accelerateSteps
	arrAccelerateStepsStr := strings.Split(this.MustGetString(accelerateSteps),",")
	arrAccelerateSteps = make([]int,len(arrAccelerateStepsStr),len(arrAccelerateStepsStr))
	for i := 0; i < len(arrAccelerateStepsStr); i++ {
		m,error := time.ParseDuration(arrAccelerateStepsStr[i])
		if error == nil {
			arrAccelerateSteps[i] = int(m.Minutes())
		}
	}

	//update task's snapshot, detect duplicated snapshot
	if snapshot != nil {

		tNow := time.Now().UTC()
		m, _ := time.ParseDuration("1m")

		if snapshot.Hash == task.SnapshotHash {
			log.Debug(fmt.Sprintf("break by same hash: %s, %s", snapshot.Hash, task.Url))
			c.Break(fmt.Sprintf("same hash: %s, %s", snapshot.Hash, task.Url))

			//extended the nextchecktime
			task.LastCheckTime = &tNow
			if task.SnapshotCreateTime == nil {
				defaultTime := tNow.Add(-m * 1)
				task.SnapshotCreateTime = &defaultTime
			}
			timeInterval := GetNextCheckTimeMinutes(false, *task.SnapshotCreateTime, tNow)
			nextT := tNow.Add(m * time.Duration(timeInterval))
			task.NextCheckTime = &nextT

			return nil
		}

		//shorten the nextchecktime
		task.LastCheckTime = &tNow
		if task.SnapshotCreateTime == nil {
			defaultTime := tNow.Add(-m * 1)
			task.SnapshotCreateTime = &defaultTime
		}
		timeInterval := GetNextCheckTimeMinutes(true, *task.SnapshotCreateTime, tNow)
		nextT := tNow.Add(m * time.Duration(timeInterval))
		task.NextCheckTime = &nextT

		task.SnapshotVersion = task.SnapshotVersion + 1
		task.SnapshotID = snapshot.ID
		task.SnapshotHash = snapshot.Hash
		task.SnapshotSimHash = snapshot.SimHash
		task.SnapshotCreateTime = snapshot.CreateTime

		snapshot.Version = task.SnapshotVersion
		snapshot.Url = task.Url
		snapshot.TaskID = task.ID
	}

	url := task.Url

	savePath := snapshot.Path
	saveFile := snapshot.File
	domain := task.Host

	saveKey := []byte(snapshot.ID)

	log.Debug("save url to db, url:", url, ",domain:", task.Host, ",path:", savePath, ",file:", saveFile, ",saveKey:", string(saveKey))

	if this.GetBool(compressEnabled, true) {
		store.AddValueCompress(this.MustGetString(bucket), saveKey, snapshot.Payload)

	} else {
		store.AddValue(this.MustGetString(bucket), saveKey, snapshot.Payload)
	}

	model.CreateSnapshot(snapshot)

	//TODO optimization algorithm
	//get current snapshot list and total num
	snapshotTotal, _, err := model.GetSnapshotList(0,1,task.ID)
	if err == nil {
		//get max snapshot num
		maxSnapshotNum := int(this.MustGetInt64(maxRevision))
		//if more than max snapshot num,delete old snapshot
		if snapshotTotal > maxSnapshotNum {
			mustDeleteNum := snapshotTotal - maxSnapshotNum
			_, snapshotsList, errReadList := model.GetSnapshotList(1,mustDeleteNum,task.ID)
			if errReadList == nil {
				for i := 0; i < len(snapshotsList); i++  {
					model.DeleteSnapshot(&snapshotsList[i])
					store.DeleteValue(this.MustGetString(bucket),[]byte(snapshotsList[i].ID),snapshotsList[i].Payload)
				}
			}
		}
	}


	stats.IncrementBy("domain.stats", domain+"."+config.STATS_STORAGE_FILE_SIZE, int64(len(snapshot.Payload)))
	stats.Increment("domain.stats", domain+"."+config.STATS_STORAGE_FILE_COUNT)



	return nil
}



func GetNextCheckTimeMinutes(fetchSuccess bool, tLastCheckTime time.Time, tNextCheckTime time.Time) int {
	timeIntervalLast := GetTimeInterval(tLastCheckTime, tNextCheckTime)
	timeIntervalNext := 24 * 60
	if fetchSuccess {
		arrTimeLength := len(arrDecelerateSteps)
		for i := 1; i < arrTimeLength; i++ {
			if timeIntervalLast > arrDecelerateSteps[0] {
				timeIntervalNext = arrDecelerateSteps[0]
				break
			}
			if timeIntervalLast <= arrDecelerateSteps[arrTimeLength-2] {
				timeIntervalNext = arrDecelerateSteps[arrTimeLength-1]
				break
			}
			if i+1 >= arrTimeLength {
				timeIntervalNext = arrDecelerateSteps[arrTimeLength-1]
				break
			}
			if timeIntervalLast <= arrDecelerateSteps[i] && timeIntervalLast > arrDecelerateSteps[i+1] {
				timeIntervalNext = arrDecelerateSteps[i+1]
				break
			}
		}
	} else {
		arrTimeLength := len(arrAccelerateSteps)
		for i := 1; i < arrTimeLength; i++ {
			if timeIntervalLast <= arrAccelerateSteps[0] {
				timeIntervalNext = arrAccelerateSteps[1]
				break
			}
			if timeIntervalLast >= arrAccelerateSteps[arrTimeLength-2] {
				timeIntervalNext = arrAccelerateSteps[arrTimeLength-1]
				break
			}
			if i+1 >= arrTimeLength {
				timeIntervalNext = arrAccelerateSteps[arrTimeLength-1]
				break
			}
			if timeIntervalLast >= arrAccelerateSteps[i] && timeIntervalLast < arrAccelerateSteps[i+1] {
				timeIntervalNext = arrAccelerateSteps[i+1]
				break
			}
		}
	}
	return timeIntervalNext
}

func GetTimeInterval(timeStart time.Time, timeEnd time.Time) int {
	ts := timeStart.Sub(timeEnd).Minutes()
	if ts < 0 {
		ts = -ts
	}
	return int(ts)
}
