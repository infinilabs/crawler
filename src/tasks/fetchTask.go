/** 
 * User: Medcl
 * Date: 13-7-25
 * Time: 下午9:18 
 */
package tasks

import (
	log "logging"
	"time"
	util "util"
	. "types"
)

//fetch url's content
func fetchUrl(url []byte, timeout time.Duration, runtimeConfig RuntimeConfig,  partition int) {
	t := time.NewTimer(timeout)
	defer t.Stop()

	resource := string(url)

	log.Debug("enter fetchUrl method:",resource)

	config:=runtimeConfig.TaskConfig

	if(runtimeConfig.Storage.CheckFetchedUrl(url)){
		return
	}


	path:=getSavedPath(runtimeConfig,url)

	if(runtimeConfig.Storage.CheckSavedFile(path)){
		log.Warn("file is already saved,skip fetch.",path)
		runtimeConfig.Storage.AddFetchedUrl(url)
		//todo re-parse local page
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

		body,err:=HttpGet(resource)

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


					Save(path, body)

				exitPage:
				} else {
					log.Debug("does not hit SavingUrlPattern ignoring,", resource)
				}
			}

			runtimeConfig.Storage.AddFetchedUrl(url)
			log.Debug("exit fetchUrl method:",resource)
		}else{
			runtimeConfig.Storage.AddFetchFailedUrl(url)
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

func FetchGo(runtimeConfig RuntimeConfig, taskC *chan []byte, quitC *chan bool, offsets *RoutingOffset, shard int) {
	 log.Info("fetch task started.shard:",shard)
	go func() {
		for {
			url := <-*taskC
			log.Debug("shard:",shard,",url received:", string(url))

				if !runtimeConfig.Storage.CheckFetchedUrl(url) {
					timeout := 10 * time.Second

//					if(fetchFilter.Lookup(url)){
//						log.Debug("hit fetch filter ,ignore,",string(url))
//						continue
//					}
//					fetchFilter.Add(url)

					log.Debug("shard:",shard,",url cool,start fetching:", string(url))
					fetchUrl(url, timeout, runtimeConfig, shard)

					//TODO
					//persist worker's offset
	//				path := config.BaseStoragePath+     "task/fetch_offset_" + strconv.FormatInt(int64(shard), 10) + ".tmp"
	//				path_new := config.BaseStoragePath+"task/fetch_offset_" + strconv.FormatInt(int64(shard), 10)
	//				fout, error := os.Create(path)
	//				if error != nil {
	//					log.Error(path, error)
	//					continue
	//				}

	//				defer fout.Close()
	//				log.Debug("partition:", shard, ",saved offset:", offsetV)
	//				fout.Write([]byte(strconv.FormatUint(msg.Offset(), 10)))
	//				utils.CopyFile(path, path_new)

				}else {
					log.Debug("hit fetch-bloomfilter,ignore,", string(url))
				}

		}
	}()

	<-*quitC
	log.Info("fetch task exit.shard:",shard)
}
