package dispatch

import (
	log "github.com/cihub/seelog"
	cfg "github.com/infinitbyte/framework/core/config"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/queue"
	"github.com/infinitbyte/framework/core/stats"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
	"time"
)

// DispatchModule handle task dispatch, include new task and update task
type DispatchModule struct {
}

type DispatchConfig struct {
	FailedTaskEnabled       bool  `config:"failure_retry"`
	NewTaskEnabled          bool  `config:"new_task"`
	UpdateTaskEnabled       bool  `config:"update_task"`
	MaxConcurrentFetchTasks int64 `config:"max_concurrent_fetch_tasks"`
	FetchSize               int   `config:"fetch_size"`
	GetHostConfigByUrl      bool  `config:"get_host_config_by_url"`
}

// Name return Dispatch
func (module DispatchModule) Name() string {
	return "Dispatch"
}

var signalChannel chan bool

func dispatchTasks(name string, cfg *DispatchConfig, tasks []model.Task, offset int64) {
	for _, v := range tasks {
		log.Trace("get task from db, ", v.ID)

		if v.Status == model.TaskPendingFetch {
			log.Tracef("task %v alrady in fetch, ingore", v.ID)
			continue
		}

		context := pipeline.Context{}
		context.Set(model.CONTEXT_TASK_ID, v.ID)

		//update offset
		if v.Created.Unix() > offset {
			offset = v.Created.Unix()
		}

		runner := "fetch"

		if v.HostConfig == nil && cfg.GetHostConfigByUrl {
			//assign pipeline config
			config, _ := model.GetHostConfigByHostAndUrl(runner, v.Host, v.Url)
			if config != nil {
				v.HostConfig = config
				context.Set(model.CONTEXT_TASK_Cookies, v.HostConfig.Cookies)
				context.Set(model.CONTEXT_TASK_PipelineConfigID, v.HostConfig.PipelineID)
				log.Trace("get host config: ", util.ToJson(config, true))
			}
		}

		if v.PipelineConfigID != "" {
			context.Set(model.CONTEXT_TASK_PipelineConfigID, v.PipelineConfigID)
		}

		//Update task status
		v.Status = model.TaskPendingFetch
		model.UpdateTask(&v)
		err := queue.Push(config.FetchChannel, util.ToJSONBytes(context))
		log.Trace("push task to fetch queue, ", v.ID)

		if err != nil {
			log.Error(err)
			continue
		}
	}

}

// Start dispatch module
func (module DispatchModule) Start(cfg *cfg.Config) {
	moduleConfig := DispatchConfig{
		FailedTaskEnabled:       false,
		UpdateTaskEnabled:       true,
		NewTaskEnabled:          true,
		MaxConcurrentFetchTasks: 100,
		FetchSize:               100,
		GetHostConfigByUrl:      false,
	}
	cfg.Unpack(&moduleConfig)

	signalChannel = make(chan bool, 2)

	go func() {
		now := time.Now().UTC()
		dd, _ := time.ParseDuration("-240h")
		defaultOffset := now.Add(dd).Unix()
		newOffset := defaultOffset
		updateOffset := defaultOffset
		failureOffset := defaultOffset

		for {
			select {
			case <-signalChannel:
				log.Trace("dispatcher exited")
				return
			case data := <-queue.ReadChan(config.DispatcherChannel):

				stats.Increment("queue."+string(config.DispatcherChannel), "pop")

				//slow down while too many task already in the queue
				depth := queue.Depth(config.FetchChannel)
				if depth > moduleConfig.MaxConcurrentFetchTasks {
					log.Tracef("too many tasks already in the queue, depth: %v, wait 1s", depth)
					time.Sleep(1 * time.Second)
					continue
				}

				log.Trace("got dispatcher signal, ", string(data))

				var total int
				var tasks []model.Task
				var err error
				//get new task
				if moduleConfig.NewTaskEnabled {
					total, tasks, err = model.GetPendingNewFetchTasks(newOffset, moduleConfig.FetchSize)
					if err != nil {
						log.Error(err)
					}
					log.Debugf("get %v new task, with offset, %v", total, newOffset)

					if tasks != nil && len(tasks) > 0 {
						dispatchTasks("new", &moduleConfig, tasks, newOffset)
						continue
					}

					if total == 0 {
						log.Tracef("%v hit 0 new task, reset offset to %v ", newOffset, defaultOffset)
						newOffset = defaultOffset
					}
				}

				//get update task
				if moduleConfig.UpdateTaskEnabled && (tasks == nil || total <= 0) {
					total, tasks, err = model.GetPendingUpdateFetchTasks(updateOffset)
					if err != nil {
						log.Error(err)
					}
					log.Debugf("get %v update task, with offset, %v", total, updateOffset)

					if tasks != nil && len(tasks) > 0 {
						dispatchTasks("update", &moduleConfig, tasks, updateOffset)
						continue
					}

					if total == 0 {
						log.Tracef("%v hit 0 update task, reset offset to %v ", updateOffset, defaultOffset)
						updateOffset = defaultOffset
					}
				}

				// get failure task
				if moduleConfig.FailedTaskEnabled && (tasks == nil || total <= 0) {
					total, tasks, err = model.GetFailedTasks(failureOffset)
					if err != nil {
						log.Error(err)
					}
					log.Debugf("get %v failure task, with offset, %v", total, failureOffset)

					if tasks != nil && len(tasks) > 0 {
						dispatchTasks("failure", &moduleConfig, tasks, failureOffset)
						continue
					}

					if total == 0 {
						log.Tracef("%v hit 0 failure task, reset offset to %v ", failureOffset, defaultOffset)
						failureOffset = defaultOffset
					}
				}
			}
		}
	}()

	go func() {
		var lastPop int64

		for {
			select {
			case <-signalChannel:
				log.Trace("auto dispatcher exited")
				return
			default:
				time.Sleep(1 * time.Second)

				if queue.Depth(config.DispatcherChannel) > 1 {
					continue
				}
				pop := stats.Stat("queue.fetch", "pop")
				push := stats.Stat("queue.fetch", "push")
				if lastPop == pop || pop == push {
					log.Tracef("fetch tasks stalled/finished, lastPop:%v,pop:%v,push:%v", lastPop, pop, push)
					err := queue.Push(config.DispatcherChannel, []byte("10s deplay"))
					if err != nil {
						log.Error(err)
						continue
					}
				}
				//update lastPop
				lastPop = pop
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
