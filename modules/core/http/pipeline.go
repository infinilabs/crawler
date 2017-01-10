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
	"net/http"
	"github.com/julienschmidt/httprouter"
	. "github.com/medcl/gopa/core/pipeline"
	log	"github.com/cihub/seelog"

	"encoding/json"
)



func   (this API)handleGetPipelineJointsRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Debug("get joints")
	joints:=GetAllRegisteredJoints()
	this.WriteJson(w,joints,http.StatusOK)

	return
}

func   (this API)handlePostPipelineJointsRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Debug("post joints")
	config := PipelineConfig{}
	//config.Name = "test_pipe_line"
	//
	//config.Context = &Context{}
	//config.Context.Init()
	//config.Context.Data["url"] = "gogol.com"
	//config.Context.Data["webpage"] = "hello world gogo "
	//
	//config.InputJoint = &JointConfig{JointName: "init_task", Parameters: map[string]interface{}{"TASK_ID": "b1pr9kqaukih1tdt5ncg"}}
	//joints := []*JointConfig{}
	//joints = append(joints, &JointConfig{JointName: "url_normalization", Parameters: map[string]interface{}{}})
	//joints = append(joints, &JointConfig{JointName: "fetch", Parameters: map[string]interface{}{}})
	//joints = append(joints, &JointConfig{JointName: "url_ext_filter", Parameters: map[string]interface{}{}})
	////joints = append(joints, &JointConfig{JointName: "save2db", Parameters: map[string]interface{}{}})
	//
	//config.ProcessJoints = joints

	data,err:=this.GetRawBody(req)
	if(err!=nil){
		this.Error(w,err)
		return
	}

	err=json.Unmarshal(data,&config)
	if(err!=nil){
		this.Error(w,err)
		return
	}
	//TODO save for later use
	pipe := NewPipelineFromConfig(&config)
	pipe.Run()

	this.WriteJson(w,config,http.StatusOK)

	return
}

func   (this API)handleGetPipelinesRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	return
}


func   (this API)handlePostPipelinesRequest(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	return
}

