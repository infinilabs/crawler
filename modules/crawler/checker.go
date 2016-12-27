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

package crawler

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/filter"
	"github.com/medcl/gopa/core/global"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/modules/config"
	. "github.com/medcl/gopa/modules/crawler/pipe"
	"runtime"
	"sync"
	"time"
)

var signalChannel chan bool
var quitChannel chan bool
var checkerStarted bool
func (this CheckerModule) Name() string {
	return "Checker"
}

func (this CheckerModule) Start(env *Env) {
	if checkerStarted {
		log.Error("url checker is already checkerStarted, please stop it first.")
		return
	}
	signalChannel = make(chan bool)
	quitChannel = make(chan bool)
	go this.runCheckerGo()
	checkerStarted = true
}

func (this CheckerModule) runCheckerGo() {

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

				this.execute(&wg)

			}
		}
	}()

	log.Trace("url checker success checkerStarted")
	<-signalChannel
	quit <- true
	wg.Wait()
	log.Trace("url checker wait end")
	quitChannel <- true
	log.Trace("url checker quit")
}

func  (this CheckerModule)execute(wg *sync.WaitGroup) {
	startTime := time.Now()
	log.Trace("waiting url to check")
	wg.Add(1)

	defer wg.Done()

	data := queue.Pop(config.CheckChannel)
	url := model.TaskSeedFromBytes(data)

	stats.Increment("checker.url", "finished")

	log.Trace("cheking url:", string(url.Url))

	//TODO 统一 url 格式 , url 目前可能是相对路径
	b, err := filter.CheckThenAdd(config.CheckFilter, []byte(url.Url))
	//checking
	if b {
		stats.Increment("checker.url", "duplicated")
		log.Debug("url already pushed to fetch queue, ignore :", string(url.Url))
		return
	}
	if err != nil {
		log.Error(err)
		panic(err)
	}

	task := model.Task{Seed: &url}
	err = model.CreateTask(&task)
	if err != nil {
		panic(err)
	}

	var pipeline *Pipeline
	defer func() {

		if !global.Env().IsDebug {
			if r := recover(); r != nil {
				if _, ok := r.(runtime.Error); ok {
					err := r.(error)
					log.Error("pipeline: ", pipeline.GetID(), ", taskId: ", task.ID, ", ", err)
				}
				log.Error("error in crawler")
			}
		}
	}()
	pipeline = NewPipeline("checker")

	pipeline.Context(&Context{Phrase:config.PhraseChecker}).
		Start(Start{Task: &task}).
		Join(UrlNormalizationJoint{FollowSubDomain: true}).
		Join(UrlFilterJoint{}).
		End(End{}).
		Run()

	//send to disk queue
	if(len(task.Domain)>0){
		stats.Increment("domain.stats", task.Domain+"."+stats.STATS_FETCH_TOTAL_COUNT)
		queue.Push(config.FetchChannel,[]byte(task.ID))
	}else{
		log.Debug("invalid domain, ",url.Url)
	}

	stats.Increment("checker.url", "valid_seed")

	log.Debugf("send url: %s ,depth: %d to  fetch queue", string(url.Url), url.Depth)
	elapsedTime := time.Now().Sub(startTime)
	stats.Timing("checker.url", "time", elapsedTime.Nanoseconds())
}

func (this CheckerModule) Stop() error {
	if checkerStarted {
		signalChannel <- true
		checkerStarted = false
	} else {
		log.Error("url checker is not checkerStarted")
	}
	return nil
}

type CheckerModule struct {
}
