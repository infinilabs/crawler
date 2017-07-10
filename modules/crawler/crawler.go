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
	. "github.com/infinitbyte/gopa/core/config"
	. "github.com/infinitbyte/gopa/core/env"
	"github.com/infinitbyte/gopa/core/global"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	. "github.com/infinitbyte/gopa/modules/crawler/config"
	. "github.com/infinitbyte/gopa/modules/crawler/pipe"
	"runtime"
	"time"
)

var signalChannels []*chan bool

var crawlerStarted bool

func (this CrawlerModule) Name() string {
	return "Crawler"
}

func (this CrawlerModule) Start(cfg *Config) {

	config := GetDefaultTaskConfig()
	cfg.Unpack(&config)
	this.config = &config

	//TODO
	InitJoints()

	if crawlerStarted {
		log.Error("crawler already crawlerStarted, please stop it first.")
		return
	}

	numGoRoutine := config.MaxGoRoutine
	signalChannels = make([]*chan bool, numGoRoutine)
	if true {
		//start fetcher
		for i := 0; i < numGoRoutine; i++ {
			log.Trace("start crawler:", i)
			signalC := make(chan bool, 1)
			signalChannels[i] = &signalC
			go this.runPipeline(global.Env(), &signalC, i)

		}
	} else {
		log.Info("crawler currently not enabled")
		return
	}

	crawlerStarted = true
}

func (this CrawlerModule) Stop() error {
	if crawlerStarted {
		crawlerStarted = false
		log.Debug("start shutting down crawler")
		for i, item := range signalChannels {
			if item != nil {
				*item <- true
			}
			log.Debug("send exit signal to fetch channel: ", i)
		}

	} else {
		log.Error("crawler is not crawlerStarted, please start it first.")
	}

	return nil
}

func (this CrawlerModule) runPipeline(env *Env, signalC *chan bool, shard int) {

	var taskID []byte
	for {
		select {
		case <-*signalC:
			return
		case taskID = <-queue.ReadChan(config.FetchChannel):
			stats.Increment("queue."+string(config.FetchChannel), "pop")
			id := string(taskID)
			log.Trace("shard:", shard, ",task received:", id)
			this.execute(id, env)
			log.Trace("shard:", shard, ",task finished:", id)
		}
	}
	log.Trace("crawler exit, shard:", shard)
}

func (this CrawlerModule) execute(taskId string, env *Env) {
	var pipeline *Pipeline
	defer func() {
		if !env.IsDebug {
			if r := recover(); r != nil {
				if e, ok := r.(runtime.Error); ok {
					log.Error("pipeline: ", pipeline.GetID(), ", taskId: ", taskId, ", ", util.GetRuntimeErrorMessage(e))
				}
				log.Debug("error in crawler,", util.ToJson(r, true), util.ToJson(pipeline.GetContext(), true))
			}
		}
	}()

	context := &Context{Phrase: config.PhraseCrawler}
	context.Set(CONTEXT_TASK_ID, taskId)

	if this.config.DefaultPipelineConfig == nil {
		panic("default pipeline config can't be null")
	}

	pipeline = NewPipelineFromConfig(this.config.DefaultPipelineConfig)
	pipeline.Context(context)
	pipeline.Run()

	if this.config.FetchThresholdInMs > 0 {
		log.Debug("sleep ", this.config.FetchThresholdInMs, "ms to control crawling speed")
		time.Sleep(time.Duration(this.config.FetchThresholdInMs) * time.Millisecond)
		log.Debug("wake up now,continue crawing")
	}

	log.Trace("end crawler")
}

type CrawlerModule struct {
	config *TaskConfig
}
