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

package url_checker

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/queue"
	. "github.com/medcl/gopa/core/filter"
	"github.com/medcl/gopa/core/stats"
	"path"
	"time"
	"github.com/medcl/gopa/core/types"
	"github.com/medcl/gopa/modules/config"
)

var quitChannel chan bool
var started = false
var filter = LeveldbFilter{}
var filterFileName = "filters/url_fetched"

func (this CheckerModule) Name() string {
	return "Checker"
}

func (this CheckerModule) Start(env *Env) {
	if started {
		log.Error("url checker is already started, please stop it first.")
		return
	}
	quitChannel = make(chan bool)

	err:=filter.Open(path.Join(env.RuntimeConfig.PathConfig.Data, filterFileName))
	if(err!=nil){
		panic(err)
	}

	go runCheckerGo(env, &quitChannel)
	started = true
}

func runCheckerGo(env *Env, quitC *chan bool) {

	go func() {
		for {
			startTime := time.Now()
			if !started {
				return
			}
			log.Trace("waiting url to check")

			data := queue.Pop(config.CheckChannel)
			url:=types.PageTaskFromBytes(data)

			stats.Increment("checker.url", "finished")

			log.Trace("cheking url:", string(url.Url))

			//TODO 统一 url 格式 , url 目前可能是相对路径
			//checking
			if filter.Exists([]byte(url.Url)) {
				log.Debug("url already pushed to fetch queue, ignore :", string(url.Url))
				continue
			}

			//add to filter
			filter.Add([]byte(url.Url))

			//send to disk queue
			queue.Push(config.FetchChannel,url.MustGetBytes())

			stats.Increment("checker.url", "get_valid_seed")

			log.Debugf("send url: %s ,depth: %d to  fetch queue", string(url.Url), url.Depth)
			elapsedTime := time.Now().Sub(startTime)
			stats.Timing("checker.url", "time", elapsedTime.Nanoseconds())
		}
	}()

	log.Trace("url checker success started")

	<-*quitC

}

func (this CheckerModule)Stop() error {
	if started {
		filter.Close()
		quitChannel <- true
		started = false
	} else {
		log.Error("url checker is not started")
	}
	return nil
}

type CheckerModule struct {

}