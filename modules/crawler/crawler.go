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
)

var fetchQuitChannels []*chan bool
var started = false

func Start(env *Env) {
	if started {
		log.Error("crawler already started, please stop it first.")
	}
	numGoRoutine := env.RuntimeConfig.MaxGoRoutine

	fetchQuitChannels = make([]*chan bool, numGoRoutine) //shutdownSignal signals for each go routing
	if env.RuntimeConfig.CrawlerConfig.Enabled {
		go func() {

			//start fetcher
			for i := 0; i < numGoRoutine; i++ {
				log.Trace("start crawler:",i)
				quitC := make(chan bool, 1)
				fetchQuitChannels[i] = &quitC
				go FetchGo(env, &quitC, i)

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
