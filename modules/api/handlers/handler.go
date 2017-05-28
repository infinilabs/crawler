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
	"github.com/medcl/gopa/core/http"
)

type API struct {
	api.Handler
}

func InitAPI() {

	apis := API{}

	user := "gopa"
	pass := "gopa"

	//Index
	api.HandleAPIMethod(api.GET, "/", apis.IndexAction)
	api.HandleAPIMethod(api.GET, "/favicon.ico", apis.IndexAction)

	//APIs
	api.HandleAPIFunc("/stats", apis.StatsAction)

	api.HandleAPIMethod(api.POST, "/task/", apis.TaskAction)
	api.HandleAPIMethod(api.GET, "/task", apis.TaskAction)
	api.HandleAPIMethod(api.GET, "/task/:id", apis.TaskGetAction)
	api.HandleAPIMethod(api.DELETE, "/task/:id", api.BasicAuth(apis.TaskDeleteAction, user, pass))

	api.HandleAPIMethod(api.GET, "/domain", apis.DomainAction)
	api.HandleAPIMethod(api.GET, "/domain/:id", apis.DomainGetAction)
	api.HandleAPIMethod(api.DELETE, "/domain/:id", api.BasicAuth(apis.DomainDeleteAction, user, pass))

	//Snapshot
	api.HandleAPIFunc("/snapshot/", apis.SnapshotAction)

	api.HandleAPIMethod(api.GET, "/joint/", apis.handleGetPipelineJointsRequest)
	api.HandleAPIMethod(api.POST, "/joint/", apis.handlePostPipelineJointsRequest)
	api.HandleAPIMethod(api.POST, "/pipeline/:id", apis.handlePostPipelinesRequest)
	api.HandleAPIMethod(api.GET, "/pipeline/", apis.handleGetPipelinesRequest)
	api.HandleAPIMethod(api.GET, "/user/", apis.handleUserLoginRequest)

}
