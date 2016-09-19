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

import (
	. "github.com/medcl/gopa/core/pipeline"
	"strings"
)

type UrlFilterJoint struct {

}


func (this UrlFilterJoint) Name() string {
	return "url_filgter"
}

func (this UrlFilterJoint) Process(context *Context) (*Context, error) {

	url := context.MustGetString(CONTEXT_URL)
	orgUrl := context.MustGetString(CONTEXT_ORIGINAL_URL)

	if(orgUrl==""){
		orgUrl=url
	}

	if((!this.valid(orgUrl))||(url!=orgUrl&&(!this.valid(url)))){
		context.Break()
	}

	return context, nil
}

func (this UrlFilterJoint)valid(url string) bool {
	if(strings.HasPrefix(url,"mailto:")){
		return false
	}

	return true
}
