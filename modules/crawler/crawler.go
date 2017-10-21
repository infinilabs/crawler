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

package crawler

import (
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/global"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	. "github.com/infinitbyte/gopa/modules/crawler/config"
	. "github.com/infinitbyte/gopa/modules/crawler/pipe"
	"runtime"
	"time"
)

var signalChannels []*chan bool

var crawlerStarted bool

// GetDefaultTaskConfig return a default TaskConfig
func GetDefaultTaskConfig() TaskConfig {
	config := PipelineConfig{}
	config.Name = "crawler"
	start := JointConfig{}
	start.Enabled = true
	start.JointName = "init_task"
	config.StartJoint = &start
	save := JointConfig{}
	save.Enabled = true
	save.JointName = "save_task"

	urlNormalization := JointConfig{}
	urlNormalization.Enabled = true
	urlNormalization.JointName = "url_normalization"
	urlNormalization.Parameters = util.MapStr{
		"follow_all_domain": false,
		"follow_sub_domain": true,
	}

	fetchJoint := JointConfig{}
	fetchJoint.Enabled = true
	fetchJoint.JointName = "fetch"

	parse := JointConfig{}
	parse.Enabled = true
	parse.JointName = "parse"

	html2text := JointConfig{}
	html2text.Enabled = true
	html2text.JointName = "html2text"

	hash := JointConfig{}
	hash.Enabled = true
	hash.JointName = "hash"

	updateCheckTime := JointConfig{}
	updateCheckTime.Enabled = true
	updateCheckTime.JointName = "update_check_time"

	contentDeduplication := JointConfig{}
	contentDeduplication.Enabled = true
	contentDeduplication.JointName = "content_deduplication"

	saveSnapshot := JointConfig{}
	saveSnapshot.Enabled = true
	saveSnapshot.JointName = "save_snapshot_db"

	task_deduplication := JointConfig{}
	task_deduplication.Enabled = true
	task_deduplication.JointName = "task_deduplication"

	config.EndJoint = &save
	config.ProcessJoints = []*JointConfig{
		&urlNormalization,
		&fetchJoint,
		&parse,
		&html2text,
		&hash,
		&updateCheckTime,
		&contentDeduplication,
		&saveSnapshot,
		&task_deduplication,
	}

	defaultCrawlerConfig := TaskConfig{
		MaxGoRoutine:          10,
		FetchThresholdInMs:    0,
		DefaultPipelineConfig: &config,
	}

	return defaultCrawlerConfig
}

func (module CrawlerModule) Name() string {
	return "Crawler"
}

func (module CrawlerModule) Start(cfg *Config) {

	config := GetDefaultTaskConfig()
	cfg.Unpack(&config)
	module.config = &config

	//TODO
	InitJoints()

	if crawlerStarted {
		log.Error("crawler already started, please stop it first.")
		return
	}

	numGoRoutine := config.MaxGoRoutine
	signalChannels = make([]*chan bool, numGoRoutine)
	if true {
		//start fetcher
		for i := 0; i < numGoRoutine; i++ {
			log.Trace("start crawler:", i)
			signalC := make(chan bool, 1)
			signalChannels[i] = &signalC
			go module.runPipeline(&signalC, i)

		}
	} else {
		log.Info("crawler currently is not enabled")
		return
	}

	crawlerStarted = true
}

func (module CrawlerModule) Stop() error {
	if crawlerStarted {
		crawlerStarted = false
		log.Debug("start shutting down crawler")
		for i, item := range signalChannels {
			if item != nil {
				*item <- true
			}
			log.Debug("send exit signal to fetch channel: ", i)
		}

	} else {
		log.Error("crawler is not crawlerStarted, please start it first.")
	}

	return nil
}

func (module CrawlerModule) runPipeline(signalC *chan bool, shard int) {

	var taskID []byte
	for {
		select {
		case <-*signalC:
			log.Trace("crawler exit, shard:", shard)
			return
		case taskID = <-queue.ReadChan(config.FetchChannel):
			stats.Increment("queue."+string(config.FetchChannel), "pop")
			id := string(taskID)
			log.Trace("shard:", shard, ",task received:", id)
			module.execute(id)
			log.Trace("shard:", shard, ",task finished:", id)
		}
	}
}

func (module CrawlerModule) execute(taskId string) {
	var pipeline *Pipeline
	defer func() {
		if !global.Env().IsDebug {
			if r := recover(); r != nil {
				if e, ok := r.(runtime.Error); ok {
					log.Error("pipeline: ", pipeline.GetID(), ", taskId: ", taskId, ", ", util.GetRuntimeErrorMessage(e))
				}
				log.Error("error in crawler,", util.ToJson(r, true), util.ToJson(pipeline.GetContext(), true))
			}
		}
	}()

	context := &Context{Phrase: config.PhraseCrawler}
	context.Set(CONTEXT_TASK_ID, taskId)

	if module.config.DefaultPipelineConfig == nil {
		panic("default pipeline config can't be null")
	}

	pipeline = NewPipelineFromConfig(module.config.DefaultPipelineConfig, context)
	pipeline.Run()

	if module.config.FetchThresholdInMs > 0 {
		log.Debug("sleep ", module.config.FetchThresholdInMs, "ms to control crawling speed")
		time.Sleep(time.Duration(module.config.FetchThresholdInMs) * time.Millisecond)
		log.Debug("wake up now,continue crawing")
	}

	log.Trace("end crawler")
}

type CrawlerModule struct {
	config *TaskConfig
}
