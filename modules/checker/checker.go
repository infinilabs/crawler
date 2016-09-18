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
)

var quitChannel chan bool
var started = false

func Start(env *Env) {
	if started {
		log.Error("url checker is already started, please stop it first.")
		return
	}
	quitChannel = make(chan bool)

	go runCheckerGo(env, &quitChannel)
	started = true
}

func runCheckerGo(env *Env, quitC *chan bool) {

	go func() {
		for {
			if !started {
				return
			}
			log.Trace("waiting url to check")

			url := env.Channels.PopUrlToCheck()
			log.Debug("cheking url:", string(url.Url))

			//checking

			//send to disk queue
			env.Channels.PushUrlToFetch(url)
			log.Debugf("send url: %s ,depth: %d to  fetch queue", string(url.Url), url.Depth)
		}
	}()

	log.Trace("url checker success started")

	<-*quitC

	log.Info("url checker success stoped")
}

func Stop() error {
	if started {
		log.Debug("start shutting down url checker")

		quitChannel <- true

		started = false
	} else {
		log.Error("url checker is not started")
	}

	return nil
}
