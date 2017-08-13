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
	"github.com/infinitbyte/gopa/core/stats"
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

// TaskAction handle task creation and return task list which support parameter: `from`, `size` and `domain`, eg:
// curl -XPOST "http://localhost:8001/task/" -d '{
//"seed":"http://elasticsearch.cn"
//}'
//curl -XGET http://127.0.0.1:8001/task?from=100&size=10&domain=elasticsearch.cn
func (handler API) TaskAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	if req.Method == api.POST.String() {
		jsonq, err := handler.GetJSON(req)
		if err != nil {
			logger.Error(err)
		}

		seed, err := jsonq.String("seed")
		if err != nil {
			logger.Error(err)
		}
		logger.Trace("receive new seed:", seed)

		task := model.NewTaskSeed(seed, "", 0, 0)

		queue.Push(config.CheckChannel, task.MustGetBytes())

		handler.WriteJSON(w, map[string]interface{}{"ok": true}, http.StatusOK)
	} else {
		logger.Trace("get all tasks")

		fr := handler.GetParameter(req, "from")
		si := handler.GetParameter(req, "size")
		domain := handler.GetParameter(req, "domain")

		from, err := strconv.Atoi(fr)
		if err != nil {
			from = 0
		}
		size, err := strconv.Atoi(si)
		if err != nil {
			size = 10
		}

		total, tasks, err := model.GetTaskList(from, size, domain)
		if err != nil {
			handler.Error(w, err)
		} else {
			handler.WriteJSONListResult(w, total, tasks, http.StatusOK)
		}
	}
}

// DomainDeleteAction handle domain deletion, only support delete by id, eg:
//curl -XDELETE http://127.0.0.1:8001/domain/1
func (handler API) DomainDeleteAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

// DomainGetAction return domain by domain id, eg:
//curl -XGET http://127.0.0.1:8001/domain/1
func (handler API) DomainGetAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	task, err := model.GetDomain(id)
	if err != nil {
		handler.Error(w, err)
	} else {
		handler.WriteJSON(w, task, http.StatusOK)

	}

}

// DomainAction return domain list, support parameter: `from`, `size` and `domain`, eg:
//curl -XGET http://127.0.0.1:8001/domain?from=0&size=10
func (handler API) DomainAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	if req.Method == api.GET.String() {

		logger.Trace("get all domain settings")

		fr := handler.GetParameter(req, "from")
		si := handler.GetParameter(req, "size")
		domain := handler.GetParameter(req, "domain")

		from, err := strconv.Atoi(fr)
		if err != nil {
			from = 0
		}
		size, err := strconv.Atoi(si)
		if err != nil {
			size = 10
		}

		total, domains, err := model.GetDomainList(from, size, domain)

		newDomains := []model.Domain{}
		for _, v := range domains {

			total := stats.Stat("domain.stats", v.Host+"."+config.STATS_FETCH_TOTAL_COUNT)
			v.LinksCount = total
			newDomains = append(newDomains, v)
		}

		if err != nil {
			handler.Error(w, err)
		} else {
			handler.WriteJSONListResult(w, total, newDomains, http.StatusOK)
		}
	}
}
