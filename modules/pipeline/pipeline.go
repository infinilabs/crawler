/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pipeline

import (
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	. "github.com/infinitbyte/gopa/modules/pipeline/config"
	. "github.com/infinitbyte/gopa/modules/pipeline/joint"
	"runtime"
	"sync"
	"time"
)

var frameworkStarted bool
var pipes map[string]*Pipe

type PipelineFrameworkModule struct {
}

type Pipe struct {
	Config         PipeConfig
	l              sync.Mutex
	signalChannels []*chan bool
}

func (pipe *Pipe) Start(config PipeConfig) {
	if !pipe.Config.Enabled {
		log.Debugf("pipeline: %s was disabled", pipe.Config.Name)
		return
	}

	pipe.l.Lock()
	defer pipe.l.Unlock()
	pipe.Config = config

	numGoRoutine := pipe.Config.MaxGoRoutine

	pipe.signalChannels = make([]*chan bool, numGoRoutine)
	//start fetcher
	for i := 0; i < numGoRoutine; i++ {
		log.Trace("start pipeline instance:", i)
		signalC := make(chan bool, 1)
		pipe.signalChannels[i] = &signalC
		go pipe.runPipeline(&signalC, i)

	}
	log.Infof("pipeline: %s was started with %v instances", pipe.Config.Name, numGoRoutine)
}

func (pipe *Pipe) Update(config PipeConfig) {
	pipe.Stop()
	pipe.Start(config)
}

func (pipe *Pipe) Stop() {
	if !pipe.Config.Enabled {
		log.Debugf("pipeline: %s was disabled", pipe.Config.Name)
		return
	}
	pipe.l.Lock()
	defer pipe.l.Unlock()

	for i, item := range pipe.signalChannels {
		if item != nil {
			*item <- true
		}
		log.Debug("send exit signal to fetch channel: ", i)
	}
}

func (pipe *Pipe) runPipeline(singal *chan bool, shard int) {

	var taskInfo []byte
	for {
		select {
		case <-*singal:
			log.Trace("pipeline exit, shard:", shard)
			return
		case taskInfo = <-queue.ReadChan(config.FetchChannel):
			stats.Increment("queue."+string(config.FetchChannel), "pop")

			taskId, pipelineConfigId := model.DecodePipelineTask(taskInfo)

			pipelineConfig := pipe.Config.DefaultConfig
			if pipelineConfigId != "" {
				var err error
				pipelineConfig, err = model.GetPipelineConfig(pipelineConfigId)
				if err != nil {
					panic(err)
				}
			}

			log.Trace("shard:", shard, ",task received:", taskId)
			pipe.execute(taskId, pipelineConfig)
			log.Trace("shard:", shard, ",task finished:", taskId)
		}
	}
}

func (pipe *Pipe) execute(taskId string, pipelineConfig *model.PipelineConfig) {
	var pipeline *model.Pipeline
	defer func() {
		if !global.Env().IsDebug {
			if r := recover(); r != nil {
				if e, ok := r.(runtime.Error); ok {
					log.Error("pipeline: ", pipeline.GetID(), ", taskId: ", taskId, ", ", util.GetRuntimeErrorMessage(e))
				}
				log.Error("error in pipeline,", util.ToJson(r, true), util.ToJson(pipeline.GetContext(), true))
			}
		}
	}()

	context := &model.Context{}
	context.Set(CONTEXT_TASK_ID, taskId)

	pipeline = model.NewPipelineFromConfig(pipe.Config.Name, pipelineConfig, context)
	pipeline.Run()

	if pipe.Config.ThresholdInMs > 0 {
		log.Debug("sleep ", pipe.Config.ThresholdInMs, "ms to control crawling speed")
		time.Sleep(time.Duration(pipe.Config.ThresholdInMs) * time.Millisecond)
		log.Debug("wake up now,continue crawing")
	}

	log.Trace("end pipeline")
}

func (module PipelineFrameworkModule) Name() string {
	return "Pipeline"
}

func (module PipelineFrameworkModule) Start(cfg *Config) {

	if frameworkStarted {
		log.Error("pipeline framework already started, please stop it first.")
		return
	}

	//init joints
	InitJoints()
	var config = struct {
		Pipes []PipeConfig `config:"pipes"`
	}{GetDefaultPipeConfig()}

	cfg.Unpack(&config)

	pipes = map[string]*Pipe{}
	for i, v := range config.Pipes {
		if v.DefaultConfig == nil {
			panic(errors.Errorf("default pipeline config can't be null, %v, %v", i, v))
		}
		p := &Pipe{Config: v}
		pipes[v.Name] = p
	}

	log.Debug("starting up pipeline framework")
	for _, v := range pipes {
		v.Start(v.Config)
	}

	frameworkStarted = true
}

func (module PipelineFrameworkModule) Stop() error {
	if frameworkStarted {
		frameworkStarted = false
		log.Debug("shutting down pipeline framework")
		for _, v := range pipes {
			v.Stop()
		}
	} else {
		log.Error("pipeline framework is not started")
	}

	return nil
}
