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

package output

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/errors"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
	"github.com/infinitbyte/gopa/pipeline/joints/common"
)

const name string = "save_task"

type SaveTaskJoint struct {
	pipeline.Parameters
}

const isCreate pipeline.ParaKey = "is_create"
const keep404 pipeline.ParaKey = "keep_404"
const keepRedirected pipeline.ParaKey = "keep_redirected"

func (joint SaveTaskJoint) IsCreate(v bool) SaveTaskJoint {
	joint.Set(isCreate, v)
	return joint
}

func (joint SaveTaskJoint) Name() string {
	return name
}

func (joint SaveTaskJoint) Process(context *pipeline.Context) error {

	log.Trace("end process")
	if context.IsExit() {
		return errors.NewWithCode(errors.New("error in process"), config.ErrorExitedPipeline, "pipeline exited")
	}

	t := common.ParseTask(context)

	if !context.GetBool(keepRedirected, false) && t.Status == model.TaskRedirected {
		if context.Has(model.CONTEXT_TASK_ID) {
			return model.DeleteTask(context.MustGetString(model.CONTEXT_TASK_ID))
		}
		return nil
	}

	if !context.GetBool(keep404, false) && t.Status == model.Task404 {
		if context.Has(model.CONTEXT_TASK_ID) {
			return model.DeleteTask(context.MustGetString(model.CONTEXT_TASK_ID))
		}
		return nil
	}

	if joint.GetBool(isCreate, false) {
		return model.CreateTask(t)
	} else {
		return model.UpdateTask(t)
	}
}
