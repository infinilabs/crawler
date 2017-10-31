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
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
)

const name string = "save_task"

type SaveTaskJoint struct {
	model.Parameters
}

const isCreate model.ParaKey = "is_create"

func (joint SaveTaskJoint) IsCreate(v bool) SaveTaskJoint {
	joint.Init()
	joint.Set(isCreate, v)
	return joint
}

func (joint SaveTaskJoint) Name() string {
	return name
}

func (joint SaveTaskJoint) Process(context *model.Context) error {

	log.Trace("end process")
	if context.IsExit() {
		return errors.NewWithCode(errors.New("error in process"), config.ErrorExitedPipeline, "pipeline exited")
	}

	task := context.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	task.Phrase = context.Phrase
	task.Message = util.ToJson(context.Payload, false)

	if joint.GetBool(isCreate, false) {
		log.Trace("create task, url:", task.Url)
		model.CreateTask(task)
	} else {
		model.UpdateTask(task)
	}

	return nil
}
