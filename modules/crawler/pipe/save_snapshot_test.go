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
	arrDecelerateSteps = initGrabVelocityArr("1m,10m,20m,30m,60m,1h30m,3h,6h,12h,24h,48h,168h,360h")
	fmt.Println(arrDecelerateSteps)

	arrAccelerateSteps = initGrabVelocityArr("1m10s,60s,50s,40s,30s,20s,10s,6s,3s")
	fmt.Println(arrAccelerateSteps)

	task := new(model.Task)
	tNow := time.Now().UTC()
	m, _ := time.ParseDuration("1s")
	setSnapNextCheckTime(task,tNow,m,true)
	fmt.Println("task.LastCheckTime     ",task.LastCheckTime)
	fmt.Println("task.SnapshotCreateTime",task.SnapshotCreateTime)
	fmt.Println("task.NextCheckTime     ",task.NextCheckTime)
}

func TestDeleteRedundantSnapShot(t *testing.T){
	deleteRedundantSnapShot(10,"","")
}
