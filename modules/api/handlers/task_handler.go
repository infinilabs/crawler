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
	logger "github.com/cihub/seelog"
	api "github.com/infinitbyte/gopa/core/http"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	"github.com/julienschmidt/httprouter"
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

// TaskAction handle task creation and return task list which support parameter: `from`, `size` and `host`, eg:

//curl -XGET http://127.0.0.1:8001/task?from=100&size=10&host=elasticsearch.cn
func (handler API) TaskAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	logger.Trace("get all tasks")

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

	total, tasks, err := model.GetTaskList(from, size, host)
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

	context := model.Context{IgnoreBroken: true}
	context.Set(model.CONTEXT_TASK_URL, url)
	context.PipelineConfigID = pipelineID

	queue.Push(config.CheckChannel, util.ToJSONBytes(context))

	handler.WriteJSON(w, map[string]interface{}{"ok": true}, http.StatusOK)

}

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
