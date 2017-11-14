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
	"github.com/infinitbyte/gopa/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInitGrabVelocityArr(t *testing.T) {
	arrDecelerateSteps := initFetchRateArr("1m,10m,20m,30m,60m,1h30m,3h,6h,12h,24h,48h,168h,360h")
	fmt.Println(arrDecelerateSteps)

	arrAccelerateSteps := initFetchRateArr("24h,12h,6h,3h,1h30m,45m,20m,10m,1m")
	fmt.Println(arrAccelerateSteps)
}

func TestSetSnapNextCheckTime(t *testing.T) {
	decelerateSteps := initFetchRateArr("1m,2m,5m,10m")
	accelerateSteps := initFetchRateArr("10m,5m,2m,1m")

	fmt.Println("decelerate,", decelerateSteps)
	fmt.Println("accelerate,", accelerateSteps)

	toBeCharge := "2017-01-01 00:00:00.0000000 +0000 UTC"
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	oneSecond, _ := time.ParseDuration("1s")
	oneMinute, _ := time.ParseDuration("1m")

	fmt.Println("update 1s with no change")
	context := model.Context{}
	tNow := theTime.Add(1 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, theTime)
	context.Set(model.CONTEXT_TASK_NextCheck, tNow)
	context.Set(model.CONTEXT_TASK_SnapshotVersion, 2)

	updateNextCheckTime(&context, tNow, decelerateSteps, false)
	new1 := context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 := context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval := getTimeInterval(new1, new2)
	fmt.Println("---- next check time          ", timeInterval)
	assert.Equal(t, 60, timeInterval)

	fmt.Println()

	fmt.Println("update 10m with no change")
	tNow = theTime.Add(10 * oneMinute)
	context.Set(model.CONTEXT_TASK_LastCheck, theTime)
	context.Set(model.CONTEXT_TASK_NextCheck, tNow)
	context.Set(model.CONTEXT_TASK_SnapshotVersion, 2)
	updateNextCheckTime(&context, tNow, decelerateSteps, false)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("---- next check time          ", timeInterval)
	assert.Equal(t, 600, timeInterval)

	fmt.Println()

	fmt.Println("update 20m with no change")
	context = model.Context{}
	tNow = theTime.Add(10 * oneMinute)
	context.Set(model.CONTEXT_TASK_LastCheck, theTime)
	context.Set(model.CONTEXT_TASK_NextCheck, tNow)
	context.Set(model.CONTEXT_TASK_SnapshotVersion, 2)
	updateNextCheckTime(&context, tNow, decelerateSteps, false)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("---- next check time          ", timeInterval)
	assert.Equal(t, 600, timeInterval)

	fmt.Println()

	fmt.Println("update 2m with change")
	tNow = theTime.Add(120 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, theTime)
	context.Set(model.CONTEXT_TASK_NextCheck, tNow)
	context.Set(model.CONTEXT_TASK_SnapshotVersion, 2)
	updateNextCheckTime(&context, tNow, accelerateSteps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("----timeInterval           ", timeInterval)
	assert.Equal(t, 60, timeInterval)

	fmt.Println("update 10s with change")
	tNow = theTime.Add(10 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, theTime)
	context.Set(model.CONTEXT_TASK_NextCheck, tNow)
	context.Set(model.CONTEXT_TASK_SnapshotVersion, 2)
	updateNextCheckTime(&context, tNow, accelerateSteps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("----timeInterval           ", timeInterval)
	assert.Equal(t, 60, timeInterval)

	fmt.Println("update 1000s with change")
	tNow = theTime.Add(1000 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, theTime)
	context.Set(model.CONTEXT_TASK_NextCheck, tNow)
	context.Set(model.CONTEXT_TASK_SnapshotVersion, 2)
	updateNextCheckTime(&context, tNow, accelerateSteps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("----timeInterval           ", timeInterval)
	assert.Equal(t, 600, timeInterval)

	fmt.Println("update 500s with change")
	tNow = theTime.Add(500 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, theTime)
	context.Set(model.CONTEXT_TASK_NextCheck, tNow)
	context.Set(model.CONTEXT_TASK_SnapshotVersion, 2)
	updateNextCheckTime(&context, tNow, accelerateSteps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("----timeInterval           ", timeInterval)
	assert.Equal(t, 300, timeInterval)

	fmt.Println("update 600s with change")
	tNow = theTime.Add(600 * oneSecond)
	context.Set(model.CONTEXT_TASK_LastCheck, theTime)
	context.Set(model.CONTEXT_TASK_NextCheck, tNow)
	context.Set(model.CONTEXT_TASK_SnapshotVersion, 2)
	updateNextCheckTime(&context, tNow, accelerateSteps, true)
	new1 = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	new2 = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	timeInterval = getTimeInterval(new1, new2)
	fmt.Println("----timeInterval           ", timeInterval)
	assert.Equal(t, 300, timeInterval)
}
