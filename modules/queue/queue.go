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

package queue

import (
	log "github.com/cihub/seelog"
	"time"
	. "github.com/medcl/gopa/core/env"
)

func Start(env *Env) {

	l:=GopaLogger{}
	dq := newDiskQueue("task", env.RuntimeConfig.PathConfig.TaskData, 1024, 4, 1<<10, 2500, 2*time.Second, l)
	defer dq.Close()


	//env.RuntimeConfig.Storage = &store
	log.Info("queue success started")

}

func Stop() error {
	return nil
}
