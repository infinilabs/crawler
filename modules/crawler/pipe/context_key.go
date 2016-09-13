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
type ContextKey string


const (
	CONTEXT_URL       ContextKey = "URL"
	CONTEXT_PAGE_ITEM ContextKey = "PAGE"
	CONTEXT_HOST ContextKey = "HOST"
	CONTEXT_URL_PATH ContextKey = "URL_PATH"
	CONTEXT_PAGE_METADATA ContextKey = "PAGE_METADATA"
	CONTEXT_SAVE_PATH ContextKey = "CONTEXT_SAVE_PATH"
	CONTEXT_SAVE_FILENAME ContextKey = "CONTEXT_SAVE_FILENAME"
)

func (this ContextKey) String() string {
	return string(this)
}
