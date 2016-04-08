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
	"github.com/medcl/gopa/core/stats"
	"net/http"
)

func getMapValue(mapData map[string]int, key string, defaultValue int32) int {
	data := mapData[key]
	return data
}

func (this *Handler) StatsAction(w http.ResponseWriter, req *http.Request) {

	statsMap := stats.StatsAll()
	fetch := stats.StatsCount{
		TotalCount:   getMapValue(statsMap, stats.STATS_FETCH_COUNT, 0),
		FailCount:    getMapValue(statsMap, stats.STATS_FETCH_FAIL_COUNT, 0),
		SuccessCount: getMapValue(statsMap, stats.STATS_FETCH_SUCCESS_COUNT, 0),
		Ignore:       getMapValue(statsMap, stats.STATS_FETCH_IGNORE_COUNT, 0),
		Timeout:      getMapValue(statsMap, stats.STATS_FETCH_TIMEOUT_COUNT, 0)}
	parse := stats.StatsCount{
		TotalCount:   getMapValue(statsMap, stats.STATS_PARSE_COUNT, 0),
		FailCount:    getMapValue(statsMap, stats.STATS_PARSE_FAIL_COUNT, 0),
		SuccessCount: getMapValue(statsMap, stats.STATS_PARSE_SUCCESS_COUNT, 0),
		Ignore:       getMapValue(statsMap, stats.STATS_PARSE_IGNORE_COUNT, 0)}

	stats := stats.TaskStatus{Fetch: fetch, Parse: parse}

	this.WriteJson(w, &stats)
}
