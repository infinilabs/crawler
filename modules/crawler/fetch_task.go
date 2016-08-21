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
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/core/util"
	"github.com/medcl/gopa/core/types"
	"strings"
)

func checkIfUrlWillBeSave(runtimeConfig *RuntimeConfig,url []byte,)bool  {

	requestUrl:=string(url)

	log.Debug("started check savingUrlPattern,", runtimeConfig.TaskConfig.SavingUrlPattern, ",", string(url))
	if runtimeConfig.TaskConfig.SavingUrlPattern.Match(url) {


		log.Debug("match saving url pattern,", requestUrl)
		if len(runtimeConfig.TaskConfig.SavingUrlMustNotContain) > 0 {
			if util.ContainStr(requestUrl, runtimeConfig.TaskConfig.SavingUrlMustNotContain) {
				log.Debug("hit SavingUrlMustNotContain,ignore,", requestUrl, " , ", runtimeConfig.TaskConfig.SavingUrlMustNotContain)
				return false
			}
		}

		if len(runtimeConfig.TaskConfig.SavingUrlMustContain) > 0 {
			if !util.ContainStr(requestUrl, runtimeConfig.TaskConfig.SavingUrlMustContain) {
				log.Debug("not hit SavingUrlMustContain,ignore,", requestUrl, " , ", runtimeConfig.TaskConfig.SavingUrlMustContain)
				return false
			}
		}

		return true

	} else {
		log.Debug("does not hit SavingUrlPattern ignoring,", requestUrl)
	}
	return false
}

//fetch url's content
func fetchUrl(url []byte, timeout time.Duration, runtimeConfig *RuntimeConfig) {
	t := time.NewTimer(timeout)
	defer t.Stop()
	requestUrl := string(url)

	if(url==nil||len(requestUrl)==0){
		log.Error("invalid fetchUrl")
		return
	}

	//adding default http protocol
	if !strings.HasPrefix(requestUrl, "http") {
		requestUrl = "http://" + requestUrl
	}



	var storage = runtimeConfig.Storage

	log.Debug("enter fetchUrl method:", requestUrl)

	config := runtimeConfig.TaskConfig

	saveDir,saveFile := getSavedPath(runtimeConfig, url)

	savePath:=saveDir+saveFile

	if storage.FileHasSaved(savePath) {
		log.Warn("file already saved,skip fetch.", savePath)
		storage.AddSavedUrl(url)
		log.Debug("file add to saved log")

		//re-parse local's previous saved page
		if runtimeConfig.ParserConfig.ReParseUrlsFromPreviousSavedPage {
			if !storage.FileHasParsed([]byte(savePath)) {
				log.Debug("previous saved page send to parse-queue:", savePath)
				storage.LogSavedFile(runtimeConfig.PathConfig.SavedFileLog, requestUrl+"|||"+savePath)
			}
		}
		storage.AddFetchedUrl(url)
		log.Debug("file add to fetched log")
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
		treasure:=types.Treasure{}
		treasure.CreateTime=time.Now().UTC()
		treasure.LastCheckTime=time.Now().UTC()

		body, err := util.HttpGetWithCookie(&treasure,requestUrl, config.Cookie)

		if err == nil {
			if body != nil {
				if(treasure.StatusCode==404||treasure.StatusCode==302){
					log.Error("error while 404 or 302:", requestUrl," ",treasure.StatusCode)
					flg <- false
					return
				}

				//check save rules
				if(checkIfUrlWillBeSave(runtimeConfig,url)){
					_, err := Save(runtimeConfig, saveDir,saveFile, body)

					treasure.Body=string(body)
					treasure.Size=len(body)
					treasure.Snapshot=savePath

					//data,_:=json.Marshal(treasure)
					//log.Error(string(data))

					if err == nil {
						log.Info("saved:", savePath)

						runtimeConfig.Storage.LogSavedFile(runtimeConfig.PathConfig.SavedFileLog, requestUrl+"|||"+savePath)
					} else {
						log.Error("error while saved:", savePath, ",", err)
						flg <- false
						return
					}
				}
			}

			storage.AddFetchedUrl(url)
			log.Debug("exit fetchUrl method:", requestUrl)
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
			log.Debug("fetching url error exit,", requestUrl)
			stats.Increment(stats.STATS_FETCH_FAIL_COUNT)
		}
		return
	}

}

func FetchGo(env *Env, quitC *chan bool, shard int) {

	go func() {
		for {
			log.Debug("ready to receive url")
			url := <-env.Channels.PendingFetchUrl

			log.Debug("shard:", shard, ",url received:", string(url))

			if !env.RuntimeConfig.Storage.UrlHasFetched(url) {


				timeout := 10 * time.Second

				log.Info("shard:", shard, ",url cool,start fetching:", string(url))

				fetchUrl(url, timeout, env.RuntimeConfig)

				if env.RuntimeConfig.TaskConfig.FetchDelayThreshold > 0 {
					log.Debug("sleep ", env.RuntimeConfig.TaskConfig.FetchDelayThreshold, "ms to control crawling speed")
					time.Sleep(time.Duration(env.RuntimeConfig.TaskConfig.FetchDelayThreshold) * time.Millisecond)
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
