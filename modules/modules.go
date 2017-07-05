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
	"github.com/infinitbyte/gopa/core/module"
	"github.com/infinitbyte/gopa/modules/api"
	"github.com/infinitbyte/gopa/modules/cluster"
	"github.com/infinitbyte/gopa/modules/crawler"
	"github.com/infinitbyte/gopa/modules/database"
	"github.com/infinitbyte/gopa/modules/dispatch"
	"github.com/infinitbyte/gopa/modules/filter"
	"github.com/infinitbyte/gopa/modules/index"
	"github.com/infinitbyte/gopa/modules/queue"
	"github.com/infinitbyte/gopa/modules/stats"
	"github.com/infinitbyte/gopa/modules/storage"
	"github.com/infinitbyte/gopa/modules/ui"
)

func Register() {
	////register modules
	module.Register(database.DatabaseModule{})
	module.Register(filter.FilterModule{})
	module.Register(storage.StorageModule{})
	module.Register(stats.StatsStoreModule{})
	module.Register(stats.StatsDModule{})
	module.Register(queue.DiskQueue{})
	module.Register(crawler.CheckerModule{})
	module.Register(crawler.CrawlerModule{})
	module.Register(dispatch.DispatchModule{})
	module.Register(index.IndexModule{})
	module.Register(cluster.ClusterModule{})
	module.Register(ui.UIModule{})
	module.Register(api.APIModule{})
}
