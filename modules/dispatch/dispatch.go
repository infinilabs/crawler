package dispatch

import (
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/modules/config"
	"time"
)

type DispatchModule struct {
}

func (this DispatchModule) Name() string {
	return "Dispatch"
}

var signalChannel chan bool

func (this DispatchModule) Start(cfg *Config) {
	signalChannel = make(chan bool, 2)
	go func() {
		now := time.Now()
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
					log.Debug("get %v update task", total)
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
						if v.CreateTime.After(*offset) && isUpdate {
							offset = v.CreateTime
						}

						queue.Push(config.FetchChannel, []byte(v.ID))
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

func (this DispatchModule) Stop() error {
	signalChannel <- true
	signalChannel <- true
	return nil
}
