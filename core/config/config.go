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

package config

import (
	"os"
	"regexp"

	log "github.com/cihub/seelog"
	cfg "github.com/robfig/config"
)

var loadingConfig *cfg.Config
var runtimeConfig RuntimeConfig


func OldGetRuntimeConfig() *RuntimeConfig {

	log.Trace("start init config")

	//parse main config
	loadingConfig, _ = cfg.ReadDefault("config.ini")
	runtimeConfig.PathConfig = new(PathConfig)
	runtimeConfig.ClusterConfig = new(ClusterConfig)
	parseConfig()

	runtimeConfig.ClusterConfig.Name = GetStringConfig("cluster", "name", "gopa")

	runtimeConfig.TaskConfig = parseConfig()

	//runtimeConfig.ParseUrlsFromSavedFileLog = GetBoolConfig("Switch", "ParseUrlsFromSavedFileLog", true)
	runtimeConfig.LoadTemplatedFetchJob = GetBoolConfig("Switch", "LoadTemplatedFetchJob", true)
	runtimeConfig.LoadRuledFetchJob = GetBoolConfig("Switch", "LoadRuledFetchJob", false)
	//runtimeConfig.LoadPendingFetchJobs = GetBoolConfig("Switch", "LoadPendingFetchJobs", true)
	runtimeConfig.HttpEnabled = GetBoolConfig("Switch", "HttpEnabled", true)
	//runtimeConfig.ParseUrlsFromPreviousSavedPage = GetBoolConfig("Switch", "ParseUrlsFromPreviousSavedPage", false)
	//runtimeConfig.ArrayStringSplitter = GetStringConfig("CrawlerRule", "ArrayStringSplitter", ",")

	//runtimeConfig.WalkBloomFilterFileName = GetStringConfig("BloomFilter", "WalkBloomFilterFileName", runtimeConfig.TaskConfig.TaskDataPath+"/filters/walk.bloomfilter")
	//runtimeConfig.FetchBloomFilterFileName = GetStringConfig("BloomFilter", "FetchBloomFilterFileName", runtimeConfig.TaskConfig.TaskDataPath+"/filters/fetch.bloomfilter")
	//runtimeConfig.ParseBloomFilterFileName = GetStringConfig("BloomFilter", "ParseBloomFilterFileName", runtimeConfig.TaskConfig.TaskDataPath+"/filters/parse.bloomfilter")
	//runtimeConfig.PendingFetchBloomFilterFileName = GetStringConfig("BloomFilter", "PendingFetchBloomFilterFileName", runtimeConfig.TaskConfig.TaskDataPath+"/filters/pending_fetch.bloomfilter")

	runtimeConfig.MaxGoRoutine = GetIntConfig("Global", "MaxGoRoutine", 1)
	if runtimeConfig.MaxGoRoutine < 2 {
		runtimeConfig.MaxGoRoutine = 2
	}

	//os.MkdirAll(runtimeConfig.PathConfig.Home, 0777)
	os.MkdirAll(runtimeConfig.PathConfig.Data, 0777)
	os.MkdirAll(runtimeConfig.PathConfig.Log, 0777)
	os.MkdirAll(runtimeConfig.PathConfig.WebData, 0777)
	os.MkdirAll(runtimeConfig.PathConfig.TaskData, 0777)

	os.MkdirAll(runtimeConfig.TaskConfig.TaskDataPath, 0777)
	os.MkdirAll(runtimeConfig.TaskConfig.TaskDataPath+"/tasks/", 0777)
	os.MkdirAll(runtimeConfig.TaskConfig.TaskDataPath+"/filters/", 0777)
	os.MkdirAll(runtimeConfig.TaskConfig.TaskDataPath+"/urls/", 0777)
	//os.MkdirAll(runtimeConfig.TaskConfig.WebDataPath, 0777)

	runtimeConfig.RuledFetchConfig = new(RuledFetchConfig)
	runtimeConfig.RuledFetchConfig.UrlTemplate = GetStringConfig("RuledFetch", "UrlTemplate", "")
	runtimeConfig.RuledFetchConfig.From = GetIntConfig("RuledFetch", "From", 0)
	runtimeConfig.RuledFetchConfig.To = GetIntConfig("RuledFetch", "To", 10)
	runtimeConfig.RuledFetchConfig.Step = GetIntConfig("RuledFetch", "Step", 1)
	runtimeConfig.RuledFetchConfig.LinkExtractPattern = GetStringConfig("RuledFetch", "LinkExtractPattern", "")
	runtimeConfig.RuledFetchConfig.LinkTemplate = GetStringConfig("RuledFetch", "LinkTemplate", "")

	return &runtimeConfig
}

