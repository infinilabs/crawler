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
	"github.com/infinitbyte/gopa/core/http"
)

// API namespace
type API struct {
	api.Handler
}

// InitAPI register apis
func InitAPI() {

	apis := API{}

	user := "gopa"
	pass := "gopa"

	//Index
	api.HandleAPIMethod(api.GET, "/", apis.IndexAction)
	api.HandleAPIMethod(api.GET, "/favicon.ico", apis.IndexAction)

	//Stats APIs
	api.HandleAPIFunc("/stats", apis.StatsAction)
	api.HandleAPIFunc("/queue/stats", apis.QueueStatsAction)

	//Task API
	api.HandleAPIMethod(api.GET, "/tasks/", apis.TaskAction)
	api.HandleAPIMethod(api.GET, "/task/", apis.TaskAction)
	api.HandleAPIMethod(api.POST, "/task/", apis.CreateTaskAction)
	api.HandleAPIMethod(api.GET, "/task/:id", apis.TaskGetAction)
	api.HandleAPIMethod(api.DELETE, "/task/:id", api.BasicAuth(apis.TaskDeleteAction, user, pass))

	//Host API
	api.HandleAPIMethod(api.GET, "/host/", apis.GetHostsAction)
	api.HandleAPIMethod(api.GET, "/host/:id", apis.GetHostAction)
	api.HandleAPIMethod(api.DELETE, "/host/:id", api.BasicAuth(apis.DeleteHostAction, user, pass))

	//Snapshot API
	api.HandleAPIMethod(api.GET, "/snapshots/", apis.SnapshotListAction)
	api.HandleAPIMethod(api.GET, "/snapshot/", apis.SnapshotListAction)
	api.HandleAPIMethod(api.GET, "/snapshot/:id", apis.SnapshotGetAction)
	api.HandleAPIMethod(api.GET, "/snapshot/:id/payload", apis.SnapshotGetPayloadAction)

	//Pipeline API
	api.HandleAPIMethod(api.GET, "/pipeline/joints/", apis.handleGetPipelineJointsRequest)
	api.HandleAPIMethod(api.GET, "/pipeline/joint/", apis.handleGetPipelineJointsRequest)

	//api.HandleAPIMethod(api.GET, "/pipeline/configs/", apis.handleGetPipelineConfigsRequest)
	//api.HandleAPIMethod(api.GET, "/pipeline/config/", apis.handleGetPipelineConfigsRequest)
	api.HandleAPIMethod(api.POST, "/pipeline/config/", apis.handleCreatePipelineConfigRequest)

	api.HandleAPIMethod(api.GET, "/pipeline/config/:id", apis.handleGetPipelineConfigRequest)
	api.HandleAPIMethod(api.POST, "/pipeline/config/:id", apis.handleUpdatePipelineConfigRequest)
	api.HandleAPIMethod(api.DELETE, "/pipeline/config/:id", apis.handleDeletePipelineConfigRequest)

	//User API
	api.HandleAPIMethod(api.GET, "/user/", apis.handleUserLoginRequest)

}
