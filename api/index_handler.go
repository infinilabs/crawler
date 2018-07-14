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
	"github.com/infinitbyte/framework/core/api/router"
	"github.com/infinitbyte/framework/core/env"
	"github.com/infinitbyte/framework/core/global"
	"net/http"
	"time"
)

// IndexAction returns cluster health information
func (handler API) IndexAction(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	if req.URL.Path != "/" {
		handler.Error404(w)
		return
	}

	data := map[string]interface{}{}
	data["cluster_name"] = global.Env().SystemConfig.ClusterConfig.Name
	data["name"] = global.Env().SystemConfig.NodeConfig.Name

	version := map[string]interface{}{}
	version["number"] = global.Env().GetVersion()
	version["build_commit"] = global.Env().GetLastCommitLog()
	version["build_date"] = global.Env().GetBuildDate()

	data["version"] = version
	data["tagline"] = "You Know, for Web"
	data["uptime"] = time.Since(env.GetStartTime()).String()

	handler.WriteJSON(w, &data, http.StatusOK)
}
