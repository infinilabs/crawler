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
	"github.com/infinitbyte/framework/core/api"
)

// API namespace
type API struct {
	api.Handler
}

// InitAPI register apis
func InitAPI() {

	apis := API{}

	//Index
	api.HandleAPIMethod(api.GET, "/", apis.IndexAction)
	api.HandleAPIMethod(api.GET, "/favicon.ico", apis.IndexAction)

	//Stats APIs
	api.HandleAPIFunc("/stats", apis.StatsAction)
	api.HandleAPIFunc("/queue/stats", apis.QueueStatsAction)

	//Project API
	api.HandleAPIMethod(api.GET, "/projects/", apis.GetProjectsAction)
	api.HandleAPIMethod(api.POST, "/project/", apis.CreateProjectAction)
	api.HandleAPIMethod(api.GET, "/project/:id", apis.GetProjectAction)
	api.HandleAPIMethod(api.PUT, "/project/:id", apis.UpdateProjectAction)
	api.HandleAPIMethod(api.DELETE, "/project/:id", apis.DeleteProjectAction)

	//Host API
	api.HandleAPIMethod(api.GET, "/hosts/", apis.GetHostsAction)
	api.HandleAPIMethod(api.GET, "/host/:id", apis.GetHostAction)
	api.HandleAPIMethod(api.DELETE, "/host/:id", apis.DeleteHostAction)

	//Host Config API
	api.HandleAPIMethod(api.GET, "/host_configs/", apis.GetHostConfigsAction)
	api.HandleAPIMethod(api.POST, "/host_config/", apis.CreateHostConfigAction)
	api.HandleAPIMethod(api.GET, "/host_config/:id", apis.GetHostConfigAction)
	api.HandleAPIMethod(api.PUT, "/host_config/:id", apis.UpdateHostConfigAction)
	api.HandleAPIMethod(api.DELETE, "/host_config/:id", apis.DeleteHostConfigAction)

	//Task API
	api.HandleAPIMethod(api.GET, "/tasks/", apis.TaskAction)
	api.HandleAPIMethod(api.POST, "/task/", apis.CreateTaskAction)
	api.HandleAPIMethod(api.GET, "/task/:id", apis.TaskGetAction)
	api.HandleAPIMethod(api.PUT, "/task/:id", apis.TaskUpdateAction)
	api.HandleAPIMethod(api.DELETE, "/task/:id", apis.TaskDeleteAction)

	//Snapshot API
	api.HandleAPIMethod(api.GET, "/snapshots/", apis.SnapshotListAction)
	api.HandleAPIMethod(api.GET, "/snapshot/:id", apis.SnapshotGetAction)
	api.HandleAPIMethod(api.GET, "/snapshot/:id/payload", apis.SnapshotGetPayloadAction)

	//Pipeline API
	api.HandleAPIMethod(api.GET, "/pipeline/joints/", apis.handleGetPipelineJointsRequest)

	api.HandleAPIMethod(api.GET, "/pipeline/configs/", apis.handleGetPipelineConfigsRequest)
	api.HandleAPIMethod(api.POST, "/pipeline/config/", apis.handleCreatePipelineConfigRequest)
	api.HandleAPIMethod(api.GET, "/pipeline/config/:id", apis.handleGetPipelineConfigRequest)
	api.HandleAPIMethod(api.PUT, "/pipeline/config/:id", apis.handleUpdatePipelineConfigRequest)
	api.HandleAPIMethod(api.DELETE, "/pipeline/config/:id", apis.handleDeletePipelineConfigRequest)

	//User API
	api.HandleAPIMethod(api.GET, "/user/", apis.handleUserLoginRequest)

}
