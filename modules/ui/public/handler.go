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

package public

import (
	"crypto/rand"
	"encoding/base64"
	log "github.com/cihub/seelog"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/google/go-github/github"
	"github.com/infinitbyte/gopa/core/http"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/ui/public/auth"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"net/http"
)

type PublicUI struct {
	api.Handler
	admins *hashset.Set
}

func (h *PublicUI) IsAdmin(user string) bool {
	if h.admins != nil && h.admins.Contains(user) {
		return true
	}
	return false
}

const oauthSession string = "oauth-session"

func (h *PublicUI) AuthHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	session, err := api.GetSessionStore(r, oauthSession)
	session.Values["state"] = state
	session.Values["redirect_url"] = h.Get(r, "redirect_url", "")
	err = session.Save(r, w)
	if err != nil {
		log.Error("error save session")
		http.Redirect(w, r, "/auth/fail/", 500)
		return
	}

	url := oauthCfg.AuthCodeURL(state)
	http.Redirect(w, r, url, 302)
}

func (h *PublicUI) CallbackHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	session, err := api.GetSessionStore(r, oauthSession)
	if err != nil {
		log.Debug(w, "aborted")
		http.Redirect(w, r, "/auth/fail/", 500)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		log.Error("no state match; possible csrf OR cookies not enabled")
		http.Redirect(w, r, "/auth/fail/", 500)
		return
	}

	tkn, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		log.Error("there was an issue getting your token")
		http.Redirect(w, r, "/auth/fail/", 500)
		return
	}

	if !tkn.Valid() {
		log.Error("retreived invalid token")
		http.Redirect(w, r, "/auth/fail/", 500)
		return
	}

	client := github.NewClient(oauthCfg.Client(oauth2.NoContext, tkn))

	user, _, err := client.Users.Get(oauth2.NoContext, "")
	if err != nil {
		log.Error("error getting name")
		http.Redirect(w, r, "/auth/fail/", 500)
		return
	}

	if user != nil {

		//session.Values["provider"] = "github"
		//session.Values["user"] = user.Login
		//session.Values["name"] = user.Name
		//session.Values["accessToken"] = tkn.AccessToken

		role := model.ROLE_GUEST

		if h.IsAdmin(*user.Login) {
			role = model.ROLE_ADMIN
		}

		//session.Values["role"] = role

		log.Debugf("%s(%s) logged in as %s", *user.Name, *user.Login, role)

		err := session.Save(r, w)
		if err != nil {
			log.Error("error save session")
			http.Redirect(w, r, "/auth/fail/", 500)
			return
		}

		api.Login(w, r, *user.Login, role)

		http.Redirect(w, r, "/auth/success/", 302)
		return
	}
	api.Logout(w, r)
	http.Redirect(w, r, "/auth/fail/", 500)
}

func (h PublicUI) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := h.Get(r, "redirect_url", "")
	if url != "" {
		url = "?redirect_url=" + util.UrlEncode(url)
	}
	auth.Login(w, url)
}

func (h PublicUI) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	api.Logout(w, r)
	auth.Logout(w, "/")
}

func (h PublicUI) LoginSuccess(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	session, err := api.GetSessionStore(r, oauthSession)
	if err != nil {
		log.Error("aborted")
		http.Redirect(w, r, "/auth/fail/", 500)
		return
	}

	url := "/"
	if session != nil {
		tmp := session.Values["redirect_url"].(string)
		if tmp != "" {
			url = tmp
		}
	}
	auth.LoginSuccess(w, url)
}

func (h PublicUI) LoginFail(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	auth.LoginFail(w)
}

func (h PublicUI) RedirectHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := h.Get(r, "url", "")
	http.Redirect(w, r, util.UrlDecode(url), 302)
	return
}
