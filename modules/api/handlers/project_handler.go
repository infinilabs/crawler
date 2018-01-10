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

package http

import (
	"encoding/json"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (api API) GetProjectsAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fr := api.GetParameter(req, "from")
	si := api.GetParameter(req, "size")

	from, err := strconv.Atoi(fr)
	if err != nil {
		from = 0
	}
	size, err := strconv.Atoi(si)
	if err != nil {
		size = 10
	}

	total, configs, err := model.GetProjectList(from, size)
	if err != nil {
		api.Error(w, err)
	} else {
		api.WriteJSONListResult(w, total, configs, http.StatusOK)
	}
}

func (api API) CreateProjectAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	config := model.Project{}

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

	err = model.CreateProject(&config)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, config, http.StatusOK)
}

func (api API) GetProjectAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	cfg, err := model.GetProject(id)

	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, cfg, http.StatusOK)
}

func (api API) UpdateProjectAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	config := model.Project{}
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
	err = model.UpdateProject(&config)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, config, http.StatusOK)
}

func (api API) DeleteProjectAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	err := model.DeleteProject(id)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, map[string]interface{}{"ok": true}, http.StatusOK)
}
