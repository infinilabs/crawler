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
	"github.com/infinitbyte/gopa/core/global"
	"net/http"
	"sync"
)

const sessionName string = "gopa-session"

func GetSessionStore(r *http.Request, key string) (*sessions.Session, error) {
	return getStore().Get(r, key)
}

// GetSession return session by session key
func GetSession(w http.ResponseWriter, r *http.Request, key string) (bool, interface{}) {
	s := getStore()
	session, err := s.Get(r, sessionName)
	if err != nil {
		log.Error(err)
		return false, nil
	}

	v := session.Values[key]
	return v != nil, v
}

// SetSession set session by session key and session value
func SetSession(w http.ResponseWriter, r *http.Request, key string, value interface{}) bool {
	s := getStore()
	session, err := s.Get(r, sessionName)
	if err != nil {
		log.Error(err)
		return false
	}
	session.Values[key] = value
	err = session.Save(r, w)
	if err != nil {
		log.Error(err)
	}
	return true
}

// GetFlash get flash value
func GetFlash(w http.ResponseWriter, r *http.Request) (bool, []interface{}) {
	log.Trace("get flash")
	session, err := getStore().Get(r, sessionName)
	if err != nil {
		log.Error(err)
		return false, nil
	}
	f := session.Flashes()
	log.Trace(f)
	return f != nil, f
}

// SetFlash set flash value
func SetFlash(w http.ResponseWriter, r *http.Request, msg string) bool {
	log.Trace("set flash")
	session, err := getStore().Get(r, sessionName)
	if err != nil {
		log.Error(err)
		return false
	}
	session.AddFlash(msg)
	session.Save(r, w)
	return true
}

var store *sessions.CookieStore
var lock sync.Mutex

func getStore() *sessions.CookieStore {
	lock.Lock()
	defer lock.Unlock()

	if store != nil {
		return store
	}

	secret := global.Env().SystemConfig.CookieSecret
	if secret == "" {
		log.Trace("use default cookie secret")
		store = sessions.NewCookieStore([]byte("GOPA-SECRET"))
	} else {
		log.Trace("get cookie secret from config,", secret)
		store = sessions.NewCookieStore([]byte(secret))
	}

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 1,
		HttpOnly: true,
	}

	return store

}
