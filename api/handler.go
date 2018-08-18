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

	handler := API{}

	//Index
	api.HandleAPIMethod(api.GET, "/", handler.IndexAction)
	api.HandleAPIMethod(api.GET, "/favicon.ico", handler.IndexAction)

	//Stats APIs
	api.HandleAPIFunc("/stats", handler.StatsAction)
	api.HandleAPIFunc("/queue/stats", handler.QueueStatsAction)

	//Project API
	api.HandleAPIMethod(api.GET, "/projects/", handler.GetProjectsAction)
	api.HandleAPIMethod(api.POST, "/project/", handler.CreateProjectAction)
	api.HandleAPIMethod(api.GET, "/project/:id", handler.GetProjectAction)
	api.HandleAPIMethod(api.PUT, "/project/:id", handler.UpdateProjectAction)
	api.HandleAPIMethod(api.DELETE, "/project/:id", handler.DeleteProjectAction)

	//Host API
	api.HandleAPIMethod(api.GET, "/hosts/", handler.GetHostsAction)
	api.HandleAPIMethod(api.GET, "/host/:id", handler.GetHostAction)
	api.HandleAPIMethod(api.DELETE, "/host/:id", handler.DeleteHostAction)

	//Host Config API
	api.HandleAPIMethod(api.GET, "/host_configs/", handler.GetHostConfigsAction)
	api.HandleAPIMethod(api.POST, "/host_config/", handler.CreateHostConfigAction)
	api.HandleAPIMethod(api.GET, "/host_config/:id", handler.GetHostConfigAction)
	api.HandleAPIMethod(api.PUT, "/host_config/:id", handler.UpdateHostConfigAction)
	api.HandleAPIMethod(api.DELETE, "/host_config/:id", handler.DeleteHostConfigAction)

	//Task API
	api.HandleAPIMethod(api.GET, "/tasks/", handler.TaskAction)
	api.HandleAPIMethod(api.POST, "/task/", handler.CreateTaskAction)
	api.HandleAPIMethod(api.GET, "/task/:id", handler.TaskGetAction)
	api.HandleAPIMethod(api.PUT, "/task/:id", handler.TaskUpdateAction)
	api.HandleAPIMethod(api.DELETE, "/task/:id", handler.TaskDeleteAction)

	//Snapshot API
	api.HandleAPIMethod(api.GET, "/snapshots/", handler.SnapshotListAction)
	api.HandleAPIMethod(api.GET, "/snapshot/:id", handler.SnapshotGetAction)
	api.HandleAPIMethod(api.GET, "/snapshot/:id/payload", handler.SnapshotGetPayloadAction)

	//User API
	api.HandleAPIMethod(api.GET, "/user/", handler.handleUserLoginRequest)

}
