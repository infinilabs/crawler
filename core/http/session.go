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
	log "github.com/cihub/seelog"
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("gopa-secret")) //TODO generate secret or configurable

const sessionName string = "session-name"

func GetSession(w http.ResponseWriter, r *http.Request, key string) (bool, interface{}) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Error(err)
		return false, nil
	}

	v := session.Values[key]
	return v != nil, v
}

func SetSession(w http.ResponseWriter, r *http.Request, key string, value interface{}) bool {
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Error(err)
		return false
	}
	session.Values[key] = value
	session.Save(r, w)
	return true
}

func GetFlash(w http.ResponseWriter, r *http.Request) (bool, []interface{}) {
	log.Trace("get flash")
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Error(err)
		return false, nil
	}
	f := session.Flashes()
	log.Trace(f)
	return f != nil, f
}

func SetFlash(w http.ResponseWriter, r *http.Request, msg string) bool {
	log.Trace("set flash")
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Error(err)
		return false
	}
	session.AddFlash(msg)
	session.Save(r, w)
	return true
}
