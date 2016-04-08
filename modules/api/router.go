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
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/medcl/gopa/core/config"
	. "github.com/medcl/gopa/modules/api/handler"
)

func getRouter(config *config.GopaConfig) (rest.App, error) {
	handler := Handler{Config: config}

	return rest.MakeRouter(
		rest.Get("/", handler.IndexAction),
		rest.Get("/stats", handler.StatsAction),
	)
}
