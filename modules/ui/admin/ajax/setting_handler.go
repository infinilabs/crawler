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

package ajax

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/config"
	api "github.com/infinitbyte/gopa/core/http"
	logging "github.com/infinitbyte/gopa/core/logger"
	"net/http"
)

// Ajax dealing with AJAX request
type Ajax struct {
	api.Handler
}

// LoggingSettingAction is the ajax request to update logging config
func (ajax Ajax) LoggingSettingAction(w http.ResponseWriter, req *http.Request) {
	if req.Method == api.GET.String() {

		cfg := logging.GetLoggingConfig()
		if cfg != nil {
			ajax.WriteJSON(w, cfg, 200)
		} else {
			ajax.Error500(w, "config not available")
		}

	} else if req.Method == api.PUT.String() || req.Method == api.POST.String() {
		body, err := ajax.GetRawBody(req)
		if err != nil {
			log.Error(err)
			ajax.Error500(w, "config update failed")
			return
		}

		configStr := string(body)

		cfg := config.LoggingConfig{}

		err = json.Unmarshal([]byte(configStr), &cfg)

		if err != nil {
			ajax.Error(w, err)

		}

		log.Info("receive new settings:", configStr)

		logging.UpdateLoggingConfig(&cfg)

		ajax.WriteJSON(w, map[string]interface{}{"ok": true}, http.StatusOK)

	}
}
