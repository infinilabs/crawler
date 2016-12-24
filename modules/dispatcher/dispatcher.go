package dispatcher

import (
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/tasks"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/modules/config"
	log "github.com/cihub/seelog"
)

type DispatcherModule  struct {

}

var started bool

func (this DispatcherModule) Name() string {
	return "Dispatcher"
}

func (this DispatcherModule)Start(env *Env) {

	go func() {
		started=true
		for {
			log.Trace("get task from db")
			if started {
				log.Trace("waiting dispatcher signal")
				v:=queue.Pop(config.DispatcherChannel)
				log.Trace("got dispatcher signal, ",string(v))

				_,tasks,err:=tasks.GetPendingFetchTasks()
				if(err!=nil){
					log.Error(err)
				}

				if(tasks!=nil){
					for _,v:=range tasks{
						log.Debug("get task from db, ",v.ID)
						queue.Push(config.FetchChannel,[]byte(v.ID))
					}
				}



			}
		}
	}()

}

func (this DispatcherModule) Stop() error {

	started=false
	return nil
}

