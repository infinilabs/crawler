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
	"github.com/medcl/gopa/core/stats"
	"net/http"
)

func getMapValue(mapData map[string]int, key string, defaultValue int32) int {
	data := mapData[key]
	return data
}

func (this API) StatsAction(w http.ResponseWriter, req *http.Request) {

	m := stats.StatsAll()
	this.WriteJsonHeader(w)
	this.Write(w, *m)
}
