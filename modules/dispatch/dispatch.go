package dispatch

import (
	log "github.com/cihub/seelog"
	cfg "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/util"
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
		offset := defaultOffset

		for {
			select {
			case <-signalChannel:
				log.Trace("dispatcher exited")
				return
			case data := <-queue.ReadChan(config.DispatcherChannel):

				//slow down while too many task already in the queue
				depth := queue.Depth(config.FetchChannel)
				if depth > 100 { //TODO configable
					log.Debugf("too many tasks already in the queue, depth: %v, wait 10s", depth)
					time.Sleep(10 * time.Second)
					return
				}

				stats.Increment("queue."+string(config.DispatcherChannel), "pop")
				log.Trace("got dispatcher signal, ", string(data))

				//get new task
				total, tasks, err := model.GetPendingNewFetchTasks()
				if err != nil {
					log.Error(err)
				}
				log.Debugf("get %v new task", total)

				isUpdate := false
				//get update task
				if tasks == nil || total <= 0 {
					total, tasks, err = model.GetPendingUpdateFetchTasks(offset)
					log.Debugf("get %v update task", total)
					isUpdate = true
					if total == 0 {
						log.Trace("reset offset, ", defaultOffset)
						offset = defaultOffset
					}
				} else {
					log.Debugf("get %v new task", total)
				}

				if tasks != nil && total > 0 {
					for _, v := range tasks {
						log.Trace("get task from db, ", v.ID)

						if err != nil {
							log.Error(err)
							panic(err)
						}

						//update offset
						if v.Created.After(offset) && isUpdate {
							offset = v.Created
						}

						runner := "fetch"
						if v.PipelineConfigID == "" {
							//assign pipeline config
							configID := model.GetPipelineIDByUrl(runner, v.Host, v.Url)
							if configID != "" {
								v.PipelineConfigID = configID
							}
						}

						context := model.Context{}
						context.Init()
						context.Set(model.CONTEXT_TASK_ID, v.ID)
						context.PipelineConfigID = v.PipelineConfigID

						queue.Push(config.FetchChannel, util.ToJSONBytes(context))
						if isUpdate {
							stats.Increment("dispatch", "update.enqueue")
						}
					}
				}

				//minimum  wait time
				time.Sleep(5 * time.Second)
			}

		}
	}()

	go func() {
		for {
			var lastPop int64
			select {
			case <-signalChannel:
				log.Trace("auto dispatcher exited")
				return
			default:
				pop := stats.Stat("queue.fetch", "pop")
				if lastPop == pop {
					lastPop = pop
					log.Trace("fetch tasks stalled, try to dispatch some tasks from db")
					err := queue.Push(config.DispatcherChannel, []byte("10s auto"))
					if err != nil {
						log.Error(err)
					}
					time.Sleep(10 * time.Second)
				} else {
					time.Sleep(20 * time.Second)
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
