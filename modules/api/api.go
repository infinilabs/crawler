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
	//"github.com/elazarl/go-bindata-assetfs"
	. "github.com/medcl/gopa/core/config"
	. "github.com/medcl/gopa/modules/api/handlers"
	websocket "github.com/medcl/gopa/modules/api/websocket"
	ui "github.com/medcl/gopa/ui"
)

func internalStart(config *GopaConfig) {
	handler := Handler{Config: config}

	websocket.InitWebSocket()

	http.HandleFunc("/", handler.IndexAction)
	http.HandleFunc("/stats", handler.StatsAction)
	http.Handle("/ui/", http.FileServer(ui.FS(false)))

	http.HandleFunc("/task/", handler.TaskAction)

	http.HandleFunc("/ws", websocket.ServeWs)

	log.Info("http server listen at: http://localhost:8001/")
	http.ListenAndServe(":8001", nil)
}

func Start(config *GopaConfig) {
	//API server
	if config.RuntimeConfig.HttpEnabled {
		go func() {
			internalStart(config)
		}()
		log.Debug("api module success started")
	}
}

func Stop() error {
	return nil
}
