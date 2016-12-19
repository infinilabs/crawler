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

package modules

import (
	//"github.com/medcl/gopa/modules/checker"
	//apiModule "github.com/medcl/gopa/modules/api"
	//crawlerModule "github.com/medcl/gopa/modules/crawler"
	//storageModule "github.com/medcl/gopa/modules/storage"
	//. "github.com/medcl/gopa/core/env"
)
import (
	"github.com/medcl/gopa/core/module"
	"github.com/medcl/gopa/modules/api"
	"github.com/medcl/gopa/modules/checker"
	"github.com/medcl/gopa/modules/parser"
	"github.com/medcl/gopa/modules/crawler"
	"github.com/medcl/gopa/modules/storage"
)

func Register() {
	//register modules
	module.Register(http.APIModule{})
	module.Register(storage.StorageModule{})
	module.Register(crawler.CrawlerModule{})
	module.Register(parser.ParserModule{})
	module.Register(url_checker.CheckerModule{})
}
