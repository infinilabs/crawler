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

const (
	CONTEXT_CRAWLER_TASK         ParaKey = "CRAWLER_TASK"
	CONTEXT_ORIGINAL_URL         ParaKey = "ORIGINAL_URL"
	CONTEXT_URL                  ParaKey = "URL"
	CONTEXT_REFERENCE_URL        ParaKey = "REFERENCE_URL"
	CONTEXT_DEPTH                ParaKey = "DEPTH"
	CONTEXT_BREADTH              ParaKey = "BREADTH"
	CONTEXT_PAGE_BODY_BYTES      ParaKey = "PAGE_BODY_BYTES"
	CONTEXT_PAGE_BODY_PLAIN_TEXT ParaKey = "PAGE_BODY_PLAIN_TEXT" //extracted plain text from html
	CONTEXT_PAGE_ITEM            ParaKey = "PAGE"
	CONTEXT_HOST                 ParaKey = "HOST"
	CONTEXT_URL_PATH             ParaKey = "URL_PATH"
	CONTEXT_PAGE_METADATA        ParaKey = "PAGE_METADATA"
	CONTEXT_PAGE_LINKS           ParaKey = "PAGE_LINKS"
	CONTEXT_SAVE_PATH            ParaKey = "SAVE_PATH"
	CONTEXT_SAVE_FILENAME        ParaKey = "SAVE_FILENAME"
	CONTEXT_PAGE_LAST_FETCH      ParaKey = "PAGE_LAST_FETCH"

	CACHE_TTL_TIMEOUT_HOST ParaKey = "CACHE_TTL_TIMEOUT_HOST"
)
