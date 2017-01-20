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

package ui

import (
	"github.com/medcl/gopa/core/api"
	"github.com/medcl/gopa/static"
	"net/http"
)

type API struct {
	api.Handler
}

func InitAPI() {

	apis := API{}

	//UI pages
	api.HandleFunc("/ui/", apis.DashboardAction)
	api.HandleFunc("/ui/dashboard/", apis.DashboardAction)
	api.HandleFunc("/ui/tasks/", apis.TasksPageAction)
	api.HandleFunc("/ui/console/", apis.ConsolePageAction)
	api.HandleFunc("/ui/explore/", apis.ExplorePageAction)
	api.HandleFunc("/ui/boltdb/", apis.BoltDBStatusAction)
	api.Handle("/static/", http.FileServer(static.FS(false)))

}
