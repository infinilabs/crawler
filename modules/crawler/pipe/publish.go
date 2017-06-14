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
	"github.com/medcl/gopa/core/model"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/modules/config"
)

const Publish JointKey = "publish"

type PublishJoint struct {
}

func (this PublishJoint) Name() string {
	return string(Publish)
}

func (this PublishJoint) Process(c *Context) error {

	//domain := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Domain)
	task := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	//snapshot := c.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	m := md5.Sum([]byte(task.Url))
	id := hex.EncodeToString(m[:])

	data := map[string]interface{}{}

	//data["original_url"] = c.MustGetString(CONTEXT_ORIGINAL_URL)
	//data["url"] = c.MustGetString(CONTEXT_URL)
	//data["host"] = c.MustGetString(CONTEXT_HOST)
	//data["summary"] = c.MustGetString(CONTEXT_PAGE_BODY_PLAIN_TEXT)
	//data["save_path"] = c.MustGetString(CONTEXT_SAVE_PATH)
	//data["save_file"] = c.MustGetString(CONTEXT_SAVE_FILENAME)
	//
	//meta, b := c.GetMap(CONTEXT_PAGE_METADATA)
	//if b {
	//	data["metadata"] = meta
	//}

	//links, b := c.GetMap(CONTEXT_PAGE_LINKS)
	//if b {
	//	maps := []model.PageLink{}
	//	for k, v := range links {
	//		item := model.PageLink{Url: k, Label: v.(string)}
	//		maps = append(maps, item)
	//	}
	//	data["links"] = maps
	//}

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
