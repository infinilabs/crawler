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
	"github.com/infinitbyte/framework/core/api"
	"github.com/infinitbyte/framework/core/api/router"
	"net/http"
	"time"
)

func (handler API) handleUserLoginRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	b, v := api.GetSession(w, req, "key")
	if !b {
		api.SetSession(w, req, "key", "hello "+time.Now().String())
		api.SetFlash(w, req, "user logged in")
	}

	b, v = api.GetFlash(w, req)
	if b {
		handler.WriteJSON(w, v, 200)
		return
	}

	handler.WriteJSON(w, v, 200)

	return
}
