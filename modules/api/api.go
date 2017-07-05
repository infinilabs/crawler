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

package api

import (
	"net/http"

	"crypto/tls"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	. "github.com/infinitbyte/gopa/core/config"
	. "github.com/infinitbyte/gopa/core/env"
	"github.com/infinitbyte/gopa/core/global"
	apis "github.com/infinitbyte/gopa/core/http"
	"github.com/infinitbyte/gopa/core/util"
	handlers "github.com/infinitbyte/gopa/modules/api/handlers"
	"github.com/julienschmidt/httprouter"
	_ "net/http/pprof"
	"path"
	"path/filepath"
)

var router *httprouter.Router
var mux *http.ServeMux

var store = sessions.NewCookieStore([]byte("1c6f2afbccef959ac5c8b81f690c1be7"))

func (this APIModule) internalStart(env *Env) {

	store.Options = &sessions.Options{
		Domain:   "localhost", //TODO config　http　domain
		Path:     "/",
		MaxAge:   60 * 15,
		Secure:   true,
		HttpOnly: true,
	}

	router = httprouter.New()
	mux = http.NewServeMux()

	mux.Handle("/", router)

	//registered handlers
	if apis.RegisteredAPIHandler != nil {
		for k, v := range apis.RegisteredAPIHandler {
			log.Debug("register custom http handler: ", k)
			mux.Handle(k, v)
		}
	}
	if apis.RegisteredAPIFuncHandler != nil {
		for k, v := range apis.RegisteredAPIFuncHandler {
			log.Debug("register custom http handler: ", k)
			mux.HandleFunc(k, v)
		}
	}
	if apis.RegisteredAPIMethodHandler != nil {
		for k, v := range apis.RegisteredAPIMethodHandler {
			for m, n := range v {
				log.Debug("register custom http handler: ", k, " ", m)
				router.Handle(k, m, n)
			}
		}
	}

	address := util.AutoGetAddress(env.SystemConfig.APIBinding)

	if this.env.SystemConfig.TLSEnabled {
		log.Debug("start ssl endpoint")

		certFile := path.Join(this.env.SystemConfig.PathConfig.Cert, "*c*rt*")
		match, err := filepath.Glob(certFile)
		if err != nil {
			panic(err)
		}
		if len(match) <= 0 {
			panic(errors.New("no cert file found, the file name must end with .crt"))
		}
		certFile = match[0]

		keyFile := path.Join(this.env.SystemConfig.PathConfig.Cert, "*key*")
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

		log.Info("api server listen at: https://", address)
		err = srv.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			log.Error(err)
			panic(err)
		}

	} else {
		log.Info("api server listen at: http://", address)
		err := http.ListenAndServe(address, context.ClearHandler(mux))
		if err != nil {
			log.Error(err)
			panic(err)
		}
	}

}

func (this APIModule) Name() string {
	return "API"
}

func (this APIModule) Start(cfg *Config) {

	this.env = global.Env()
	//API server
	go func() {
		handlers.InitAPI()
		this.internalStart(global.Env())
	}()

}

func (this APIModule) Stop() error {
	return nil
}

type APIModule struct {
	env *Env
}
