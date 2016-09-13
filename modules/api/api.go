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
	. "github.com/medcl/gopa/core/env"
	. "github.com/medcl/gopa/modules/api/http"
	websocket "github.com/medcl/gopa/modules/api/websocket"
	ui "github.com/medcl/gopa/ui"
)

func internalStart(env *Env) {
	handler := Handler{Env: env}

	mux := http.NewServeMux()
	websocket.InitWebSocket(env)
	mux.HandleFunc("/ws", websocket.ServeWs)

	//Index
	mux.HandleFunc("/", handler.IndexAction)


	//APIs
	mux.HandleFunc("/stats", handler.StatsAction)

	mux.HandleFunc("/task", handler.TaskAction)
	mux.HandleFunc("/task/", handler.TaskAction)
	mux.HandleFunc("/setting/seelog", handler.LoggingSettingAction)
	mux.HandleFunc("/setting/seelog/", handler.LoggingSettingAction)

	//UI pages
	mux.Handle("/ui/", http.FileServer(ui.FS(false)))
	mux.HandleFunc("/ui/boltdb", handler.BoltDBStatusAction)

	log.Info("http server listen at: http://localhost:8001/")
	http.ListenAndServe(":8001", mux)
}

func Start(config *Env) {
	//API server
	//if config.RuntimeConfig.HttpEnabled {
		go func() {
			internalStart(config)
		}()
		log.Debug("api module success started")
	//}
}

func Stop() error {
	return nil
}
