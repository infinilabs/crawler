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
	"time"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/env"
	. "github.com/medcl/gopa/core/pipeline"
	. "github.com/medcl/gopa/modules/crawler/pipe"
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
			log.Debug("ready to receive url")
			url := string(<-env.Channels.PendingFetchUrl)
			log.Debug("shard:", shard, ",url received:", url)

			//if !env.RuntimeConfig.Storage.UrlHasFetched([]byte(url)) {

			//log.Info("shard:", shard, ",url cool,start fetching:", url)

			pipeline := Pipeline{}
			pipeline.Context(&Context{Env: env}).
				Start(FetchJoint{Url: url}).
				//Join(ParserJoint{}).
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
			//} else {
			//	log.Debug("shard:", shard, ",url received,but already fetched,skip: ", url)
			//}

		}
	}()

	log.Trace("fetch task started.shard:", shard)

	<-*quitC

	log.Trace("fetch task exit.shard:", shard)

}
