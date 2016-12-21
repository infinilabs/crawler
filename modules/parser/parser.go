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

package parser

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/env"
)

var quitChannels []*chan bool
var started = false

type ParserModule  struct {

}

func (this ParserModule) Name() string {
	return "Parser"
}

func (this ParserModule)Start(env *Env) {
	if started {
		log.Error("parser is already started, please stop it first.")
	}
	parseQuitChannels := make([]*chan bool, 2) //shutdownSignal signals for each go routing
	c2 := make(chan bool, 1)

	parseQuitChannels[0] = &c2

	//start local saved file parser
	if env.RuntimeConfig.ParserConfig.ParseUrlsFromSavedFileLog {
		go ParseGo(env, &c2)
		started = true
	}
}

func (this ParserModule) Stop() error {
	if started {
		log.Debug("start shutting down parser")

		for i, item := range quitChannels {
			if item != nil {
				*item <- true
			}
			log.Error("send exit signal to parser channel: ", i)
		}
		started = false
	} else {
		log.Error("parser is not started")
	}

	return nil
}
