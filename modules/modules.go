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
	"github.com/medcl/gopa/core/module"
	"github.com/medcl/gopa/modules/api"
	"github.com/medcl/gopa/modules/cluster"
	"github.com/medcl/gopa/modules/crawler"
	"github.com/medcl/gopa/modules/dispatcher"
	"github.com/medcl/gopa/modules/filter"
	"github.com/medcl/gopa/modules/queue"
	"github.com/medcl/gopa/modules/stats"
	"github.com/medcl/gopa/modules/storage"
	"github.com/medcl/gopa/modules/core"
	"github.com/medcl/gopa/modules/ui"
)

func Register() {
	//register modules
	module.Register(core.CoreModule{})
	module.Register(filter.FilterModule{})
	module.Register(storage.StorageModule{})
	module.Register(stats.StatsStoreModule{})
	//module.Register(statsd.StatsDModule{})
	module.Register(queue.DiskQueue{})
	module.Register(crawler.CheckerModule{})
	module.Register(crawler.CrawlerModule{})
	module.Register(dispatcher.DispatcherModule{})
	module.Register(cluster.ClusterModule{})
	module.Register(ui.UIModule{})
	module.Register(http.APIModule{})
}
