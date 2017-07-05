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
	"fmt"
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/model"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	. "github.com/infinitbyte/gopa/modules/crawler/pipe"
	"runtime"
	"time"
)

var signalChannel chan bool

var checkerStarted bool

func (this CheckerModule) Name() string {
	return "Checker"
}

func (this CheckerModule) Start(cfg *Config) {
	if checkerStarted {
		log.Error("url checker is already checkerStarted, please stop it first.")
		return
	}
	signalChannel = make(chan bool)
	go this.runCheckerGo()
	checkerStarted = true
	log.Trace("Checker started")
}

func (this CheckerModule) runCheckerGo() {

	var data []byte
	for {
		select {
		case data = <-queue.ReadChan(config.CheckChannel):
			stats.Increment("queue."+string(config.CheckChannel), "pop")
			this.execute(data)
		case <-signalChannel:
			fmt.Println("url checker quit")
		}

	}
}

func (this CheckerModule) execute(data []byte) {
	startTime := time.Now()
	seed := model.TaskSeedFromBytes(data)

	task := model.Task{}
	task.OriginalUrl = seed.Url
	task.Url = seed.Url
	task.Reference = seed.Reference
	task.Depth = seed.Depth
	task.Breadth = seed.Breadth

	pipeline := this.runPipe(global.Env().IsDebug, &task)

	//send to disk queue
	if len(task.Host) > 0 && !pipeline.GetContext().IsErrorExit() && !pipeline.GetContext().IsBreak() {
		stats.Increment("domain.stats", task.Host+"."+config.STATS_FETCH_TOTAL_COUNT)

		err := model.IncrementDomainLinkCount(task.Host)
		if err != nil {
			log.Error(err)
		}
		log.Trace("load host settings, ", task.Host)

		queue.Push(config.FetchChannel, []byte(task.ID))

		stats.Increment("checker.url", "valid_seed")

		log.Debugf("send url: %s ,depth: %d, breadth: %d, to fetch queue", string(seed.Url), seed.Depth, seed.Breadth)
		elapsedTime := time.Now().Sub(startTime)
		stats.Timing("checker.url", "time", elapsedTime.Nanoseconds())
	} else {
		log.Debug("ignored url, ", seed.Url)
	}

	stats.Increment("checker.url", "finished")

}

func (this CheckerModule) runPipe(debug bool, task *model.Task) *Pipeline {
	var pipeline *Pipeline
	defer func() {

		if !debug {
			if r := recover(); r != nil {
				if e, ok := r.(runtime.Error); ok {
					log.Error("pipeline: ", pipeline.GetID(), ", taskId: ", task.ID, ", ", util.GetRuntimeErrorMessage(e))
				}
				log.Error("error in checker")
			}
		}
	}()
	pipeline = NewPipeline("checker")

	context := &Context{Phrase: config.PhraseChecker, IgnoreBroken: true}
	pipeline.Context(context).
		Start(InitTaskJoint{Task: task}).
		Join(UrlNormalizationJoint{FollowAllDomain: false, FollowSubDomain: true}).
		Join(UrlExtFilterJoint{}).
		Join(UrlCheckFilterJoint{}).
		End(SaveTaskJoint{IsCreate: true}).
		Run()

	return pipeline
}

func (this CheckerModule) Stop() error {
	log.Trace("start stop checker")

	if checkerStarted {
		log.Trace("send signal to checker")
		signalChannel <- true
		log.Trace("close queue checker")
		checkerStarted = false
		log.Debug("closed queue checker")

	} else {
		log.Error("url checker is not checkerStarted")
	}
	log.Debug("done stop checker")

	return nil
}

type CheckerModule struct {
}
