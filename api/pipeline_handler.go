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
	"github.com/infinitbyte/framework/core/api/router"
	"github.com/infinitbyte/framework/core/pipeline"
	"net/http"

	"encoding/json"
	"strconv"
)

func (api API) handleGetPipelineJointsRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	joints := pipeline.GetAllRegisteredJoints()
	api.WriteJSON(w, joints, http.StatusOK)
}

func (api API) handleGetPipelineConfigsRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

	total, configs, err := pipeline.GetPipelineList(from, size)
	if err != nil {
		api.Error(w, err)
	} else {
		api.WriteJSONListResult(w, total, configs, http.StatusOK)
	}
}

func (api API) handleCreatePipelineConfigRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	config := pipeline.PipelineConfig{}

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

	err = pipeline.CreatePipelineConfig(&config)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, config, http.StatusOK)
}

func (api API) handleGetPipelineConfigRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")
	cfg, err := pipeline.GetPipelineConfig(id)
	if err != nil {
		api.Error(w, err)
	} else {
		api.WriteJSON(w, cfg, http.StatusOK)
	}
}

func (api API) handleUpdatePipelineConfigRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	config := pipeline.PipelineConfig{}
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

	err = pipeline.UpdatePipelineConfig(id, &config)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, config, http.StatusOK)
}

func (api API) handleDeletePipelineConfigRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	err := pipeline.DeletePipelineConfig(id)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, map[string]interface{}{"ok": true}, http.StatusOK)
}
