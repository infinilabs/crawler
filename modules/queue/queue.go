package queue

import (
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/global"
	. "github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/modules/config"
	. "github.com/medcl/gopa/modules/queue/disk_queue"
	"time"
)

var queues map[QueueKey]*BackendQueue

type DiskQueue struct {
}

func (this DiskQueue) Name() string {
	return "Queue"
}

func (this DiskQueue) Start(env *Env) {
	queues = make(map[QueueKey]*BackendQueue)
	path := global.Env().RuntimeConfig.PathConfig.QueueData
	pendingFetchDiskQueue := NewDiskQueue("pending_fetch", path, 100*1024*1024, 4, 1<<10, 2500, 10*time.Second)
	pendingCheckDiskQueue := NewDiskQueue("pending_check", path, 100*1024*1024, 4, 1<<10, 2500, 10*time.Second)
	queues[config.FetchChannel] = &pendingFetchDiskQueue
	queues[config.CheckChannel] = &pendingCheckDiskQueue
	//TODO configable
	Register(this)
}

func (this DiskQueue) Push(k QueueKey, v []byte) error {
	return (*queues[k]).Put(v)
}

func (this DiskQueue) Pop(k QueueKey) []byte {
	b := <-(*queues[k]).ReadChan()
	return b
}

func (this DiskQueue) Stop() error {
	for _, v := range queues {
		err := (*v).Close()
		if err != nil {
			return err
		}
	}
	return nil
}
