/** 
 * User: Medcl
 * Date: 13-7-25
 * Time: 下午9:18 
 */
package tasks

import (
	log "github.com/cihub/seelog"
	"time"
	util "github.com/medcl/gopa/src/util"
	. "github.com/medcl/gopa/src/config"
)

//fetch url's content
func fetchUrl(url []byte, timeout time.Duration, runtimeConfig *RuntimeConfig,  offsets *RoutingOffset) {
	t := time.NewTimer(timeout)
	defer t.Stop()

	resource := string(url)

	var storage = runtimeConfig.Storage

	log.Debug("enter fetchUrl method:",resource)

	config:=runtimeConfig.TaskConfig

	path:=getSavedPath(runtimeConfig,url)

	if(storage.CheckSavedFile(path)){
		log.Warn("file already saved,skip fetch.",path)
		storage.AddSavedUrl(url)

		//re-parse local's previous saved page
		if(runtimeConfig.ParseUrlsFromPreviousSavedPage){
			if(!storage.CheckParsedFile([]byte(path))){
				log.Debug("previous saved page send to parse-queue:",path)
				storage.LogSavedFile(runtimeConfig.PathConfig.SavedFileLog,resource+"|||"+path)
			}
		}
	   return
	}


	//checking fetchUrlPattern
	log.Debug("started check fetchUrlPattern,", config.FetchUrlPattern, ",", resource)
	if config.FetchUrlPattern.Match(url) {
		log.Debug("match fetch url pattern,", resource)
		if len(config.FetchUrlMustNotContain) > 0 {
			if util.ContainStr(resource, config.FetchUrlMustNotContain) {
				log.Debug("hit FetchUrlMustNotContain,ignore,", resource, " , ", config.FetchUrlMustNotContain)
				return
			}
		}

		if len(config.FetchUrlMustContain) > 0 {
			if !util.ContainStr(resource, config.FetchUrlMustContain) {
				log.Debug("not hit FetchUrlMustContain,ignore,", resource, " , ", config.FetchUrlMustContain)
				return
			}
		}
	} else {
		log.Debug("does not hit FetchUrlPattern ignoring,", resource)
		return
	}

	log.Debug("start fetch url,", resource)
	flg := make(chan bool, 1)

	go func() {

		body,err:=HttpGetWithCookie(resource,config.Cookie)

		if err == nil {
			if body != nil {
				//todo parse urls from this page
				log.Debug("started check savingUrlPattern,", config.SavingUrlPattern, ",", string(url))
				if config.SavingUrlPattern.Match(url) {
					log.Debug("match saving url pattern,", resource)
					if len(config.SavingUrlMustNotContain) > 0 {
						if util.ContainStr(resource, config.SavingUrlMustNotContain) {
							log.Debug("hit SavingUrlMustNotContain,ignore,", resource, " , ", config.SavingUrlMustNotContain)
							goto exitPage
						}
					}

					if len(config.SavingUrlMustContain) > 0 {
						if !util.ContainStr(resource, config.SavingUrlMustContain) {
							log.Debug("not hit SavingUrlMustContain,ignore,", resource, " , ", config.SavingUrlMustContain)
							goto exitPage
						}
					}


					_,err:=Save(runtimeConfig,path, body)
					if(err==nil){
						log.Info("saved:",path)
						//todo saved per shard
						storage.LogSavedFile(runtimeConfig.PathConfig.SavedFileLog,resource+"|||"+path)
					}else{
						log.Info("error while saved:",path,",",err)
						goto exitPage
					}

				} else {
					log.Debug("does not hit SavingUrlPattern ignoring,", resource)
				}
			}
			storage.AddFetchedUrl(url)
		exitPage:
			log.Debug("exit fetchUrl method:",resource)
		}else{
//			storage.AddFetchFailedUrl(url)
			storage.LogFetchFailedUrl(runtimeConfig.PathConfig.FetchFailedLog,resource)
		}
		flg <- true
	}()

	//监听通道，由于设有超时，不可能泄露
	select {
	case <-t.C:
		log.Error("fetching url time out,", resource)
	case <-flg:
		log.Debug("fetching url normal exit,", resource)
		return
	}

}

//var fetchFilter  *bloom.Filter
func init() {
//	fetchFilter = bloom.NewFilter(fnv.New64(), 1000000)
//	log.Warn("init bloom filter")
}

func FetchGo(runtimeConfig *RuntimeConfig, taskC *chan []byte, quitC *chan bool, offsets *RoutingOffset) {
	shard:=offsets.Shard
	log.Info("fetch task started.shard:",shard)
//	var storage = runtimeConfig.Storage
	go func() {
		for {
			url := <-*taskC
			log.Debug("shard:",shard,",url received:", string(url))

//				if !storage.CheckFetchedUrl(url) {
					timeout := 10 * time.Second

//					if(fetchFilter.Lookup(url)){
//						log.Debug("hit fetch filter ,ignore,",string(url))
//						continue
//					}
//					fetchFilter.Add(url)

					log.Debug("shard:",shard,",url cool,start fetching:", string(url))
					fetchUrl(url, timeout, runtimeConfig, offsets)
					if(runtimeConfig.TaskConfig.FetchDelayThreshold >0){
						log.Debug("sleep ",runtimeConfig.TaskConfig.FetchDelayThreshold,"ms to control crawling speed")
						time.Sleep( time.Duration (runtimeConfig.TaskConfig.FetchDelayThreshold) * time.Millisecond)
						log.Debug("wake up now,continue crawing")
					}
//
//				}else {
//					log.Debug("hit fetch-bloomfilter,ignore,", string(url))
//				}

		}
	}()

	log.Info("fetch task exit.shard:",shard)
	<-*quitC
}
