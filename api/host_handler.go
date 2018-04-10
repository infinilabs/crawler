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
	"encoding/json"
	api "github.com/infinitbyte/framework/core/api"
	"github.com/infinitbyte/framework/core/api/router"
	"github.com/infinitbyte/gopa/model"
	"net/http"
	"strconv"
)

// DeleteHostAction handle host deletion, only support delete by id, eg:
//curl -XDELETE http://127.0.0.1:8001/host/1
func (handler API) DeleteHostAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method == api.DELETE.String() {
		id := ps.ByName("id")
		err := model.DeleteTask(id)
		if err != nil {
			handler.Error(w, err)
		} else {
			handler.WriteJSON(w, map[string]interface{}{"ok": true}, http.StatusOK)
		}
	} else {
		handler.Error404(w)
	}
}

// GetHostAction return host by id, eg:
//curl -XGET http://127.0.0.1:8001/host/1
func (handler API) GetHostAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	task, err := model.GetHost(id)
	if err != nil {
		handler.Error(w, err)
	} else {
		handler.WriteJSON(w, task, http.StatusOK)

	}

}

// GetHostsAction return host list, support parameter: `from`, `size` and `host`, eg:
//curl -XGET http://127.0.0.1:8001/host?from=0&size=10
func (handler API) GetHostsAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	if req.Method == api.GET.String() {

		fr := handler.GetParameter(req, "from")
		si := handler.GetParameter(req, "size")
		host := handler.GetParameter(req, "host")

		from, err := strconv.Atoi(fr)
		if err != nil {
			from = 0
		}
		size, err := strconv.Atoi(si)
		if err != nil {
			size = 10
		}

		total, hosts, err := model.GetHostList(from, size, host)

		newDomains := []model.Host{}
		for _, v := range hosts {

			//total := stats.Stat("host.stats", v.Host+"."+config.STATS_FETCH_TOTAL_COUNT)
			newDomains = append(newDomains, v)
		}

		if err != nil {
			handler.Error(w, err)
		} else {
			handler.WriteJSONListResult(w, total, newDomains, http.StatusOK)
		}
	}
}

func (api API) GetHostConfigsAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fr := api.GetParameter(req, "from")
	si := api.GetParameter(req, "size")
	host := api.GetParameter(req, "host")

	from, err := strconv.Atoi(fr)
	if err != nil {
		from = 0
	}
	size, err := strconv.Atoi(si)
	if err != nil {
		size = 10
	}

	total, configs, err := model.GetHostConfigList(from, size, host)
	if err != nil {
		api.Error(w, err)
	} else {
		api.WriteJSONListResult(w, total, configs, http.StatusOK)
	}
}

func (api API) CreateHostConfigAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	config := model.HostConfig{}

	data, err := api.GetRawBody(req)
	if err != nil {
		api.Error(w, err)
		return
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		api.Error(w, err)
		return
	}

	err = model.CreateHostConfig(&config)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, config, http.StatusOK)
}

func (api API) GetHostConfigAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	cfg, err := model.GetHostConfigByID(id)

	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, cfg, http.StatusOK)
}

func (api API) UpdateHostConfigAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	config := model.HostConfig{}
	id := ps.ByName("id")

	data, err := api.GetRawBody(req)
	if err != nil {
		api.Error(w, err)
		return
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		api.Error(w, err)
		return
	}

	config.ID = id
	err = model.UpdateHostConfig(&config)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, config, http.StatusOK)
}

func (api API) DeleteHostConfigAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	err := model.DeleteHostConfig(id)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, map[string]interface{}{"ok": true}, http.StatusOK)
}
