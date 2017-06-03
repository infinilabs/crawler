package queue

import (
	"errors"
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/global"
	. "github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/modules/config"
	. "github.com/medcl/gopa/modules/queue/disk_queue"
	."github.com/medcl/gopa/core/config"
	"os"
	"time"
)

var queues map[QueueKey]*BackendQueue

type DiskQueue struct {
}

func (this DiskQueue) Name() string {
	return "Queue"
}

func (this DiskQueue) Start(cfg *Config) {
	queues = make(map[QueueKey]*BackendQueue)
	path := global.Env().SystemConfig.GetDataDir() + "/queue"
	os.Mkdir(path, 0777)
	pendingFetchDiskQueue := NewDiskQueue("pending_fetch", path, 100*1024*1024, 1, 1<<10, 2500, 5*time.Second)
	pendingCheckDiskQueue := NewDiskQueue("pending_check", path, 100*1024*1024, 1, 1<<10, 2500, 5*time.Second)
	pendingDispatchDiskQueue := NewDiskQueue("pending_dispatch", path, 100*1024*1024, 1, 1<<10, 2500, 5*time.Second)
	pendingIndexDiskQueue := NewDiskQueue("pending_index", path, 100*1024*1024, 1, 1<<20, 2500, 5*time.Second)
	queues[config.FetchChannel] = &pendingFetchDiskQueue
	queues[config.CheckChannel] = &pendingCheckDiskQueue
	queues[config.DispatcherChannel] = &pendingDispatchDiskQueue
	queues[config.IndexChannel] = &pendingIndexDiskQueue
	//TODO configable
	Register(this)
}

func (this DiskQueue) Push(k QueueKey, v []byte) error {
	return (*queues[k]).Put(v)
}

func (this DiskQueue) Pop(k QueueKey, timeoutInSeconds time.Duration) (error, []byte) {

	if timeoutInSeconds < 1 {
		timeoutInSeconds = 5
	}

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(timeoutInSeconds) // sleep 3 second
		timeout <- true
	}()
	select {
	case b := <-(*queues[k]).ReadChan():
		return nil, b
	case <-timeout:
		return errors.New("time out"), nil
	}
}

func (this DiskQueue) Close(k QueueKey) error {
	b := (*queues[k]).Close()
	return b
}

func (this DiskQueue) Stop() error {
	for _, v := range queues {
		err := (*v).Close()
		if err != nil {
			log.Debug(err)
		}
	}
	return nil
}
