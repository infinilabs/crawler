/*
Copyright Medcl (m AT medcl.net)

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

package model

type Domain struct {
	ID      string
	Host    string
	Port    string
	Favicon string
	Enabled bool
}

type Url struct {
	Domain Domain

	FirstPath  string
	SecondPath string
	ThirdPath  string

	FileExt string
}

type FetchTask struct {
	Url Url

	CurrentStatus int
	CurrentStage  int
	StageStatus   map[int]int
}

func (task *FetchTask) UpdateStatus(status int) {
	task.CurrentStage++
	task.CurrentStatus = status
	task.StageStatus[task.CurrentStage] = status
}

const StagePreFetch = 0
const StageFetch = 1
const StageAfterFetch = 2

const PreFetchPendingCheck = 3
const PreFetchCheck = 4
const PreFetchChecking = 5
const PreFetchCheckError = 6
