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
	"time"

	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/stats"
	task "github.com/medcl/gopa/core/tasks"
	util "github.com/medcl/gopa/core/util"
)

//fetch url's content
func fetchUrl(url []byte, timeout time.Duration, runtimeConfig *RuntimeConfig, offsets *RoutingParameter) {
	t := time.NewTimer(timeout)
	defer t.Stop()

	requestUrl := string(url)

	var storage = runtimeConfig.Storage

	log.Debug("enter fetchUrl method:", requestUrl)

	config := runtimeConfig.TaskConfig

	savePath := getSavedPath(runtimeConfig, url)

	if storage.FileHasSaved(savePath) {
		log.Warn("file already saved,skip fetch.", savePath)
		storage.AddSavedUrl(url)

		//re-parse local's previous saved page
		if runtimeConfig.ParseUrlsFromPreviousSavedPage {
			if !storage.FileHasParsed([]byte(savePath)) {
				log.Debug("previous saved page send to parse-queue:", savePath)
				storage.LogSavedFile(runtimeConfig.PathConfig.SavedFileLog, requestUrl+"|||"+savePath)
			}
		}
		storage.AddFetchedUrl(url)
		stats.Increment(stats.STATS_FETCH_IGNORE_COUNT)
		return
	}

	//checking fetchUrlPattern
	log.Debug("started check fetchUrlPattern,", config.FetchUrlPattern, ",", requestUrl)
	if config.FetchUrlPattern.Match(url) {
		log.Debug("match fetch url pattern,", requestUrl)
		if len(config.FetchUrlMustNotContain) > 0 {
			if util.ContainStr(requestUrl, config.FetchUrlMustNotContain) {
				log.Debug("hit FetchUrlMustNotContain,ignore,", requestUrl, " , ", config.FetchUrlMustNotContain)
				storage.AddFetchedUrl(url)
				stats.Increment(stats.STATS_FETCH_IGNORE_COUNT)
				return
			}
		}

		if len(config.FetchUrlMustContain) > 0 {
			if !util.ContainStr(requestUrl, config.FetchUrlMustContain) {
				log.Debug("not hit FetchUrlMustContain,ignore,", requestUrl, " , ", config.FetchUrlMustContain)
				storage.AddFetchedUrl(url)
				stats.Increment(stats.STATS_FETCH_IGNORE_COUNT)
				return
			}
		}
	} else {
		log.Debug("does not hit FetchUrlPattern ignoring,", requestUrl)
		storage.AddFetchedUrl(url)
		stats.Increment(stats.STATS_FETCH_IGNORE_COUNT)
		return
	}

	log.Debug("start fetch url,", requestUrl)
	flg := make(chan bool, 1)

	go func() {

		body, err := task.HttpGetWithCookie(requestUrl, config.Cookie)

		if err == nil {
			if body != nil {
				//todo parse urls from this page
				log.Debug("started check savingUrlPattern,", config.SavingUrlPattern, ",", string(url))
				if config.SavingUrlPattern.Match(url) {
					log.Debug("match saving url pattern,", requestUrl)
					if len(config.SavingUrlMustNotContain) > 0 {
						if util.ContainStr(requestUrl, config.SavingUrlMustNotContain) {
							log.Debug("hit SavingUrlMustNotContain,ignore,", requestUrl, " , ", config.SavingUrlMustNotContain)
							goto exitPage
						}
					}

					if len(config.SavingUrlMustContain) > 0 {
						if !util.ContainStr(requestUrl, config.SavingUrlMustContain) {
							log.Debug("not hit SavingUrlMustContain,ignore,", requestUrl, " , ", config.SavingUrlMustContain)
							goto exitPage
						}
					}

					_, err := Save(runtimeConfig, savePath, body)
					if err == nil {
						log.Info("saved:", savePath)
						//todo saved per shard
						storage.LogSavedFile(runtimeConfig.PathConfig.SavedFileLog, requestUrl+"|||"+savePath)
					} else {
						log.Debug("error while saved:", savePath, ",", err)
						goto exitPage
					}

				} else {
					log.Debug("does not hit SavingUrlPattern ignoring,", requestUrl)
				}
			}
			storage.AddFetchedUrl(url)
		exitPage:
			log.Debug("exit fetchUrl method:", requestUrl)
			storage.AddFetchedUrl(url)
			flg <- true
		} else {
			storage.LogFetchFailedUrl(runtimeConfig.PathConfig.FetchFailedLog, requestUrl)
			flg <- false
		}
	}()

	//监听通道，由于设有超时，不可能泄露
	select {
	case <-t.C:
		log.Error("fetching url time out,", requestUrl)
		stats.Increment(stats.STATS_FETCH_TIMEOUT_COUNT)
	case value := <-flg:
		if value {
			log.Debug("fetching url normal exit,", requestUrl)
			stats.Increment(stats.STATS_FETCH_SUCCESS_COUNT)
		} else {
			log.Debug("fetching url normal exit,", requestUrl)
			stats.Increment(stats.STATS_FETCH_FAIL_COUNT)
		}
		return
	}

}

func init() {}

func FetchGo(runtimeConfig *RuntimeConfig, taskC *chan []byte, quitC *chan bool, offsets *RoutingParameter) {

	shard := offsets.Shard

	go func() {
		for {
			url := <-*taskC

			if !runtimeConfig.Storage.UrlHasFetched(url) {

				log.Debug("shard:", shard, ",url received:", string(url))

				timeout := 10 * time.Second

				log.Info("shard:", shard, ",url cool,start fetching:", string(url))

				fetchUrl(url, timeout, runtimeConfig, offsets)

				if runtimeConfig.TaskConfig.FetchDelayThreshold > 0 {
					log.Debug("sleep ", runtimeConfig.TaskConfig.FetchDelayThreshold, "ms to control crawling speed")
					time.Sleep(time.Duration(runtimeConfig.TaskConfig.FetchDelayThreshold) * time.Millisecond)
					log.Debug("wake up now,continue crawing")
				}
			} else {
				log.Debug("shard:", shard, ",url received,but already fetched,skip: ", string(url))
			}

		}
	}()

	log.Trace("fetch task started.shard:", shard)

	<-*quitC

	log.Trace("fetch task exit.shard:", shard)

}

type FetchTask struct {
	innerTaskConfig *InnerTaskConfig
}

func (this *FetchTask) Init(config *InnerTaskConfig) {
	this.innerTaskConfig = config
}

func (this *FetchTask) Start() error {
	log.Info("fetch task is started, shard: ", this.innerTaskConfig.Parameter.Shard)
	FetchGo(this.innerTaskConfig.RuntimeConfig, this.innerTaskConfig.MessageChan, this.innerTaskConfig.QuitChan, this.innerTaskConfig.Parameter)
	return nil
}

func (this *FetchTask) Stop() error {
	log.Info("fetch task is stoped")
	return nil
}
