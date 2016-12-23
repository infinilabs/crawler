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
	"github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/logging"
	"strings"
	"github.com/medcl/gopa/core/types"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/modules/config"
)

type Command struct{
	Env *env.Env
}

func (this *Command) Help(c *WebsocketConnection,a []string) ()  {
	c.WriteMessage([]byte("COMMAND LIST\nseed [url] eg: seed http://elastic.co\nlog [level]  eg: log debug"))
}


func (this *Command) AddSeed(c *WebsocketConnection,a []string) ()  {

	url:=a[1]
	if(len(url)>0){
		queue.Push(config.CheckChannel,types.NewTaskSeed(url,"",0).MustGetBytes())
		c.WriteMessage([]byte("url "+url+" success added to pending fetch queue"))
		return
	}
	c.WriteMessage([]byte("invalid url"))
}

func (this *Command) UpdateLogLevel(c *WebsocketConnection,a []string) ()  {

	level :=a[1]
	if(len(level)>0){
		level:=strings.ToLower(level)
		logging.SetInitLogging(this.Env,level)
		c.WriteMessage([]byte("setting log level to  "+ level ))
		return
	}
	c.WriteMessage([]byte("invalid setting"))
}