//parse config setting
func parseConfig() *TaskConfig {
	log.Debug("start parsing taskConfig")
	taskConfig := new(TaskConfig)
	taskConfig.LinkUrlExtractRegex = regexp.MustCompile(
		GetStringConfig("CrawlerRule", "LinkUrlExtractRegex", "(\\s+(src2|src|href|HREF|SRC))\\s*=\\s*[\"']?(.*?)[\"']"))

	//taskConfig.SplitByUrlParameter = GetStringConfig("CrawlerRule", "SplitByUrlParameter", "p,pn,page,start,index")

	taskConfig.LinkUrlExtractRegexGroupIndex = GetIntConfig("CrawlerRule", "LinkUrlExtractRegexGroupIndex", 3)
	//taskConfig.Name = GetStringConfig("CrawlerRule", "Name", "GopaTask")

	taskConfig.FollowSameDomain = GetBoolConfig("CrawlerRule", "FollowSameDomain", true)
	taskConfig.FollowSubDomain = GetBoolConfig("CrawlerRule", "FollowSubDomain", true)
	taskConfig.LinkUrlMustContain = GetStringConfig("CrawlerRule", "LinkUrlMustContain", "")
	taskConfig.LinkUrlMustNotContain = GetStringConfig("CrawlerRule", "LinkUrlMustNotContain", "")

	taskConfig.SkipPageParsePattern = regexp.MustCompile(GetStringConfig("CrawlerRule", "SkipPageParsePattern", ".*?\\.((js)|(css)|(rar)|(gz)|(zip)|(exe)|(bmp)|(jpeg)|(gif)|(png)|(jpg)|(apk))\\b")) //end with js,css,apk,zip,ignore

	taskConfig.FetchUrlPattern = regexp.MustCompile(GetStringConfig("CrawlerRule", "FetchUrlPattern", ".*"))
	taskConfig.FetchUrlMustContain = GetStringConfig("CrawlerRule", "FetchUrlMustContain", "")
	taskConfig.FetchUrlMustNotContain = GetStringConfig("CrawlerRule", "FetchUrlMustNotContain", "")

	taskConfig.SavingUrlPattern = regexp.MustCompile(GetStringConfig("CrawlerRule", "SavingUrlPattern", ".*"))
	taskConfig.SavingUrlMustContain = GetStringConfig("CrawlerRule", "SavingUrlMustContain", "")
	taskConfig.SavingUrlMustNotContain = GetStringConfig("CrawlerRule", "SavingUrlMustNotContain", "")

	taskConfig.Cookie = GetStringConfig("CrawlerRule", "Cookie", "")
	taskConfig.FetchDelayThreshold = GetIntConfig("CrawlerRule", "FetchDelayThreshold", 0)

	log.Debug("finished parsing taskConfig")
	return taskConfig
}

func GetStringConfig(configSection string, configKey string, defaultValue string) string {
	if loadingConfig == nil {
		log.Trace("loadingConfig is nil,just return")
		return defaultValue
	}

	//loading or initializing bloom filter
	value, error := loadingConfig.String(configSection, configKey)
	if error != nil {
		value = defaultValue
	}
	log.Trace("get config value,", configSection, ".", configKey, ":", value)
	return value
}

func GetFloatConfig(configSection string, configKey string, defaultValue float64) float64 {
	if loadingConfig == nil {
		log.Trace("loadingConfig is nil,just return")
		return defaultValue
	}

	//loading or initializing bloom filter
	value, error := loadingConfig.Float(configSection, configKey)
	if error != nil {
		value = defaultValue
	}
	log.Trace("get config value,", configSection, ".", configKey, ":", value)
	return value
}

func GetIntConfig(configSection string, configKey string, defaultValue int) int {
	if loadingConfig == nil {
		log.Trace("loadingConfig is nil,just return")
		return defaultValue
	}

	//loading or initializing bloom filter
	value, error := loadingConfig.Int(configSection, configKey)
	if error != nil {
		value = defaultValue
	}
	log.Trace("get config value,", configSection, ".", configKey, ":", value)
	return value
}

func GetBoolConfig(configSection string, configKey string, defaultValue bool) bool {
	if loadingConfig == nil {
		log.Trace("loadingConfig is nil,just return")
		return defaultValue
	}

	//loading or initializing bloom filter
	value, error := loadingConfig.Bool(configSection, configKey)
	if error != nil {
		value = defaultValue
	}
	log.Trace("get config value,", configSection, ".", configKey, ":", value)
	return value
}
