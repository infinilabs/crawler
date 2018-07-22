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

package common

import (
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/model"
)

func ParseTask(context *pipeline.Context) *model.Task {
	task := model.Task{}
	task.ID = context.GetStringOrDefault(model.CONTEXT_TASK_ID, "")
	task.Url = context.MustGetString(model.CONTEXT_TASK_URL)
	task.Reference = context.GetStringOrDefault(model.CONTEXT_TASK_Reference, "")
	task.Depth = context.GetIntOrDefault(model.CONTEXT_TASK_Depth, 0)
	task.Breadth = context.GetIntOrDefault(model.CONTEXT_TASK_Breadth, 0)
	task.Host = context.GetStringOrDefault(model.CONTEXT_TASK_Host, "")
	task.Schema = context.GetStringOrDefault(model.CONTEXT_TASK_Schema, "")
	task.OriginalUrl = context.GetStringOrDefault(model.CONTEXT_TASK_OriginalUrl, "")

	if context.Has(model.CONTEXT_TASK_Status) {
		task.Status = context.MustGetInt(model.CONTEXT_TASK_Status)
	} else if context.IsEnd() {
		task.Status = model.TaskInterrupted
	}

	if context.Has(model.CONTEXT_TASK_Message) {
		task.Message = context.GetStringOrDefault(model.CONTEXT_TASK_Message, "")
	} else {
		task.Message = util.ToJson(context.Payload, false)
	}

	if context.Has(model.CONTEXT_TASK_Created) {
		task.Created = context.MustGetTime(model.CONTEXT_TASK_Created)
	}
	if context.Has(model.CONTEXT_TASK_Updated) {
		task.Updated = context.MustGetTime(model.CONTEXT_TASK_Updated)
	}
	if context.Has(model.CONTEXT_TASK_LastFetch) {
		task.LastFetch = context.MustGetTime(model.CONTEXT_TASK_LastFetch)
	}
	if context.Has(model.CONTEXT_TASK_LastCheck) {
		task.LastCheck = context.MustGetTime(model.CONTEXT_TASK_LastCheck)
	}
	if context.Has(model.CONTEXT_TASK_NextCheck) {
		task.NextCheck = context.MustGetTime(model.CONTEXT_TASK_NextCheck)
	}
	if context.Has(model.CONTEXT_TASK_SnapshotID) {
		task.SnapshotID = context.GetStringOrDefault(model.CONTEXT_TASK_SnapshotID, "")
	}
	if context.Has(model.CONTEXT_TASK_SnapshotSimHash) {
		task.SnapshotSimHash = context.GetStringOrDefault(model.CONTEXT_TASK_SnapshotSimHash, "")
	}
	if context.Has(model.CONTEXT_TASK_SnapshotHash) {
		task.SnapshotHash = context.GetStringOrDefault(model.CONTEXT_TASK_SnapshotHash, "")
	}
	if context.Has(model.CONTEXT_TASK_SnapshotCreated) {
		task.SnapshotCreated = context.MustGetTime(model.CONTEXT_TASK_SnapshotCreated)
	}
	if context.Has(model.CONTEXT_TASK_SnapshotVersion) {
		task.SnapshotVersion = context.GetIntOrDefault(model.CONTEXT_TASK_SnapshotVersion, 0)
	}
	if context.Has(model.CONTEXT_TASK_LastScreenshotID) {
		task.LastScreenshotID = context.GetStringOrDefault(model.CONTEXT_TASK_LastScreenshotID, "")
	}
	if context.Has(model.CONTEXT_TASK_PipelineConfigID) {
		task.PipelineConfigID = context.GetStringOrDefault(model.CONTEXT_TASK_PipelineConfigID, "")
	}
	return &task
}
