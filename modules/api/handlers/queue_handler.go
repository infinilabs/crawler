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
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	"net/http"
)

// QueueStatsAction return queue stats information
func (handler API) QueueStatsAction(w http.ResponseWriter, req *http.Request) {

	data := map[string]int64{}
	data["check"] = queue.Depth(config.CheckChannel)
	data["fetch"] = queue.Depth(config.FetchChannel)
	data["dispatch"] = queue.Depth(config.DispatcherChannel)
	data["index"] = queue.Depth(config.IndexChannel)
	handler.WriteJSON(w, util.MapStr{
		"depth": data,
	}, 200)
}
