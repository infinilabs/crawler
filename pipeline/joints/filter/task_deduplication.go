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

package filter

import (
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/gopa/model"
)

// TaskDeduplicationJoint is used to find whether the task already in the database
type TaskDeduplicationJoint struct {
}

// Name return task_deduplication
func (joint TaskDeduplicationJoint) Name() string {
	return "task_deduplication"
}

// Process deduplication
func (joint TaskDeduplicationJoint) Process(c *pipeline.Context) error {
	url := c.MustGetString(model.CONTEXT_TASK_URL)

	if global.Env().IsDebug {
		log.Trace("check duplication, ", url)
	}

	items, err := model.GetTaskByField("url", url)
	if err != nil {
		panic(err)
	}
	if len(items) > 0 {

		log.Trace("duplication: ", items)

		//double check
		for _, v := range items {
			if v.Url == url {
				//if v.Status==model.TaskCreated{
				//	return nil
				//}

				msg := fmt.Sprintf("task already exists, %s", url)
				c.Set(model.CONTEXT_TASK_Status, model.TaskDuplicated)
				c.Exit(msg)
				return errors.New(msg)
			} else {
				log.Error(url, " vs ", v.Url, ", mismatch, but said was duplication")
			}
		}

		//panic(errors.New(fmt.Sprintln("duplication false alarm,",url,",",items)))

	}

	return nil
}
