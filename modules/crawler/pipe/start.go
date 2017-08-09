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
	"errors"
	"github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/model"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/util"
	"time"
)

type InitTaskJoint struct {
	Parameters
}

func (joint InitTaskJoint) Name() string {
	return "init_task"
}

func (joint InitTaskJoint) Process(context *Context) error {

	seelog.Trace("start process")

	var task *model.Task

	if context.Has(CONTEXT_CRAWLER_TASK) {
		task = context.Get(CONTEXT_CRAWLER_TASK).(*model.Task)
	} else if context.Has(CONTEXT_TASK_ID) {
		//init task record
		t, err := model.GetTask(context.MustGetString(CONTEXT_TASK_ID))
		if err != nil {
			context.ErrorExit("task init error")
			panic(err)
		}
		task = &t
		context.Set(CONTEXT_CRAWLER_TASK, task)

	} else {
		context.ErrorExit("task init error")
		panic(errors.New("task not set"))
	}

	if task == nil {
		context.ErrorExit("task init error")
		panic(errors.New("nil task"))
	}

	t1 := time.Now().UTC()

	//init snapshot
	var snapshot = &model.Snapshot{
		ID:         util.GetUUID(),
		CreateTime: &t1,
	}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, snapshot)

	return nil
}
