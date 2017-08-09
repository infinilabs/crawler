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

func (this SaveSnapshotToDBJoint) Name() string {
	return string(SaveSnapshotToDB)
}

func (this SaveSnapshotToDBJoint) Process(c *Context) error {
	task := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := c.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	//init decelerateSteps
	arrDecelerateSteps = initGrabVelocityArr(this.MustGetString(decelerateSteps))
	//init accelerateSteps
	arrAccelerateSteps = initGrabVelocityArr(this.MustGetString(accelerateSteps))

	//update task's snapshot, detect duplicated snapshot
	if snapshot != nil {

		tNow := time.Now().UTC()
		m, _ := time.ParseDuration("1s")

		if snapshot.Hash == task.SnapshotHash {
			log.Debug(fmt.Sprintf("break by same hash: %s, %s", snapshot.Hash, task.Url))
			c.Break(fmt.Sprintf("same hash: %s, %s", snapshot.Hash, task.Url))

			//extend the nextchecktime
			setSnapNextCheckTime(task,tNow,m,false)

			return nil
		}

		//shorten the nextchecktime
		setSnapNextCheckTime(task,tNow,m,true)

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

	deleteRedundantSnapShot(int(this.MustGetInt64(maxRevision)),this.MustGetString(bucket),task.ID)

	stats.IncrementBy("domain.stats", domain+"."+config.STATS_STORAGE_FILE_SIZE, int64(len(snapshot.Payload)))
	stats.Increment("domain.stats", domain+"."+config.STATS_STORAGE_FILE_COUNT)

	return nil
}

//unit is the second
var arrDecelerateSteps []int
var arrAccelerateSteps []int

//init the grab velocity array by cfg parameters
func initGrabVelocityArr(velocityStr string) []int {
	arrVelocityStr := strings.Split(velocityStr,",")
	var velocityArr = make([]int,len(arrVelocityStr),len(arrVelocityStr))
	for i := 0; i < len(arrVelocityStr); i++ {
		m,err := time.ParseDuration(arrVelocityStr[i])
		if err == nil {
			velocityArr[i] = int(m.Seconds())
		}
	}
	return velocityArr
}

//set snapshot nextchecktime
func setSnapNextCheckTime(task *model.Task,timeNow time.Time,timeDuration time.Duration,fetchSuccess bool){
	task.LastCheckTime = &timeNow
	if task.SnapshotCreateTime == nil {
		defaultTime := timeNow.Add(-timeDuration * 1)
		task.SnapshotCreateTime = &defaultTime
	}
	timeInterval := GetNextCheckTimeSeconds(fetchSuccess, *task.SnapshotCreateTime, timeNow)
	nextT := timeNow.Add(timeDuration * time.Duration(timeInterval))
	task.NextCheckTime = &nextT
}

func GetNextCheckTimeSeconds(fetchSuccess bool, tSnapshotCreateTime time.Time, tTimeNow time.Time) int {
	timeIntervalLast := GetTimeInterval(tSnapshotCreateTime, tTimeNow)
	//set one day as default time,unit is the second
	timeIntervalNext := 24 * 60 * 60
	if fetchSuccess {
		arrTimeLength := len(arrAccelerateSteps)
		for i := 1; i < arrTimeLength; i++ {
			if timeIntervalLast >= arrAccelerateSteps[0] {
				timeIntervalNext = arrAccelerateSteps[0]
				break
			}
			if timeIntervalLast < arrAccelerateSteps[arrTimeLength-2] {
				timeIntervalNext = arrAccelerateSteps[arrTimeLength-1]
				break
			}
			if i+1 >= arrTimeLength {
				timeIntervalNext = arrAccelerateSteps[arrTimeLength-1]
				break
			}
			if timeIntervalLast <= arrAccelerateSteps[i-1] && timeIntervalLast > arrAccelerateSteps[i] {
				timeIntervalNext = arrAccelerateSteps[i]
				break
			}
		}

	} else {
		arrTimeLength := len(arrDecelerateSteps)
		for i := 1; i < arrTimeLength; i++ {
			if timeIntervalLast <= arrDecelerateSteps[0] {
				timeIntervalNext = arrDecelerateSteps[0]
				break
			}
			if timeIntervalLast > arrDecelerateSteps[arrTimeLength-2] {
				timeIntervalNext = arrDecelerateSteps[arrTimeLength-1]
				break
			}
			if i+1 >= arrTimeLength {
				timeIntervalNext = arrDecelerateSteps[arrTimeLength-1]
				break
			}
			if timeIntervalLast >= arrDecelerateSteps[i-1] && timeIntervalLast < arrDecelerateSteps[i] {
				timeIntervalNext = arrDecelerateSteps[i]
				break
			}
		}
	}
	return timeIntervalNext
}

func GetTimeInterval(timeStart time.Time, timeEnd time.Time) int {
	ts := timeStart.Sub(timeEnd).Seconds()
	if ts < 0 {
		ts = -ts
	}
	return int(ts)
}

//TODO optimization algorithm
func deleteRedundantSnapShot(maxRevisionNum int,bucketStr string,taskId string){
	//get current snapshot list and total num
	snapshotTotal, _, err := model.GetSnapshotList(0,1,taskId)
	if err == nil {
		//get max snapshot num
		maxSnapshotNum := maxRevisionNum
		//if more than max snapshot num,delete old snapshot
		if snapshotTotal > maxSnapshotNum {
			mustDeleteNum := snapshotTotal - maxSnapshotNum
			_, snapshotsList, errReadList := model.GetSnapshotList(1,mustDeleteNum,taskId)
			if errReadList == nil {
				for i := 0; i < len(snapshotsList); i++  {
					model.DeleteSnapshot(&snapshotsList[i])
					store.DeleteValue(bucketStr,[]byte(snapshotsList[i].ID),snapshotsList[i].Payload)
				}
			}
		}
	}
}