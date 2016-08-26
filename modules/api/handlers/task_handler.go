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
	logger "github.com/cihub/seelog"
	_ "github.com/jmoiron/jsonq"
	"net/http"
)

func (this *Handler) TaskAction(w http.ResponseWriter, req *http.Request) {

	if req.Method == POST.String() {
		jsonq, err := this.GetJson(req)
		if err != nil {
			logger.Error(err)
		}

		seed, err := jsonq.String("seed")
		if err != nil {
			logger.Error(err)
		}

		logger.Info("receive new seed:", seed)

		this.Env.Channels.PendingFetchUrl <- []byte(seed)

		this.WriteJson(w, map[string]interface{}{"ok": true}, http.StatusOK)
	} else {
		this.error404(w)
	}
}
