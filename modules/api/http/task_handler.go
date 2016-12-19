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
	"github.com/medcl/gopa/core/tasks"
	"github.com/medcl/gopa/core/types"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"strconv"
)

func (this *Handler) TaskDeleteAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	 if(req.Method == DELETE.String()) {
		 id:=ps.ByName("id")
		 id1,_:=strconv.Atoi(id)
		 tasks.DeleteTask(id1)
		 this.WriteJson(w, map[string]interface{}{"ok": true}, http.StatusOK)
	 }else{
		 this.error404(w)
	 }
}
func (this *Handler) TaskGetAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		 id:=ps.ByName("id")
		 id1,_:=strconv.Atoi(id)
		 task,err:=tasks.GetTask(id1)
		if(err!=nil){
			this.error(w,err)
		}else
		{
			this.WriteJson(w, task, http.StatusOK)

		}

}

func (this *Handler) TaskAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

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

		task := types.NewPageTask(seed, "", 0)

		tasks.CreateTask(task)

		this.WriteJson(w, map[string]interface{}{"ok": true}, http.StatusOK)
	} else {
		logger.Trace("get all tasks")

		tasks:=tasks.GetTaskList()
		this.WriteJson(w,tasks,http.StatusOK)
	}
}
