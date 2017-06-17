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
	. "github.com/medcl/gopa/core/config"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/global"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/util"
	"github.com/medcl/gopa/modules/config"
	. "github.com/medcl/gopa/modules/crawler/config"
	. "github.com/medcl/gopa/modules/crawler/pipe"
	"runtime"
	"time"
)

var signalChannels []*chan bool
var quitChannels []*chan bool

var crawlerStarted bool

func (this CrawlerModule) Name() string {
	return "Crawler"
}

func (this CrawlerModule) Start(cfg *Config) {

	config := GetDefaultCrawlerConfig()
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
	quitChannels = make([]*chan bool, numGoRoutine)
	if true {
		//start fetcher
		for i := 0; i < numGoRoutine; i++ {
			log.Trace("start crawler:", i)
			signalC := make(chan bool, 1)
			quitC := make(chan bool, 1)
			signalChannels[i] = &signalC
			quitChannels[i] = &quitC
			go this.runPipeline(global.Env(), &signalC, &quitC, i)

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

		//waiting for quit
		for i, item := range quitChannels {
			log.Debug("get final exit signal from fetch channel: ", i)
			if item != nil {
				<-*item
			}
		}

	} else {
		log.Error("crawler is not crawlerStarted, please start it first.")
	}

	return nil
}

func (this CrawlerModule) runPipeline(env *Env, signalC *chan bool, quitC *chan bool, shard int) {

	quit := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				log.Trace("waiting url to fetch, shard:", shard)
				err, taskID := queue.Pop(config.FetchChannel)
				if err != nil {
					log.Trace(err)
					continue
				}
				id := string(taskID)
				log.Trace("shard:", shard, ",task received:", id)
				this.execute(id, env)
				log.Trace("shard:", shard, ",task finished:", id)

			}
		}
	}()
	log.Trace("crawler Started, shard:", shard)
	<-*signalC
	log.Trace("crawler gonna exit, waiting task to finish, shard:", shard)
	quit <- true
	log.Trace("crawler finished, shard:", shard)
	*quitC <- true
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
				log.Debug("error in crawler,", util.ToJson(pipeline.GetContext(), true))
			}
		}
	}()

	pipeline = NewPipeline("crawler")

	init := InitTaskJoint{}
	init.Set(TaskID, taskId)

	pipeline.Context(&Context{Phrase: config.PhraseCrawler}).
		Start(init).
		Join(UrlNormalizationJoint{FollowAllDomain: true, FollowSubDomain: true}).
		Join(LoadMetadataJoint{}).
		Join(IgnoreTimeoutJoint{IgnoreTimeoutAfterCount: 10}).
		Join(FetchJoint{}).
		Join(ParsePageJoint{DispatchLinks: true, MaxDepth: 30, MaxBreadth: 3}).
		Join(HtmlToTextJoint{MergeWhitespace: false}).
		Join(HashJoint{Simhash: false}).
		//Join(SaveSnapshotToFileSystemJoint{}).
		Join(SaveSnapshotToDBJoint{CompressBody: true, Bucket: "Global"}).
		Join(PublishJoint{}).
		End(SaveTaskJoint{}).
		Run()

	if this.config.FetchThresholdInMs > 0 {
		log.Debug("sleep ", this.config.FetchThresholdInMs, "ms to control crawling speed")
		time.Sleep(time.Duration(this.config.FetchThresholdInMs) * time.Millisecond)
		log.Debug("wake up now,continue crawing")
	}

	log.Trace("end crawler")
}

type CrawlerModule struct {
	config *CrawlerConfig
}
