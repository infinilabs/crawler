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
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/errors"
	"github.com/medcl/gopa/core/model"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/util"
	"github.com/medcl/gopa/modules/config"
)

const SaveTask JointKey = "save_task"

type SaveTaskJoint struct {
	IsCreate bool
}

func (this SaveTaskJoint) Name() string {
	return string(SaveTask)
}

func (this SaveTaskJoint) Process(context *Context) error {

	log.Trace("end process")
	if context.IsErrorExit() {
		return errors.NewWithCode(errors.New("error in process"), config.ErrorExitedPipeline, "pipeline exited")
	}

	task := context.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	task.Status = model.TaskFetchSuccess
	task.Phrase = context.Phrase

	if context.IsBreak() {
		log.Trace("broken pipeline,", context.Payload)
		task.Status = model.TaskFetchFailed
		task.Message = util.ToJson(context.Payload, false)
	}

	if this.IsCreate {
		log.Trace("create task, url:", task.Url)
		model.CreateTask(task)
	} else {
		model.UpdateTask(task)
	}

	return nil
}
