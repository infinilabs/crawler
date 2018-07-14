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

package websocket

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/global"
	logging "github.com/infinitbyte/framework/core/logger"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/queue"
	"github.com/infinitbyte/framework/core/ui/websocket"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
	"strings"
)

// AddSeed handle task creation by send a seed, eg:
//SEED http://elastic.co
func AddSeed(c *websocket.WebsocketConnection, a []string) {
	url := a[1]
	if len(url) > 0 {
		context := pipeline.Context{IgnoreBroken: true}
		context.Set(model.CONTEXT_TASK_URL, url)
		err := queue.Push(config.CheckChannel, util.ToJSONBytes(context))
		if err != nil {
			log.Error(err)
		}
		c.WritePrivateMessage("url " + url + " success added to pending fetch queue")
		return
	}
	c.WritePrivateMessage("invalid url")
}

// UpdateLogLevel update the logging level, usually used for debug
func UpdateLogLevel(c *websocket.WebsocketConnection, a []string) {

	level := a[1]
	if len(level) > 0 {
		level := strings.ToLower(level)
		logging.SetLogging(global.Env(), level, "")
		c.WritePrivateMessage("setting log level to  " + level)
		return
	}
	c.WritePrivateMessage("invalid setting")
}

//TODO
// Dispatch just send a dispatch signal to dispatch service
func Dispatch(c *websocket.WebsocketConnection, a []string) {

	err := queue.Push(config.DispatcherChannel, []byte("go"))
	if err != nil {
		panic(err)
	}
	c.WritePrivateMessage("trigger tasks")
}

// GetTask return task information by send task_id, or task field and value, eg:
//GET_TASK host elasticsearch.cn
//GET_TASK 596
func GetTask(c *websocket.WebsocketConnection, a []string) {

	if len(a) == 2 {
		para1 := a[1]
		task, err := model.GetTask(para1)
		if err != nil {
			c.WritePrivateMessage(err.Error())
		}

		b, _ := json.MarshalIndent(task, "", " ")

		c.WritePrivateMessage(string(b))
		c.WritePrivateMessage("get task by taskId," + para1 + "\n")
		return
	}

	if len(a) == 3 {
		para1 := a[1]
		para2 := a[2]
		tasks, err := model.GetTaskByField(para1, para2)
		if err != nil {
			c.WritePrivateMessage(err.Error())
		}

		b, _ := json.MarshalIndent(tasks, "", " ")

		c.WritePrivateMessage(string(b))

		c.WritePrivateMessage("get task by," + para1 + ", " + para2 + "\n")

		return
	}

	c.WritePrivateMessage("invalid taskId")
}
