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
	"github.com/julienschmidt/httprouter"
	. "github.com/medcl/gopa/core/api"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/modules/config"
	"net/http"
	"strconv"
)

func (this API) TaskDeleteAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method == DELETE.String() {
		id := ps.ByName("id")
		err := model.DeleteTask(id)
		if err != nil {
			this.Error(w, err)
		} else {
			this.WriteJson(w, map[string]interface{}{"ok": true}, http.StatusOK)
		}
	} else {
		this.Error404(w)
	}
}
func (this API) TaskGetAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	task, err := model.GetTask(id)
	if err != nil {
		this.Error(w, err)
	} else {
		this.WriteJson(w, task, http.StatusOK)

	}

}

func (this API) TaskAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	if req.Method == POST.String() {
		jsonq, err := this.GetJson(req)
		if err != nil {
			logger.Error(err)
		}

		seed, err := jsonq.String("seed")
		if err != nil {
			logger.Error(err)
		}
		logger.Trace("receive new seed:", seed)

		task := model.NewTaskSeed(seed, "", 0,0)

		queue.Push(config.CheckChannel, task.MustGetBytes())

		this.WriteJson(w, map[string]interface{}{"ok": true}, http.StatusOK)
	} else {
		logger.Trace("get all tasks")

		fr := this.GetParameter(req, "from")
		si := this.GetParameter(req, "size")
		domain := this.GetParameter(req, "domain")

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
			this.Error(w, err)
		} else {
			this.WriteListResultJson(w, total, tasks, http.StatusOK)
		}
	}
}

func (this API) DomainDeleteAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method == DELETE.String() {
		id := ps.ByName("id")
		err := model.DeleteTask(id)
		if err != nil {
			this.Error(w, err)
		} else {
			this.WriteJson(w, map[string]interface{}{"ok": true}, http.StatusOK)
		}
	} else {
		this.Error404(w)
	}
}

func (this API) DomainGetAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	task, err := model.GetTask(id)
	if err != nil {
		this.Error(w, err)
	} else {
		this.WriteJson(w, task, http.StatusOK)

	}

}

func (this API) DomainAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	if req.Method == POST.String() {
		//jsonq, err := this.GetJson(req)
		//if err != nil {
		//	logger.Error(err)
		//}
		//
		//seed, err := jsonq.String("seed")
		//if err != nil {
		//	logger.Error(err)
		//}
		//logger.Trace("receive new seed:", seed)
		//
		//task := model.NewTaskSeed(seed, "", 0)
		//
		//queue.Push(config.CheckChannel, task.MustGetBytes())

		this.WriteJson(w, map[string]interface{}{"ok": true}, http.StatusOK)
	} else {
		logger.Trace("get all domain settings")

		fr := this.GetParameter(req, "from")
		si := this.GetParameter(req, "size")
		domain := this.GetParameter(req, "domain")

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

			total := stats.Stat("domain.stats", v.Host+"."+stats.STATS_FETCH_TOTAL_COUNT)
			v.LinksCount = total
			newDomains = append(newDomains, v)
		}

		if err != nil {
			this.Error(w, err)
		} else {
			this.WriteListResultJson(w, total, newDomains, http.StatusOK)
		}
	}
}
