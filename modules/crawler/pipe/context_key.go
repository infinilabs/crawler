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

import . "github.com/infinitbyte/gopa/core/pipeline"

const (
	CONTEXT_TASK_ID          ParaKey = "CRAWLER_TASK_ID"
	CONTEXT_CRAWLER_DOMAIN   ParaKey = "CRAWLER_DOMAIN"
	CONTEXT_CRAWLER_TASK     ParaKey = "CRAWLER_TASK"
	CONTEXT_CRAWLER_SNAPSHOT ParaKey = "CRAWLER_SNAPSHOT"

	CONTEXT_PAGE_LINKS ParaKey = "PAGE_LINKS"
)
