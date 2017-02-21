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
	"github.com/medcl/gopa/core/global"
	logging "github.com/medcl/gopa/core/logger"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/modules/config"
	"strings"
)

type Command struct {
}

func getHelpMessage() string {
	help := "COMMAND LIST\n" +
		"seed [url] eg: seed http://elastic.co\n" +
		"log [level]  eg: log debug"
	return help
}

func (this Command) Help(c *WebsocketConnection, a []string) {
	c.WritePrivateMessage(getHelpMessage())
}

func (this Command) AddSeed(c *WebsocketConnection, a []string) {
	url := a[1]
	if len(url) > 0 {
		queue.Push(config.CheckChannel, model.NewTaskSeed(url, "", 0, 0).MustGetBytes())
		c.WritePrivateMessage("url " + url + " success added to pending fetch queue")
		return
	}
	c.WritePrivateMessage("invalid url")
}

func (this Command) UpdateLogLevel(c *WebsocketConnection, a []string) {

	level := a[1]
	if len(level) > 0 {
		level := strings.ToLower(level)
		logging.SetLogging(global.Env(), level, "")
		c.WritePrivateMessage("setting log level to  " + level)
		return
	}
	c.WritePrivateMessage("invalid setting")
}

func (this Command) Dispatch(c *WebsocketConnection, a []string) {

	err := queue.Push(config.DispatcherChannel, []byte("go"))
	if err != nil {
		panic(err)
	}
	c.WritePrivateMessage("trigger tasks")
}

func (this Command) GetTask(c *WebsocketConnection, a []string) {

	if len(a) == 2 {
		para1 := a[1]
		task, err := model.GetTask(para1)
		if err != nil {
			c.WritePrivateMessage(err.Error())
		}

		b, err := json.MarshalIndent(task, "", " ")

		c.WritePrivateMessage(string(b))
		c.WritePrivateMessage("get task by taskId," + para1 + "\n")
		return
	}

	if len(a) == 3 {
		para1 := a[1]
		para2 := a[2]
		task, err := model.GetTaskByField(para1, para2)
		if err != nil {
			c.WritePrivateMessage(err.Error())
		}

		b, err := json.MarshalIndent(task, "", " ")

		c.WritePrivateMessage(string(b))

		c.WritePrivateMessage("get task by," + para1 + ", " + para2 + "\n")

		return
	}

	c.WritePrivateMessage("invalid taskId")
}
