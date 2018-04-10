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
	"github.com/infinitbyte/framework/core/api"
	"github.com/infinitbyte/framework/core/index"
	"github.com/infinitbyte/framework/core/ui"
	core "github.com/infinitbyte/framework/modules/ui/common"
	"github.com/infinitbyte/gopa/ui/search"
	"github.com/infinitbyte/gopa/ui/search/common"
	"github.com/infinitbyte/gopa/ui/websocket"
)

var (
	defaultConfig = common.IndexConfig{
		Elasticsearch: &index.ElasticsearchConfig{
			Endpoint:    "http://localhost:9200",
			IndexPrefix: "gopa-",
		},
		UIConfig: &common.UIConfig{
			Enabled:     true,
			SiteName:    "GOPA",
			SiteFavicon: "/static/assets/img/favicon.ico",
			SiteLogo:    "/static/assets/img/logo.svg",
		},
	}
)

// InitUI register ui handlers
func InitUI() {
	//Nav init
	core.RegisterNav("Tasks", "Tasks", "/admin/tasks/")

	//UI pages init
	admin := AdminUI{}

	ui.HandleUIMethod(api.GET, "/screenshot/:id", admin.GetScreenshotAction)

	ui.HandleUIMethod(api.GET, "/admin/tasks/", api.NeedPermission(api.PERMISSION_ADMIN_MINIMAL, admin.TasksPageAction))
	ui.HandleUIMethod(api.GET, "/admin/task/view/:id", api.NeedPermission(api.PERMISSION_ADMIN_MINIMAL, admin.TaskViewPageAction))

	indexConfig := defaultConfig

	//TODO resolve cfg
	//cfg.Unpack(&indexConfig)

	//register UI
	if indexConfig.UIConfig.Enabled {
		search := search.UserUI{}
		search.Config = indexConfig.UIConfig
		search.SearchClient = &index.ElasticsearchClient{Config: indexConfig.Elasticsearch}
		ui.HandleUIMethod(api.GET, "/", search.IndexPageAction)
		ui.HandleUIMethod(api.GET, "/m/", search.MobileIndexPageAction)
		ui.HandleUIMethod(api.GET, "/ajax_more_item/", search.AJAXMoreItemAction)
		ui.HandleUIMethod(api.GET, "/snapshot/:id", api.NeedPermission(api.PERMISSION_SNAPSHOT_VIEW, search.GetSnapshotPayloadAction))
		ui.HandleUIMethod(api.GET, "/suggest/", search.SuggestAction)
	}

	ui.HandleWebSocketCommand("SEED", "seed [url] eg: seed http://elastic.co", websocket.AddSeed)
	ui.HandleWebSocketCommand("LOG", "log [level]  eg: log debug", websocket.UpdateLogLevel)
	//ui.HandleWebsocketCommand("DIS", websocket.Dispatch)
	//ui.HandleWebsocketCommand("GET_TASK", websocket.GetTask)
}
