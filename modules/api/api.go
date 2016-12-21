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
	"github.com/julienschmidt/httprouter"
	. "github.com/medcl/gopa/core/env"
	. "github.com/medcl/gopa/modules/api/http"
	websocket "github.com/medcl/gopa/modules/api/websocket"
	ui "github.com/medcl/gopa/ui"
	_ "net/http/pprof"
)

func internalStart(env *Env) {
	handler := Handler{Env: env}
	router := httprouter.New()

	user := "gopa"
	pass := "gopa"

	mux := http.NewServeMux()
	websocket.InitWebSocket(env)

	mux.HandleFunc("/ws", websocket.ServeWs)

	//Index
	router.GET("/", handler.IndexAction)
	router.GET("/favicon.ico", handler.IndexAction)

	//APIs
	mux.HandleFunc("/stats", handler.StatsAction)

	router.GET("/tasks", handler.TaskAction)
	router.GET("/task/:id", handler.TaskGetAction)
	router.DELETE("/task/:id", BasicAuth(handler.TaskDeleteAction, user, pass))

	mux.HandleFunc("/setting/seelog", handler.LoggingSettingAction)
	mux.HandleFunc("/setting/seelog/", handler.LoggingSettingAction)

	//Snapshot
	mux.HandleFunc("/snapshot/", handler.SnapshotAction)

	//UI pages
	mux.Handle("/ui/", http.FileServer(ui.FS(false)))
	mux.HandleFunc("/ui/boltdb", handler.BoltDBStatusAction)

	mux.Handle("/", router)

	log.Info("http server listen at: http://localhost:8001/")
	http.ListenAndServe(":8001", mux)
}

func BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func (this APIModule) Start(config *Env) {

	this.env = config
	//API server
	go func() {
		internalStart(config)
	}()
	log.Info("api module success started")
}

func (this APIModule) Stop() error {
	return nil
}

type APIModule struct {
	env *Env
}
