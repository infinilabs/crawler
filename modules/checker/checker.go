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
	"github.com/medcl/gopa/core/filter"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/core/tasks"
	"github.com/medcl/gopa/core/types"
	"github.com/medcl/gopa/modules/config"
	"sync"
	"time"
)

var signalChannel chan bool
var quitChannel chan bool
var started = false

func (this CheckerModule) Name() string {
	return "Checker"
}

func (this CheckerModule) Start(env *Env) {
	if started {
		log.Error("url checker is already started, please stop it first.")
		return
	}
	signalChannel = make(chan bool)
	quitChannel = make(chan bool)
	go runCheckerGo()
	started = true
}

func runCheckerGo() {

	quit := make(chan bool, 1)
	var wg sync.WaitGroup

	go func() {
		for {
			select {
			case <-quit:
				wg.Wait()
				log.Trace("url checker success stoped")
				return
			default:
				{
					startTime := time.Now()
					if !started {
						return
					}
					log.Trace("waiting url to check")

					wg.Add(1)
					defer wg.Done()
					data := queue.Pop(config.CheckChannel)
					url := types.TaskSeedFromBytes(data)

					stats.Increment("checker.url", "finished")

					log.Trace("cheking url:", string(url.Url))

					//TODO 统一 url 格式 , url 目前可能是相对路径
					//checking
					if filter.Exists(config.CheckFilter, []byte(url.Url)) {
						stats.Increment("checker.url", "duplicated")
						log.Debug("url already pushed to fetch queue, ignore :", string(url.Url))
						continue
					}

					//add to filter
					filter.Add(config.CheckFilter, []byte(url.Url))

					task := types.Task{Seed: &url}
					tasks.CreateTask(&task)

					//send to disk queue
					//queue.Push(config.FetchChannel,url.MustGetBytes())

					stats.Increment("checker.url", "get_valid_seed")

					log.Debugf("send url: %s ,depth: %d to  fetch queue", string(url.Url), url.Depth)
					elapsedTime := time.Now().Sub(startTime)
					stats.Timing("checker.url", "time", elapsedTime.Nanoseconds())
				}
			}
		}
	}()

	log.Trace("url checker success started")
	<-signalChannel
	quit <- true
	wg.Wait()
	log.Trace("url checker wait end")
	quitChannel <- true
	log.Trace("url checker quit")
}

func (this CheckerModule) Stop() error {
	if started {
		signalChannel <- true
		started = false
	} else {
		log.Error("url checker is not started")
	}
	return nil
}

type CheckerModule struct {
}
