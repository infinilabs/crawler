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
	"encoding/json"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/api"
	"github.com/medcl/gopa/core/config"
	logging "github.com/medcl/gopa/core/logger"
	"net/http"
)

func (this API) LoggingSettingAction(w http.ResponseWriter, req *http.Request) {
	if req.Method == GET.String() {

		cfg := logging.GetLoggingConfig()
		if cfg != nil {
			this.WriteJson(w, cfg, 200)
		} else {
			this.Error500(w, "config not available")
		}

	} else if req.Method == PUT.String() || req.Method == POST.String() {
		body, err := this.GetRawBody(req)
		if err != nil {
			log.Error(err)
			this.Error500(w, "config update failed")
			return
		}

		configStr := string(body)

		cfg := config.LoggingConfig{}

		err = json.Unmarshal([]byte(configStr), &cfg)

		if err != nil {
			this.Error(w, err)

		}

		log.Info("receive new settings:", configStr)

		logging.UpdateLoggingConfig(&cfg)

		this.WriteJson(w, map[string]interface{}{"ok": true}, http.StatusOK)

	}
}
