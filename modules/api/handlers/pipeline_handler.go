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
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/julienschmidt/httprouter"
	"net/http"

	"encoding/json"
)

func (this API) handleGetPipelineJointsRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Debug("get joints")
	joints := GetAllRegisteredJoints()
	this.WriteJSON(w, joints, http.StatusOK)

	return
}

func (this API) handlePostPipelineJointsRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Debug("post joints")
	config := PipelineConfig{}

	data, err := this.GetRawBody(req)
	if err != nil {
		this.Error(w, err)
		return
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		this.Error(w, err)
		return
	}
	//TODO save for later use
	context := &Context{}
	context.Init()
	pipe := NewPipelineFromConfig(&config, context)
	pipe.Run()

	this.WriteJSON(w, config, http.StatusOK)

	return
}

func (this API) handleGetPipelinesRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	return
}

func (this API) handlePostPipelinesRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	return
}
