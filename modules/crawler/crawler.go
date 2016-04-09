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
	"math/rand"
)

var fetchQuitChannels []*chan bool
var started = false

func Start(config *GopaConfig) {
	if started {
		log.Error("crawler already started, please stop it first.")
	}
	numGoRoutine := config.RuntimeConfig.MaxGoRoutine
	fetchQuitChannels = make([]*chan bool, numGoRoutine)    //shutdownSignal signals for each go routing
	fetchTaskChannels := make([]*chan []byte, numGoRoutine) //fetchTask channels
	fetchOffsets := make([]*RoutingParameter, numGoRoutine) //kafka fetchOffsets
	if config.RuntimeConfig.HttpEnabled {
		go func() {

			//start fetcher
			for i := 0; i < numGoRoutine; i++ {
				quitC := make(chan bool, 1)
				taskC := make(chan []byte)

				fetchQuitChannels[i] = &quitC
				fetchTaskChannels[i] = &taskC
				parameter := new(RoutingParameter)
				parameter.Shard = i
				fetchOffsets[i] = parameter

				fetchTask := new(FetchTask)
				innerTaskConfig := new(InnerTaskConfig)
				innerTaskConfig.RuntimeConfig = config.RuntimeConfig
				innerTaskConfig.MessageChan = &taskC
				innerTaskConfig.QuitChan = &quitC
				innerTaskConfig.Parameter = parameter

				fetchTask.Init(innerTaskConfig)
				go fetchTask.Start()

			}
		}()
	}

	//redistribute pendingFetchUrls to sharded workers
	go func() {
		for {
			url := <-config.Channels.PendingFetchUrl
			if !config.RuntimeConfig.Storage.UrlHasWalked(url) {

				if config.RuntimeConfig.Storage.UrlHasFetched(url) {
					log.Warn("don't hit walk filter but hit fetch filter, also ignore,", string(url))
					config.RuntimeConfig.Storage.AddWalkedUrl(url)
					continue
				}

				randomShard := 0
				if numGoRoutine >= 1 {
					randomShard = rand.Intn(numGoRoutine)
				}
				log.Debug("publish:", string(url), ",shard:", randomShard)
				config.RuntimeConfig.Storage.AddWalkedUrl(url)
				*fetchTaskChannels[randomShard] <- url
			} else {
				log.Trace("hit walk or fetch filter,just ignore,", string(url))
			}
		}
	}()
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
			log.Error("send exit signal to fetch channel: ", i)
		}

		log.Info("crawler success stoped")

		started = false
	} else {
		log.Error("crawler is not started, please start it first.")
	}

	return nil
}
