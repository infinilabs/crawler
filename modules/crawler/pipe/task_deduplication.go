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
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/model"
	api "github.com/infinitbyte/gopa/core/pipeline"
)

// TaskDeduplicationJoint is used to find whether the task already in the database
type TaskDeduplicationJoint struct {
}

// Name return task_deduplication
func (joint TaskDeduplicationJoint) Name() string {
	return "task_deduplication"
}

// Process deduplication
func (joint TaskDeduplicationJoint) Process(c *api.Context) error {
	task := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)

	log.Trace("check duplication, ", task.Url)

	items, err := model.GetTaskByField("url", task.Url)

	if err != nil {
		panic(err)
	}
	if len(items) > 0 {
		msg := fmt.Sprintf("task already exists, %s, %s", task.ID, task.Url)
		c.Exit(msg)
		return errors.New(msg)
	}

	return nil
}
