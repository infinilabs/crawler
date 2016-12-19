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

package pipe

import (. "github.com/medcl/gopa/core/pipeline"
. "github.com/medcl/gopa/core/types"
	"github.com/medcl/gopa/core/tasks"
)

type StartSeed struct {
	Seed TaskSeed
}

func (this StartSeed) Name() string {
	return "start_seed"
}

func (this StartSeed) Process(context *Context) (*Context, error) {

	//init task record
	task:=CrawlerTask{}
	task.Seed=&this.Seed
	tasks.CreateTask(task)

	context.Set(CONTEXT_ORIGINAL_URL,this.Seed.Url)
	context.Set(CONTEXT_URL,this.Seed.Url)
	context.Set(CONTEXT_DEPTH,this.Seed.Depth)
	context.Set(CONTEXT_REFERENCE_URL,this.Seed.Reference)

	return context, nil
}

