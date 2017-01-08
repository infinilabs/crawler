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
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/modules/config"
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

func (this CrawlerModule) Start(env *Env) {
	if crawlerStarted {
		log.Error("crawler already crawlerStarted, please stop it first.")
		return
	}

	numGoRoutine := env.RuntimeConfig.MaxGoRoutine
	signalChannels = make([]*chan bool, numGoRoutine)
	quitChannels = make([]*chan bool, numGoRoutine)
	if env.RuntimeConfig.CrawlerConfig.Enabled {
		//start fetcher
		for i := 0; i < numGoRoutine; i++ {
			log.Trace("start crawler:", i)
			signalC := make(chan bool, 1)
			quitC := make(chan bool, 1)
			signalChannels[i] = &signalC
			quitChannels[i] = &quitC
			go this.runPipeline(env, &signalC, &quitC, i)

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
				log.Trace("shard:", shard, ",task received:", string(taskID))
				this.execute(string(taskID), env)
				log.Trace("shard:", shard, ",task finished:", string(taskID))

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
				if _, ok := r.(runtime.Error); ok {
					err := r.(error)
					log.Error("pipeline: ", pipeline.GetID(), ", taskId: ", taskId, ", ", err)
				}
				log.Error("error in crawler")
			}
		}
	}()

	pipeline = NewPipeline("crawler")

	pipeline.Context(&Context{Phrase: config.PhraseCrawler}).
		Start(InitTask{ID: &taskId}).
		Join(UrlNormalizationJoint{FollowSubDomain: true}).
		Join(UrlExtFilterJoint{}).
		Join(LoadMetadataJoint{}).
		Join(IgnoreTimeoutJoint{IgnoreTimeoutAfterCount: 100}).
		Join(FetchJoint{}).
		Join(ParserJoint{DispatchLinks: true, MaxDepth: 3}).
		Join(HtmlToTextJoint{MergeWhitespace: true}).
		//Join(SaveToFileSystemJoint{}).
		Join(SaveToDBJoint{CompressBody: true}).
		//Join(PublishJoint{}).
		End(SaveTask{}).
		Run()

	if env.RuntimeConfig.TaskConfig.FetchDelayThreshold > 0 {
		log.Debug("sleep ", env.RuntimeConfig.TaskConfig.FetchDelayThreshold, "ms to control crawling speed")
		time.Sleep(time.Duration(env.RuntimeConfig.TaskConfig.FetchDelayThreshold) * time.Millisecond)
		log.Debug("wake up now,continue crawing")
	}

	log.Trace("end crawler")
}

type CrawlerModule struct {
}
