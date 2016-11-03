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

package env

import (
	"github.com/medcl/gopa/core/types"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/queue"

	"time"
	"github.com/medcl/gopa/core/stats"
)

type Channels struct {
	pendingFetchDiskQueue BackendQueue
	pendingCheckDiskQueue BackendQueue
}

func (this *Channels) Init(path string)  {
	this.pendingFetchDiskQueue = NewDiskQueue("pending_fetch", path, 100*1024*1024, 4, 1<<10, 2500, 10*time.Second)
	this.pendingCheckDiskQueue = NewDiskQueue("pending_check", path, 100*1024*1024, 4, 1<<10, 2500, 10*time.Second)
}

func (this *Channels) PushUrlToCheck(url types.PageTask) error {
	err := this.pendingCheckDiskQueue.Put(url.MustGetBytes())
	stats.Increment("global", stats.STATS_CHECKER_PUSH_DISK_COUNT)
	return err
}

func (this *Channels) PushUrlToFetch(url types.PageTask) error {
	err := this.pendingFetchDiskQueue.Put(url.MustGetBytes())
	stats.Increment("global", stats.STATS_FETCH_PUSH_DISK_COUNT)
	return err
}

func (this *Channels) PopUrlToCheck() (types.PageTask, error) {
	b := <-this.pendingCheckDiskQueue.ReadChan()
	url := types.PageTaskFromBytes(b)
	stats.Increment("global", stats.STATS_CHECKER_POP_DISK_COUNT)
	return url, nil
}

func (this *Channels) PopUrlToFetch() (types.PageTask, error) {
	log.Trace("start pop url from queue")
	b := <-this.pendingFetchDiskQueue.ReadChan()
	url := types.PageTaskFromBytes(b)
	log.Trace("end pop url from queue")
	stats.Increment("global", stats.STATS_FETCH_POP_DISK_COUNT)
	return url, nil
}

func (this *Channels) Close() {
	this.pendingFetchDiskQueue.Close()
	this.pendingCheckDiskQueue.Close()
}
