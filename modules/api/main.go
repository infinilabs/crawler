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

package http

import (
	"net/http"

	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	. "github.com/medcl/gopa/modules/api/handler"
)

func internalStart(runtimeConfig *RuntimeConfig) {

	handler := Handler{Config: runtimeConfig}

	http.HandleFunc("/", handler.IndexAction)
	http.HandleFunc("/stats", handler.StatsAction)

	log.Info("http server listen at: http://localhost:8001/")
	http.ListenAndServe(":8001", nil)
}

func Start(runtimeConfig *RuntimeConfig) {
	//API server
	if runtimeConfig.HttpEnabled {
		go func() {
			internalStart(runtimeConfig)
		}()
		log.Debug("api module success started")
	}
}

func Stop() error {
	return nil
}
