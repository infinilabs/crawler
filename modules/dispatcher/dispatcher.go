package dispatcher

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/filter"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/modules/config"
	"github.com/medcl/gopa/core/stats"
	"time"
)

type DispatcherModule struct {
}

var started bool

func (this DispatcherModule) Name() string {
	return "Dispatcher"
}

func (this DispatcherModule) Start(env *Env) {

	go func() {
		started = true
		for {
			log.Trace("get task from db")
			if started {
				log.Trace("waiting dispatcher signal")
				v := queue.Pop(config.DispatcherChannel)
				log.Trace("got dispatcher signal, ", string(v))

				_, tasks, err := model.GetPendingFetchTasks()
				if err != nil {
					log.Debug(err)
					return
				}

				if tasks != nil {
					for _, v := range tasks {
						log.Debug("get task from db, ", v.ID)

						b,err:= filter.CheckThenAdd(config.FetchFilter, []byte(v.ID))

						if(err!=nil){
							log.Error(err)
							panic(err)
						}

						if b{
							log.Debug("url seems already fetched, ignore now, ",v.ID)
							continue
						}

						queue.Push(config.FetchChannel, []byte(v.ID))
					}
				}

			}
		}
	}()

	go func() {
		started = true
		for {
			if started {

				pop:=stats.Stat("queue.fetch","pop")
				push:=stats.Stat("queue.fetch","push")

				time.Sleep(10*time.Second)

				pop2:=stats.Stat("queue.fetch","pop")
				push2:=stats.Stat("queue.fetch","push")
				if(push==push2&&pop==pop2){
					//log.Debug("fetch tasks stalled after 5 seconds, try to dispatch some tasks from db")
					err:=queue.Push(config.DispatcherChannel,[]byte("10s auto"))
					if(err!=nil){
						panic(err)
					}
				}else{
					time.Sleep(30*time.Second)
					continue
				}
			}
		}
	}()
}

func (this DispatcherModule) Stop() error {

	started = false
	return nil
}
