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
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/model"
	"strings"
	"time"
)

type UpdateCheckTimeJoint struct {
	model.Parameters
}

const decelerateSteps model.ParaKey = "decelerate_steps"
const accelerateSteps model.ParaKey = "accelerate_steps"

func (this UpdateCheckTimeJoint) Name() string {
	return "update_check_time"
}

var oneSecond, _ = time.ParseDuration("1s")

func (this UpdateCheckTimeJoint) Process(c *model.Context) error {
	snapshot := c.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	taskID := c.MustGetString(model.CONTEXT_TASK_ID)
	taskUrl := c.MustGetString(model.CONTEXT_TASK_URL)

	if snapshot == nil {
		return errors.Errorf("snapshot is nil, %s", taskID)
	}

	lastSnapshotHash := c.GetStringOrDefault(model.CONTEXT_TASK_SnapshotHash, "")

	//this control how page are updated, or update frequency, for example, by default, a page will be checked after 24h,
	//if the page doesn't change during this update fetch,
	//the next fetch time will be changed to 48h later, which means it will automatically delayed from 24h to 48h,
	//and if the page still not change after that 48h, then it will fetch the page again but 168h later
	decelerateSteps := initFetchRateArr(this.GetStringOrDefault(decelerateSteps, "24h,48h,72h,168h,360h,720h"))
	//you may consider this is reverse of decelerateSteps, by default a page will be updated after 24h,
	//if the page changed, it will try to fetch after 12h, if still change, will shorten the period again, but 10m is minim wait time,
	//but you can change the configuration
	accelerateSteps := initFetchRateArr(this.GetStringOrDefault(accelerateSteps, "24h,12h,6h,3h,1h30m,45m,30m,20m,10m"))

	current := time.Now().UTC()

	//update task's snapshot, detect duplicated snapshot
	if snapshot.Hash == lastSnapshotHash {

		//increase next check time
		updateNextCheckTime(c, current, decelerateSteps, false)

		msg := fmt.Sprintf("content unchanged, snapshot with same hash: %s, %s", snapshot.Hash, taskUrl)

		c.End(msg)

		return nil
	}

	updateNextCheckTime(c, current, accelerateSteps, true)

	return nil
}

//init the fetch rate array by cfg
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
func updateNextCheckTime(c *model.Context, current time.Time, steps []int, changed bool) {

	if len(steps) < 1 {
		panic(errors.New("invalid steps"))
	}

	lastSnapshotHash := c.GetStringOrDefault(model.CONTEXT_TASK_SnapshotHash, "")
	lastSnapshotVer := c.GetIntOrDefault(model.CONTEXT_TASK_SnapshotVersion, 0)
	taskLastCheck, b1 := c.GetTime(model.CONTEXT_TASK_LastCheck)
	taskNextCheck, b2 := c.GetTime(model.CONTEXT_TASK_NextCheck)

	if lastSnapshotHash == "" {

	}

	//set one day as default next check time, unit is the second
	var timeIntervalNext = 24 * 60 * 60

	if lastSnapshotVer <= 1 && !b1 && !b2 {

		timeIntervalNext = steps[0]

	} else {
		timeIntervalLast := getTimeInterval(taskLastCheck, taskNextCheck)

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

	c.Set(model.CONTEXT_TASK_LastCheck, current)
	nextT := current.Add(oneSecond * time.Duration(timeIntervalNext))
	c.Set(model.CONTEXT_TASK_NextCheck, nextT)
}

func getTimeInterval(timeStart time.Time, timeEnd time.Time) int {
	ts := timeStart.Sub(timeEnd).Seconds()
	if ts < 0 {
		ts = -ts
	}
	return int(ts)
}
