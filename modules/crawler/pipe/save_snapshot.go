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
	"strconv"
	"strings"
)

const SaveSnapshotToDB JointKey = "save_snapshot_db"

type SaveSnapshotToDBJoint struct {
	Parameters
}

const compressEnabled ParaKey = "compress_enabled"
const bucket ParaKey = "bucket"
const snapshottimeToless ParaKey = "snapshottime_toless"
const snapshottimeTomore ParaKey = "snapshottime_tomore"
const snapshotMaxnum ParaKey = "snapshot_maxnum"

//minutes
var arrTimeToLess []int
var arrTimeToMore []int

func (this SaveSnapshotToDBJoint) Name() string {
	return string(SaveSnapshotToDB)
}

func (this SaveSnapshotToDBJoint) Process(c *Context) error {
	task := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := c.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	//init snapshottimeToless
	arrTimeToLessStr := strings.Split(this.MustGetString(snapshottimeToless),",")
	arrTimeToLess = make([]int,len(arrTimeToLessStr),len(arrTimeToLessStr))
	for i := 0; i < len(arrTimeToLessStr); i++ {
		m,error := strconv.Atoi(arrTimeToLessStr[i])
		if error == nil {
			arrTimeToLess[i] = m
		}
	}
	//init snapshottimeTomore
	arrTimeToMoreStr := strings.Split(this.MustGetString(snapshottimeTomore),",")
	arrTimeToMore = make([]int,len(arrTimeToMoreStr),len(arrTimeToMoreStr))
	for i := 0; i < len(arrTimeToMoreStr); i++ {
		m,error := strconv.Atoi(arrTimeToMoreStr[i])
		if error == nil {
			arrTimeToMore[i] = m
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

	//delete old snapshot
	//get current snapshot list and total num
	snapshotTotal, snapshotsList, err := model.GetSnapshotAllList(task.ID)
	if err == nil {
		//get max snapshot num
		maxSnapshotNum := this.MustGetInt64(snapshotMaxnum)
		//if more than max snapshot num,delete old snapshot
		if int64(snapshotTotal) > maxSnapshotNum {
			mustDeleteNum := int64(snapshotTotal) - maxSnapshotNum
			for i := 0; i < len(snapshotsList); i++  {
				if i > 0 && i < len(snapshotsList)-1 && mustDeleteNum > 0 {
					model.DeleteSnapshot(&snapshotsList[i])
					store.DeleteValue(this.MustGetString(bucket),[]byte(snapshotsList[i].ID),snapshotsList[i].Payload)
					mustDeleteNum -= 1
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
		arrTimeLength := len(arrTimeToLess)
		for i := 1; i < arrTimeLength; i++ {
			if timeIntervalLast > arrTimeToLess[0] {
				timeIntervalNext = arrTimeToLess[0]
				break
			}
			if timeIntervalLast <= arrTimeToLess[arrTimeLength-2] {
				timeIntervalNext = arrTimeToLess[arrTimeLength-1]
				break
			}
			if i+1 >= arrTimeLength {
				timeIntervalNext = arrTimeToLess[arrTimeLength-1]
				break
			}
			if timeIntervalLast <= arrTimeToLess[i] && timeIntervalLast > arrTimeToLess[i+1] {
				timeIntervalNext = arrTimeToLess[i+1]
				break
			}
		}
	} else {
		arrTimeLength := len(arrTimeToMore)
		for i := 1; i < arrTimeLength; i++ {
			if timeIntervalLast <= arrTimeToMore[0] {
				timeIntervalNext = arrTimeToMore[1]
				break
			}
			if timeIntervalLast >= arrTimeToMore[arrTimeLength-2] {
				timeIntervalNext = arrTimeToMore[arrTimeLength-1]
				break
			}
			if i+1 >= arrTimeLength {
				timeIntervalNext = arrTimeToMore[arrTimeLength-1]
				break
			}
			if timeIntervalLast >= arrTimeToMore[i] && timeIntervalLast < arrTimeToMore[i+1] {
				timeIntervalNext = arrTimeToMore[i+1]
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
