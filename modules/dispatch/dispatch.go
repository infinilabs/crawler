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

type DispatchConfig struct {
	FailedTaskEnabled bool `config:"failed_retry"`
	NewTaskEnabled    bool `config:"new_task"`
	UpdateTaskEnabled bool `config:"update_task"`
}

// Name return Dispatch
func (module DispatchModule) Name() string {
	return "Dispatch"
}

var signalChannel chan bool

// Start dispatch module
func (module DispatchModule) Start(cfg *cfg.Config) {
	moduleConfig := DispatchConfig{FailedTaskEnabled: false, UpdateTaskEnabled: true, NewTaskEnabled: true}
	cfg.Unpack(&moduleConfig)

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

				var total int
				var tasks []model.Task
				var err error
				//get new task
				if moduleConfig.NewTaskEnabled {
					total, tasks, err = model.GetPendingNewFetchTasks()
					if err != nil {
						log.Error(err)
					}
					log.Debugf("get %v new task", total)
				}

				isUpdate := false
				//get update task
				if moduleConfig.UpdateTaskEnabled && (tasks == nil || total <= 0) {
					total, tasks, err = model.GetPendingUpdateFetchTasks(offset)
					log.Debugf("get %v update task, with offset, %v", total, offset)
					isUpdate = true
					if total == 0 {
						log.Tracef("%v hit 0 update task, reset offset to %v ", offset, defaultOffset)
						offset = defaultOffset
					}
				}

				// get failed task
				if moduleConfig.FailedTaskEnabled && (tasks == nil || total <= 0) {
					total, tasks, err = model.GetFailedTasks(offset)
					log.Debugf("get %v failed task, with offset, %v", total, offset)
					if total == 0 {
						log.Tracef("%v hit 0 failed task, reset offset to %v ", offset, defaultOffset)
						offset = defaultOffset
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
				push := stats.Stat("queue.fetch", "push")
				if lastPop == pop || pop == push {
					lastPop = pop
					log.Trace("fetch tasks stalled/finished, try to dispatch some tasks from db")
					err := queue.Push(config.DispatcherChannel, []byte("5s auto"))
					if err != nil {
						log.Error(err)
					}
					time.Sleep(5 * time.Second)
				} else {
					log.Trace("no new data in fetch queue, wait 20s")
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
