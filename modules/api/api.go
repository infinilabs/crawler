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
	apis "github.com/medcl/gopa/core/api"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/logger"
	"github.com/medcl/gopa/core/util"
	. "github.com/medcl/gopa/modules/api/http"
	"github.com/medcl/gopa/modules/api/websocket"
	_ "net/http/pprof"
	"path"
	"path/filepath"
	"errors"
	"crypto/tls"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
)

var router *httprouter.Router
var mux *http.ServeMux

var store = sessions.NewCookieStore([]byte("1c6f2afbccef959ac5c8b81f690c1be7"))

func (this APIModule) internalStart(env *Env) {

	store.Options = &sessions.Options{
		Domain:     "localhost", //TODO config　http　domain
		Path:       "/",
		MaxAge:     60 * 15,
		Secure:     true,
		HttpOnly:   true,
	}

	handler := API{}
	router = httprouter.New()

	user := "gopa"
	pass := "gopa"

	mux = http.NewServeMux()
	websocket.InitWebSocket(env)

	mux.HandleFunc("/ws", websocket.ServeWs)
	mux.Handle("/", router)

	//Index
	router.GET("/", handler.IndexAction)
	router.GET("/favicon.ico", handler.IndexAction)

	//APIs
	mux.HandleFunc("/stats", handler.StatsAction)

	router.POST("/task/", handler.TaskAction)
	router.GET("/task", handler.TaskAction)
	router.GET("/task/:id", handler.TaskGetAction)
	router.DELETE("/task/:id", BasicAuth(handler.TaskDeleteAction, user, pass))

	router.GET("/domain", handler.DomainAction)
	router.GET("/domain/:id", handler.DomainGetAction)
	router.DELETE("/domain/:id", BasicAuth(handler.DomainDeleteAction, user, pass))


	mux.HandleFunc("/setting/logger", handler.LoggingSettingAction)
	mux.HandleFunc("/setting/logger/", handler.LoggingSettingAction)

	//Snapshot
	mux.HandleFunc("/snapshot/", handler.SnapshotAction)

	//registered handlers
	if apis.RegisteredHandler != nil {
		for k, v := range apis.RegisteredHandler {
			log.Debug("register custom http handler: ", k)
			mux.Handle(k, v)
		}
	}
	if apis.RegisteredFuncHandler != nil {
		for k, v := range apis.RegisteredFuncHandler {
			log.Debug("register custom http handler: ", k)
			mux.HandleFunc(k, v)
		}
	}
	if apis.RegisteredMethodHandler != nil {
		for k, v := range apis.RegisteredMethodHandler {
			for m, n := range v {
				log.Debug("register custom http handler: ", k," ",m)
				router.Handle(k,m,n)
			}
		}
	}

	address := util.AutoGetAddress(env.SystemConfig.HttpBinding)

	if(len(this.env.SystemConfig.CertPath)>0){
		log.Debug("start ssl endpoint")

		certFile:=path.Join(this.env.SystemConfig.CertPath,"*c*rt*")
		match,err:=filepath.Glob(certFile)
		if(err!=nil){
			panic(err)
		}
		if(len(match)<=0){
			panic(errors.New("no cert file found, the file name must end with .crt"))
		}
		certFile=match[0]


		keyFile:=path.Join(this.env.SystemConfig.CertPath,"*key*")
		match,err=filepath.Glob(keyFile)
		if(err!=nil){
			panic(err)
		}
		if(len(match)<=0){
			panic(errors.New("no key file found, the file name must end with .key"))
		}
		keyFile=match[0]


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
		err=srv.ListenAndServeTLS(certFile, keyFile)
		if(err!=nil){
			log.Error(err)
			panic(err)
		}

	}else{
		log.Info("http server listen at: http://", address)
		err:=http.ListenAndServe(address, context.ClearHandler(mux))
		if(err!=nil){
			log.Error(err)
			panic(err)
		}
	}

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

func (this APIModule) Name() string {
	return "API"
}

func LoggerReceiver(message string, level log.LogLevel, context log.LogContextInterface) {

	websocket.BroadcastMessage(message)
}

func (this APIModule) Start(config *Env) {

	this.env = config
	//API server
	go func() {
		this.internalStart(config)
	}()

	logger.RegisterWebsocketHandler(LoggerReceiver)

}

func (this APIModule) Stop() error {
	return nil
}

type APIModule struct {
	env *Env
}
