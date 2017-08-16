package queue

import (
	"errors"
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/global"
	. "github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/modules/config"
	. "github.com/infinitbyte/gopa/modules/queue/disk_queue"
	"os"
	"time"
)

var queues map[QueueKey]*BackendQueue

type DiskQueue struct {
}

func (module DiskQueue) Name() string {
	return "Queue"
}

func (module DiskQueue) Start(cfg *Config) {
	queues = make(map[QueueKey]*BackendQueue)
	path := global.Env().SystemConfig.GetWorkingDir() + "/queue"
	os.Mkdir(path, 0777)
	pendingFetchDiskQueue := NewDiskQueue("pending_fetch", path, 100*1024*1024, 1, 1<<10, 2500, 5*time.Second,10)
	pendingCheckDiskQueue := NewDiskQueue("pending_check", path, 100*1024*1024, 1, 1<<10, 2500, 5*time.Second,10)
	pendingDispatchDiskQueue := NewDiskQueue("pending_dispatch", path, 100*1024*1024, 1, 1<<10, 2500, 5*time.Second,10)
	pendingIndexDiskQueue := NewDiskQueue("pending_index", path, 100*1024*1024, 1, 1<<20, 2500, 5*time.Second,10)
	queues[config.FetchChannel] = &pendingFetchDiskQueue
	queues[config.CheckChannel] = &pendingCheckDiskQueue
	queues[config.DispatcherChannel] = &pendingDispatchDiskQueue
	queues[config.IndexChannel] = &pendingIndexDiskQueue
	//TODO configable
	Register(module)
}

func (module DiskQueue) Push(k QueueKey, v []byte) error {
	return (*queues[k]).Put(v)
}

func (module DiskQueue) ReadChan(k QueueKey) chan []byte {

	return (*queues[k]).ReadChan()
}

func (module DiskQueue) Pop(k QueueKey, timeoutInSeconds time.Duration) (error, []byte) {

	if timeoutInSeconds > 0 {
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
	} else {
		b := <-(*queues[k]).ReadChan()
		return nil, b
	}
}

func (module DiskQueue) Close(k QueueKey) error {
	b := (*queues[k]).Close()
	return b
}

func (module DiskQueue) Depth(k QueueKey) int64 {
	b := (*queues[k]).Depth()
	return b
}

func (module DiskQueue) Stop() error {
	for _, v := range queues {
		err := (*v).Close()
		if err != nil {
			log.Debug(err)
		}
	}
	return nil
}
