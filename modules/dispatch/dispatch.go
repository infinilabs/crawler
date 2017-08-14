package dispatch

import (
	log "github.com/cihub/seelog"
	cfg "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/modules/config"
	"time"
)

// DispatchModule handle task dispatch, include new task and update task
type DispatchModule struct {
}

// Name return Dispatch
func (module DispatchModule) Name() string {
	return "Dispatch"
}

var signalChannel chan bool

// Start dispatch module
func (module DispatchModule) Start(cfg *cfg.Config) {
	signalChannel = make(chan bool, 2)
	go func() {
		now := time.Now().UTC()
		dd, _ := time.ParseDuration("-240h")
		defaultOffset := now.Add(dd)
		offset := &defaultOffset

		for {
			select {
			case <-signalChannel:
				log.Trace("dispatcher exited")
				return
			case data := <-queue.ReadChan(config.DispatcherChannel):
				stats.Increment("queue."+string(config.DispatcherChannel), "pop")
				log.Trace("got dispatcher signal, ", string(data))

				//get new task
				total, tasks, err := model.GetPendingNewFetchTasks()
				if err != nil {
					log.Error(err)
				}

				isUpdate := false
				//get update task
				if tasks == nil || total <= 0 {
					total, tasks, err = model.GetPendingUpdateFetchTasks(offset)
					log.Debugf("get %v update task", total)
					isUpdate = true
					if total == 0 {
						log.Trace("reset offset, ", defaultOffset)
						offset = &defaultOffset
					}
				}

				if tasks != nil && total > 0 {
					for _, v := range tasks {
						log.Trace("get task from db, ", v.ID)

						if err != nil {
							log.Error(err)
							panic(err)
						}

						//update offset
						if v.Created.After(*offset) && isUpdate {
							offset = v.Created
						}

						queue.Push(config.FetchChannel, []byte(v.ID))
						if isUpdate {
							stats.Increment("dispatch", "update.enqueue")
						}
					}
				}
			}

		}
	}()

	go func() {
		for {
			select {
			case <-signalChannel:
				log.Trace("auto dispatcher exited")
				return
			default:
				pop := stats.Stat("queue.fetch", "pop")
				push := stats.Stat("queue.fetch", "push")

				time.Sleep(10 * time.Second)

				pop2 := stats.Stat("queue.fetch", "pop")
				push2 := stats.Stat("queue.fetch", "push")
				if push == push2 && pop == pop2 {
					log.Trace("fetch tasks stalled after 5 seconds, try to dispatch some tasks from db")
					err := queue.Push(config.DispatcherChannel, []byte("10s auto"))
					if err != nil {
						log.Error(err)
					}
				} else {
					time.Sleep(10 * time.Second)
					continue
				}
			}
		}

	}()
}

// Stop dispatch module
func (module DispatchModule) Stop() error {
	signalChannel <- true
	signalChannel <- true
	return nil
}
