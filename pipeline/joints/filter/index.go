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
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/index"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/queue"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
	"github.com/infinitbyte/gopa/pipeline/joints/common"
)

// IndexJoint is used to send snapshot and task info into index
type IndexJoint struct {
}

// Name return index
func (joint IndexJoint) Name() string {
	return "index"
}

// Process wrapper index document and send to queue
func (joint IndexJoint) Process(c *pipeline.Context) error {

	snapshot := c.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	url := c.MustGetString(model.CONTEXT_TASK_URL)
	host := c.MustGetString(model.CONTEXT_TASK_Host)
	m := md5.Sum([]byte(url))
	id := hex.EncodeToString(m[:])

	data := map[string]interface{}{}

	data["host"] = host
	data["task"] = common.ParseTask(c)
	data["snapshot"] = snapshot

	docs := index.IndexDocument{
		Index:  "index",
		ID:     id,
		Source: data,
	}

	bytes, err := json.Marshal(docs)
	if err != nil {
		log.Error(err)
		return err
	}

	err = queue.Push(config.IndexChannel, bytes)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
