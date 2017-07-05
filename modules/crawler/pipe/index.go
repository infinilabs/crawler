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
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/model"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
)

const Publish JointKey = "index"

type IndexJoint struct {
}

func (this IndexJoint) Name() string {
	return string(Publish)
}

func (this IndexJoint) Process(c *Context) error {

	task := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := c.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	m := md5.Sum([]byte(task.Url))
	id := hex.EncodeToString(m[:])

	data := map[string]interface{}{}

	data["domain"] = util.MapStr{
		"host": task.Host,
	}

	data["task"] = task
	data["snapshot"] = snapshot

	docs := model.IndexDocument{
		Index:  "gopa",
		Type:   "doc",
		Id:     id,
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
