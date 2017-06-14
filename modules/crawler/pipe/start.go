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
	"github.com/medcl/gopa/core/model"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/util"
	"time"
)

const InitTask JointKey = "init_task"

const TaskID ParaKey = "TASK_ID"

type InitTaskJoint struct {
	Parameters
	Task *model.Task
}

func NewTaskJoint(task *model.Task) Joint {
	return InitTaskJoint{Task: task}
}

func (this InitTaskJoint) Name() string {
	return string(InitTask)
}

func (this InitTaskJoint) Process(context *Context) error {

	seelog.Trace("start process")

	var task *model.Task

	if this.Task != nil {
		task = this.Task
	} else if this.Has(TaskID) {
		//init task record
		t, err := model.GetTask(this.MustGetString(TaskID))
		if err != nil {
			context.ErrorExit("task init error")
			panic(err)
		}
		task = &t
	} else {
		context.ErrorExit("task init error")
		panic(errors.New("task not set"))
	}

	if task == nil {
		context.ErrorExit("task init error")
		panic(errors.New("nil task"))
	}

	//init snapshot
	var snapshot = &model.Snapshot{
		ID: util.GetUUID(),
	}

	//update last check time
	t1 := time.Now().UTC()
	task.LastCheckTime = &t1

	//update next check time //TODO

	context.Set(CONTEXT_CRAWLER_TASK, task)
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, snapshot)

	return nil
}
