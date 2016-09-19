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

import . "github.com/medcl/gopa/core/pipeline"

type UrlSource struct {
	Url string
	Reference string
	Depth int
}

func (this UrlSource) Name() string {
	return "url_source"
}

func (this UrlSource) Process(context *Context) (*Context, error) {

	context.Set(CONTEXT_ORIGINAL_URL,this.Url)
	context.Set(CONTEXT_URL,this.Url)
	context.Set(CONTEXT_DEPTH,this.Depth)
	context.Set(CONTEXT_REFERENCE_URL,this.Reference)
	return context, nil
}

