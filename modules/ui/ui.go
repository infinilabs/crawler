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

package ui

import (
	uis "github.com/infinitbyte/gopa/core/http"
	"github.com/infinitbyte/gopa/modules/ui/admin"
	"github.com/infinitbyte/gopa/modules/ui/user"
	"github.com/infinitbyte/gopa/modules/ui/websocket"
	"github.com/infinitbyte/gopa/static"
	"net/http"

	"crypto/tls"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/logger"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/julienschmidt/httprouter"
	_ "net/http/pprof"
	"path"
	"path/filepath"
)

var router *httprouter.Router
var mux *http.ServeMux

var store = sessions.NewCookieStore([]byte("1c6f2afbccef959ac5c8b81f690c1be7"))

var faviconAction = func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Location", "/static/assets/img/favicon.ico")
	w.WriteHeader(301)
}

func (this UIModule) internalStart(cfg *Config) {

	store.Options = &sessions.Options{
		Domain:   "localhost", //TODO config　http　domain
		Path:     "/",
		MaxAge:   60 * 15,
		Secure:   true,
		HttpOnly: true,
	}

	router = httprouter.New()
	mux = http.NewServeMux()
	websocket.InitWebSocket(global.Env())

	mux.HandleFunc("/ws", websocket.ServeWs)
	mux.Handle("/", router)

	//Index
	router.GET("/favicon.ico", faviconAction)

	//init common
	mux.Handle("/static/", http.FileServer(static.FS(false)))

	//registered handlers
	if uis.RegisteredUIHandler != nil {
		for k, v := range uis.RegisteredUIHandler {
			log.Debug("register custom http handler: ", k)
			mux.Handle(k, v)
		}
	}
	if uis.RegisteredUIFuncHandler != nil {
		for k, v := range uis.RegisteredUIFuncHandler {
			log.Debug("register custom http handler: ", k)
			mux.HandleFunc(k, v)
		}
	}
	if uis.RegisteredUIMethodHandler != nil {
		for k, v := range uis.RegisteredUIMethodHandler {
			for m, n := range v {
				log.Debug("register custom http handler: ", k, " ", m)
				router.Handle(k, m, n)
			}
		}
	}

	address := util.AutoGetAddress(global.Env().SystemConfig.HTTPBinding)

	if global.Env().SystemConfig.TLSEnabled {
		log.Debug("start ssl endpoint")

		certFile := path.Join(global.Env().SystemConfig.PathConfig.Cert, "*c*rt*")
		match, err := filepath.Glob(certFile)
		if err != nil {
			panic(err)
		}
		if len(match) <= 0 {
			panic(errors.New("no cert file found, the file name must end with .crt"))
		}
		certFile = match[0]

		keyFile := path.Join(global.Env().SystemConfig.PathConfig.Cert, "*key*")
		match, err = filepath.Glob(keyFile)
		if err != nil {
			panic(err)
		}
		if len(match) <= 0 {
			panic(errors.New("no key file found, the file name must end with .key"))
		}
		keyFile = match[0]

		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		srv := &http.Server{
			Addr:         address,
			Handler:      context.ClearHandler(mux),
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		log.Info("https server listen at: https://", address)
		err = srv.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			log.Error(err)
			panic(err)
		}

	} else {
		log.Info("http server listen at: http://", address)
		err := http.ListenAndServe(address, context.ClearHandler(mux))
		if err != nil {
			log.Error(err)
			panic(err)
		}
	}

}

type UIModule struct {
}

func (this UIModule) Name() string {
	return "Web"
}

func (this UIModule) Start(cfg *Config) {

	//init admin ui //TODO ui module enable/disable config
	admin.InitUI()

	//init search ui
	user.InitUI()

	//register websocket logger
	logger.RegisterWebsocketHandler(LoggerReceiver)

	go func() {
		this.internalStart(cfg)
	}()

}

func (this UIModule) Stop() error {

	return nil
}

func LoggerReceiver(message string, level log.LogLevel, context log.LogContextInterface) {

	websocket.BroadcastMessage(message)
}
