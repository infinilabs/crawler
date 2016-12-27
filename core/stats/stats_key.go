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

package stats

const STATS_FETCH_TOTAL_COUNT = "fetch.total"
const STATS_FETCH_SUCCESS_COUNT = "fetch.success"
const STATS_FETCH_FAIL_COUNT = "fetch.fail"
const STATS_FETCH_TIMEOUT_COUNT = "fetch.timeout"
const STATS_FETCH_IGNORE_COUNT = "fetch.ignore"

const STATS_STORAGE_FILE_SIZE = "stats.sum.file.size"
const STATS_STORAGE_FILE_COUNT = "stats.sum.file.count"


type StatsCount struct {
	TotalCount   int `json:"total,omitempty"`
	SuccessCount int `json:"success,omitempty"`
	FailCount    int `json:"fail,omitempty"`
	Ignore       int `json:"ignore,omitempty"`
	Timeout      int `json:"timeout,omitempty"`
}

type TaskStatus struct {
	Stats map[string]map[string]int `json:"stats"`
}
