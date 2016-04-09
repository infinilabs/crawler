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

package handler

import (
	log "github.com/cihub/seelog"
	logging "github.com/medcl/gopa/core/logging"
	"net/http"
)

func (this *Handler) LoggingSettingAction(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {

		str := logging.GetLoggingConfig(this.Env)
		if len(str) > 0 {
			this.Write(w, []byte(str))
		} else {
			this.error500(w, "empty setting")
		}

	} else if req.Method == "PUT" || req.Method == "POST" {
		body, err := this.GetRawBody(req)
		if err != nil {
			log.Error(err)
			this.error500(w, "config replace failed")
			return
		}

		configStr := string(body)

		log.Info("receive new settings:", configStr)

		logging.ReplaceConfig(this.Env, configStr)

		this.WriteJson(w, map[string]interface{}{"ok": true}, http.StatusOK)
	}
}
