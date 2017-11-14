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
	"errors"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/util"
	"time"
)

// InitTaskJoint basically start the pipeline process, construct a model.Task, may loaded from db with CONTEXT_TASK_ID or manually passed in with CONTEXT_CRAWLER_TASK
type InitTaskJoint struct {
	model.Parameters
}

// Name return: init_task
func (joint InitTaskJoint) Name() string {
	return "init_task"
}

// Process task load, init a new snapshot instance
func (joint InitTaskJoint) Process(context *model.Context) error {

	if context.Has(model.CONTEXT_TASK_ID) && !context.Has(model.CONTEXT_TASK_URL) {
		//init task record
		t, err := model.GetTask(context.MustGetString(model.CONTEXT_TASK_ID))
		if err != nil {
			context.Exit("task init error")
			panic(err)
		}
		context.Set(model.CONTEXT_TASK_ID, t.ID)
		context.Set(model.CONTEXT_TASK_URL, t.Url)
		context.Set(model.CONTEXT_TASK_Reference, t.Reference)
		context.Set(model.CONTEXT_TASK_Depth, t.Depth)
		context.Set(model.CONTEXT_TASK_Breadth, t.Breadth)
		context.Set(model.CONTEXT_TASK_Host, t.Host)
		context.Set(model.CONTEXT_TASK_Schema, t.Schema)
		context.Set(model.CONTEXT_TASK_OriginalUrl, t.OriginalUrl)
		context.Set(model.CONTEXT_TASK_Status, t.Status)
		context.Set(model.CONTEXT_TASK_Message, t.Message)
		context.Set(model.CONTEXT_TASK_Created, t.Created)
		context.Set(model.CONTEXT_TASK_Updated, t.Updated)
		context.Set(model.CONTEXT_TASK_LastFetch, t.LastFetch)
		context.Set(model.CONTEXT_TASK_LastCheck, t.LastCheck)
		context.Set(model.CONTEXT_TASK_NextCheck, t.NextCheck)
		context.Set(model.CONTEXT_TASK_SnapshotID, t.SnapshotID)
		context.Set(model.CONTEXT_TASK_SnapshotSimHash, t.SnapshotSimHash)
		context.Set(model.CONTEXT_TASK_SnapshotHash, t.SnapshotHash)
		context.Set(model.CONTEXT_TASK_SnapshotCreated, t.SnapshotCreated)
		context.Set(model.CONTEXT_TASK_SnapshotVersion, t.SnapshotVersion)
		context.Set(model.CONTEXT_TASK_PipelineConfigID, t.PipelineConfigID)

	} else if !context.Has(model.CONTEXT_TASK_URL) {
		context.Exit("task init error")
		panic(errors.New("task not set"))
	}

	t1 := time.Now().UTC()

	//init snapshot
	var snapshot = &model.Snapshot{
		ID:      util.GetUUID(),
		Created: t1,
	}
	context.Set(model.CONTEXT_SNAPSHOT, snapshot)

	return nil
}
