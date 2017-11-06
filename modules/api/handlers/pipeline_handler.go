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
	"github.com/infinitbyte/gopa/core/model"
	"github.com/julienschmidt/httprouter"
	"net/http"

	"encoding/json"
)

func (api API) handleGetPipelineJointsRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	joints := model.GetAllRegisteredJoints()
	api.WriteJSON(w, joints, http.StatusOK)
}

func (api API) handleCreatePipelineConfigRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	config := model.PipelineConfig{}

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

	err = model.CreatePipelineConfig(&config)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, config, http.StatusOK)
}

func (api API) handleGetPipelineConfigRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")
	cfg, err := model.GetPipelineConfig(id)
	if err != nil {
		api.Error(w, err)
	} else {
		api.WriteJSON(w, cfg, http.StatusOK)
	}
}

func (api API) handleUpdatePipelineConfigRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	config := model.PipelineConfig{}
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

	err = model.UpdatePipelineConfig(id, &config)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, config, http.StatusOK)
}

func (api API) handleDeletePipelineConfigRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	err := model.DeletePipelineConfig(id)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, map[string]interface{}{"ok": true}, http.StatusOK)
}
