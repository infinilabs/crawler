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
	log "github.com/cihub/seelog"
	logger "github.com/cihub/seelog"
	api "github.com/infinitbyte/framework/core/api"
	"github.com/infinitbyte/framework/core/api/router"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/queue"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
	"net/http"
	"strconv"
)

// TaskDeleteAction handle task delete by id, eg:
//curl -XDELETE http://127.0.0.1:8001/task/1
func (handler API) TaskDeleteAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

// TaskGetAction return task model by task_id, eg:
//curl -XGET http://127.0.0.1:8001/task/1
func (handler API) TaskGetAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	task, err := model.GetTask(id)
	if err != nil {
		handler.Error(w, err)
	} else {
		handler.WriteJSON(w, task, http.StatusOK)

	}

}

func (api API) TaskUpdateAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	task := model.Task{}
	id := ps.ByName("id")

	data, err := api.GetRawBody(req)
	if err != nil {
		api.Error(w, err)
		return
	}

	err = json.Unmarshal(data, &task)
	if err != nil {
		api.Error(w, err)
		return
	}

	task.ID = id
	err = model.UpdateTask(&task)
	if err != nil {
		api.Error(w, err)
		return
	}

	api.WriteJSON(w, task, http.StatusOK)
}

// TaskAction handle task creation and return task list which support parameter: `from`, `size` and `host`, eg:
//curl -XGET http://127.0.0.1:8001/task?from=100&size=10&host=elasticsearch.cn
func (handler API) TaskAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	logger.Trace("get all tasks")

	fr := handler.GetParameter(req, "from")
	si := handler.GetParameter(req, "size")
	host := handler.GetParameter(req, "host")
	status := handler.GetIntOrDefault(req, "status", -1)

	from, err := strconv.Atoi(fr)
	if err != nil {
		from = 0
	}
	size, err := strconv.Atoi(si)
	if err != nil {
		size = 10
	}

	total, tasks, err := model.GetTaskList(from, size, host, status)
	if err != nil {
		handler.Error(w, err)
	} else {
		handler.WriteJSONListResult(w, total, tasks, http.StatusOK)
	}
}

// curl -XPOST "http://localhost:8001/task/" -d '{
//"url":"http://elasticsearch.cn",
//"pipeline_id":"1231231212312"
//}'
func (handler API) CreateTaskAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	jsonq, err := handler.GetJSON(req)
	if err != nil {
		logger.Error(err)
	}

	url, err := jsonq.String("url")
	if err != nil {
		logger.Error(err)
	}
	pipelineID, err := jsonq.String("pipeline_id")
	if err == nil {
		logger.Error(err)
	}

	context := pipeline.Context{IgnoreBroken: true}
	context.Set(model.CONTEXT_TASK_URL, url)
	context.Set(model.CONTEXT_TASK_PipelineConfigID, pipelineID)

	err = queue.Push(config.CheckChannel, util.ToJSONBytes(context))
	if err != nil {
		log.Error(err)
		return
	}

	handler.WriteJSON(w, map[string]interface{}{"ok": true}, http.StatusOK)

}
