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
	"testing"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/stretchr/testify/assert"
	"time"
	"fmt"
)

func TestInitGrabVelocityArr(t *testing.T) {
	arrDecelerateSteps = initGrabVelocityArr("1m,10m,20m,30m,60m,1h30m,3h,6h,12h,24h,48h,168h,360h")
	fmt.Println(arrDecelerateSteps)

	arrAccelerateSteps = initGrabVelocityArr("24h,12h,6h,3h,1h30m,45m,20m,10m,1m")
	fmt.Println(arrAccelerateSteps)
}

func TestSetSnapNextCheckTime(t *testing.T){
	arrDecelerateSteps = initGrabVelocityArr("3s,6s,10s,20s,30s,40s,50s,60s,70s")
	arrAccelerateSteps = initGrabVelocityArr("70s,60s,50s,40s,30s,20s,10s,6s,3s")

	toBeCharge := "2017-01-01 00:00:00.0000000 +0000 UTC"
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)

	task := new(model.Task)
	task.SnapshotCreateTime = &theTime
	fmt.Println("----task.SnapshotCreateTime",task.SnapshotCreateTime)
	m, _ := time.ParseDuration("1s")
	tNow := theTime.Add(1*m)
	setSnapNextCheckTime(task,tNow,m,false)
	fmt.Println("    task.LastCheckTime     ",task.LastCheckTime)
	fmt.Println("    task.NextCheckTime     ",task.NextCheckTime)
	timeInterval := GetTimeInterval(*task.LastCheckTime,*task.NextCheckTime)
	fmt.Println("----timeInterval           ",timeInterval)
	assert.Equal(t, 3, timeInterval)

	tNow = theTime.Add(71*m)
	setSnapNextCheckTime(task,tNow,m,true)
	fmt.Println("    task.LastCheckTime     ",task.LastCheckTime)
	fmt.Println("    task.NextCheckTime     ",task.NextCheckTime)
	timeInterval = GetTimeInterval(*task.LastCheckTime,*task.NextCheckTime)
	fmt.Println("----timeInterval           ",timeInterval)
	assert.Equal(t, 70, timeInterval)
}

func TestDeleteRedundantSnapShot(t *testing.T){
	deleteRedundantSnapShot(10,"","")
}
