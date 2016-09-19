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
	"github.com/medcl/gopa/core/types"
	. "github.com/medcl/gopa/modules/crawler/pipe"
	"time"
	"runtime"
)

var fetchQuitChannels []*chan bool
var started = false

func Start(env *Env) {
	if started {
		log.Error("crawler already started, please stop it first.")
	}
	numGoRoutine := env.RuntimeConfig.MaxGoRoutine
	//shutdownSignal signals for each go routing
	fetchQuitChannels = make([]*chan bool, numGoRoutine)
	if env.RuntimeConfig.CrawlerConfig.Enabled {
		go func() {

			//start fetcher
			for i := 0; i < numGoRoutine; i++ {
				log.Trace("start crawler:", i)
				quitC := make(chan bool, 1)
				fetchQuitChannels[i] = &quitC
				go RunPipeline(env, &quitC, i)

			}
		}()
	}

	started = true
	log.Info("crawler success started")
}

func Stop() error {
	if started {
		log.Debug("start shutting down crawler")
		for i, item := range fetchQuitChannels {
			if item != nil {
				*item <- true
			}
			log.Debug("send exit signal to fetch channel: ", i)
		}

		log.Info("crawler success stoped")
		started = false
	} else {
		log.Error("crawler is not started, please start it first.")
	}

	return nil
}

func RunPipeline(env *Env, quitC *chan bool, shard int) {

	go func() {
		for {
			log.Trace("waiting url to fetch")

			url, err := env.Channels.PopUrlToFetch()
			if err != nil {
				continue
			}
			urlStr := string(url.Url)
			log.Debug("shard:", shard, ",url received:", urlStr)

			execute(url, env)
		}
	}()

	log.Trace("fetch task started.shard:", shard)

	<-*quitC

	log.Trace("fetch task exit.shard:", shard)

}

func execute(task types.PageTask, env *Env) {

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				err := r.(error)
				log.Error(task.Url, " , ", err)
			}
		}
	}()

	pipeline := Pipeline{}
	pipeline.Context(&Context{Env: env}).
		Start(UrlSource{Url: task.Url, Depth: task.Depth, Reference: task.Reference}).
		Join(UrlNormalizationJoint{FollowSubDomain: true}).
		Join(UrlFilterJoint{}).
		Join(LoadMetadataJoint{}).
		Join(IgnoreTimeoutJoint{IgnoreTimeoutAfterCount: 100}).
		Join(FetchJoint{}).
		Join(ParserJoint{DispatchLinks: true}).
		//Join(SaveToFileSystemJoint{}).
		Join(SaveToDBJoint{}).
		//Join(PublishJoint{}).
		End().
		Run()

	if env.RuntimeConfig.TaskConfig.FetchDelayThreshold > 0 {
		log.Debug("sleep ", env.RuntimeConfig.TaskConfig.FetchDelayThreshold, "ms to control crawling speed")
		time.Sleep(time.Duration(env.RuntimeConfig.TaskConfig.FetchDelayThreshold) * time.Millisecond)
		log.Debug("wake up now,continue crawing")
	}
}
