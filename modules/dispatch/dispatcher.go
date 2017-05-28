package dispatch

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/filter"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/modules/config"
	"time"
)

type DispatchModule struct {
}

func (this DispatchModule) Name() string {
	return "Dispatcher"
}

var signalChannel chan bool
var quitChannel chan bool

func (this DispatchModule) Start(env *Env) {
	signalChannel = make(chan bool, 2)
	quitChannel = make(chan bool, 3)
	go func() {
		for {
			select {
			case <-signalChannel:
				quitChannel <- true
				log.Trace("dispatcher exited")
				return
			default:
				log.Trace("waiting dispatcher signal")
				er, v := queue.Pop(config.DispatcherChannel)
				log.Trace("got dispatcher signal, ", string(v))
				if er != nil {
					log.Trace(er)
					continue
				}

				_, tasks, err := model.GetPendingFetchTasks()
				if err != nil {
					log.Trace(err)
					continue
				}

				if tasks != nil {
					for _, v := range tasks {
						log.Debug("get task from db, ", v.ID)

						b, err := filter.CheckThenAdd(config.FetchFilter, []byte(v.ID))

						if err != nil {
							log.Error(err)
							panic(err)
						}

						if b {
							log.Debug("url seems already fetched, ignore now, ", v.ID)
							continue
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
				quitChannel <- true
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
