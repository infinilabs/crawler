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

package env

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	. "github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/core/types"
	"github.com/medcl/gopa/core/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type Env struct {
	//Logger logging.Logger
	*Registrar
	SystemConfig  *SystemConfig
	RuntimeConfig *RuntimeConfig
	Channels      *Channels

	ESClient util.ElasticsearchClient
}

func Environment(sysConfig SystemConfig) *Env {
	//if logger == nil {
	//	logger = logging.NullLogger{}
	//}

	env := Env{}

	env.SystemConfig = &sysConfig
	config, err := env.loadRuntimeConfig()
	if err != nil {
		panic(err)
	}
	env.RuntimeConfig = &config

	//override logging level
	if len(sysConfig.LogLevel) > 0 {
		env.RuntimeConfig.LoggingConfig.Level = sysConfig.LogLevel
	}

	env.Channels = &Channels{}
	env.Channels.pendingFetchUrl = make(chan types.PageTask, 1) //buffer number is 10
	env.Channels.pendingCheckUrl = make(chan types.PageTask, 1) //buffer number is 10
	env.Registrar = &Registrar{values: map[string]interface{}{}}

	env.init()

	return &env
}

func (this *Env) loadRuntimeConfig() (RuntimeConfig, error) {

	var configFile = "./gopa.yml"
	if this.SystemConfig != nil && len(this.SystemConfig.ConfigFile) > 0 {
		configFile = this.SystemConfig.ConfigFile
	}

	//load external yaml config
	filename, _ := filepath.Abs(configFile)
	var config RuntimeConfig

	if util.FileExists(filename) {
		log.Debug("load configFile:", filename)

		yamlFile, err := ioutil.ReadFile(filename)

		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(err)
		}
	} else {
		//init default Config
		config = RuntimeConfig{}
		config.PathConfig = (&PathConfig{}).Init()
		config.ClusterConfig = (&ClusterConfig{}).Init()
		config.LoggingConfig = (&LoggingConfig{}).Init()
		config.IndexingConfig = (&IndexingConfig{}).Init()
		config.ChannelConfig = (&ChannelConfig{}).Init()
		config.CrawlerConfig = (&CrawlerConfig{})       //.Init()
		config.ParserConfig = (&ParserConfig{})         //.Init()
		config.TaskConfig = (&TaskConfig{})             //.Init()
		config.RuledFetchConfig = (&RuledFetchConfig{}) //.Init()
	}

	//override built-in config
	config.PathConfig.SavedFileLog = config.PathConfig.Data + "/tasks/pending_parse.files"
	config.PathConfig.PendingFetchLog = config.PathConfig.Data + "/tasks/pending_fetch.urls"
	config.PathConfig.FetchFailedLog = config.PathConfig.Data + "/tasks/failed_fetch.urls"

	config.PathConfig.WebData = config.PathConfig.Data + "/web/"
	config.PathConfig.TaskData = config.PathConfig.Data + "/tasks/"
	config.PathConfig.QueueData = config.PathConfig.Data + "/queue/"

	config.TaskConfig.LinkUrlExtractRegex = regexp.MustCompile(config.TaskConfig.LinkUrlExtractRegexStr)
	config.TaskConfig.FetchUrlPattern = regexp.MustCompile(config.TaskConfig.FetchUrlPatternStr)
	config.TaskConfig.SavingUrlPattern = regexp.MustCompile(config.TaskConfig.SavingUrlPatternStr)
	config.TaskConfig.SkipPageParsePattern = regexp.MustCompile(config.TaskConfig.SkipPageParsePatternStr)

	return config, nil
}

func (this *Env) init() error {

	if this.RuntimeConfig.MaxGoRoutine < 1 {
		this.RuntimeConfig.MaxGoRoutine = 1
	}
	os.MkdirAll(this.RuntimeConfig.PathConfig.Data, 0777)
	os.MkdirAll(this.RuntimeConfig.PathConfig.Log, 0777)
	os.MkdirAll(this.RuntimeConfig.PathConfig.QueueData, 0777)
	os.MkdirAll(this.RuntimeConfig.PathConfig.WebData, 0777)
	os.MkdirAll(this.RuntimeConfig.PathConfig.TaskData, 0777)

	this.ESClient = util.ElasticsearchClient{Host: this.RuntimeConfig.IndexingConfig.Host, Index: this.RuntimeConfig.IndexingConfig.Index}
	this.Channels.pendingFetchDiskQueue = NewDiskQueue("pending_fetch", this.RuntimeConfig.PathConfig.QueueData, 100*1024*1024, 4, 1<<10, 2500, 10*time.Second)
	this.Channels.pendingCheckDiskQueue = NewDiskQueue("pending_check", this.RuntimeConfig.PathConfig.QueueData, 100*1024*1024, 4, 1<<10, 2500, 10*time.Second)

	return nil
}

func EmptyEnv() *Env {
	return &Env{}
}

type Channels struct {
	pendingCheckUrl chan types.PageTask //check if the url need to fetch or not
	pendingFetchUrl chan types.PageTask //the urls pending to fetch

	pendingFetchDiskQueue BackendQueue
	pendingCheckDiskQueue BackendQueue
}

func (this *Channels) PushUrlToCheck(url types.PageTask) error {
	//push chan first and then push diskqueue
	select {
	case this.pendingCheckUrl <- url:
		stats.Increment("global", stats.STATS_CHECKER_PUSH_CHAN_COUNT)
		return nil
	default:
		err := this.pendingCheckDiskQueue.Put(url.MustGetBytes())
		stats.Increment("global", stats.STATS_CHECKER_PUSH_DISK_COUNT)
		return err
	}
}

func (this *Channels) PushUrlToFetch(url types.PageTask) error {

	//push chan first and then push diskqueue
	select {
	case this.pendingFetchUrl <- url:
		stats.Increment("global", stats.STATS_FETCH_PUSH_CHAN_COUNT)
		return nil
	default:
		err := this.pendingFetchDiskQueue.Put(url.MustGetBytes())
		stats.Increment("global", stats.STATS_FETCH_PUSH_DISK_COUNT)
		return err
	}

}

func (this *Channels) PopUrlToCheck() (types.PageTask, error) {
	b := <-this.pendingCheckDiskQueue.ReadChan()
	url := types.PageTaskFromBytes(b)
	stats.Increment("global", stats.STATS_CHECKER_POP_DISK_COUNT)
	return url, nil
}

func (this *Channels) PopUrlToFetch() (types.PageTask, error) {
	b := <-this.pendingFetchDiskQueue.ReadChan() //:
	url := types.PageTaskFromBytes(b)
	stats.Increment("global", stats.STATS_FETCH_POP_DISK_COUNT)
	return url, nil
}

func (this *Channels) Close() {
	this.pendingFetchDiskQueue.Close()
	this.pendingCheckDiskQueue.Close()
}

//high priority config, init from the environment or startup, can't be changed
type SystemConfig struct {
	Version    string `0.0.1`
	ConfigFile string `gopa.yml`
	LogLevel   string `info`
}
