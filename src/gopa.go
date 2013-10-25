/**
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午4:45
 */
package main

import (
	"flag"
	"fmt"
	log "logging"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"util"
	"math/rand"
	"runtime"
	task "tasks"
	. "types"
	"net/http"
	fsstore "store/fs"
	config "config"
)

var seedUrl string
var logLevel string
var runtimeConfig RuntimeConfig

func getSeqStr(start []byte, end []byte, mix bool) []byte {
	if (len(start)) == len(end) {
		for i := range start {
			fmt.Println(start[i])
		}
		//		if(start>64 && end < 123){

		//		}
	}

	return nil
}

func init() {

}

func initOffset(runtimeConfig RuntimeConfig, typeName string, shard int) uint64 {
	log.Info("start init offsets,", typeName, ",shard:", shard)

	path := runtimeConfig.TaskConfig.TaskDataPath + "/" + typeName + "_offset_" + strconv.FormatInt(int64(shard), 10)
	if util.CheckFileExists(path) {
		log.Debug("found offset file,start loading")
		n, err := ioutil.ReadFile(path)
		if err != nil {
			log.Error("offset", err)
			return 0
		}
		ret, err := strconv.ParseInt(string(n), 10, 64)
		if err != nil {
			log.Error("offset", err)
			return 0
		}
		log.Info("init offsets successfully,", shard, ":", ret)
		return uint64(ret)
	}

	return 0
}

func shutdown(offsets []*RoutingOffset, quitChannels []*chan bool, offsets2 []*RoutingOffset, quitChannels2 []*chan bool, quit chan bool) {
	log.Debug("start shutting down")
	for i := range quitChannels {
		log.Debug("send exit signal to channel,", i)
		*quitChannels[i] <- true
	}

	for i := range quitChannels2 {
		log.Debug("send exit signal to channel,", i)
		*quitChannels2[i] <- true
	}

	log.Info("sent quit signal to go routings done")

	//	for i:=range offsets{
	//		//TODO
	//		log.Info("persist offset,",i,":",offsets[i].Offset,",",offsets[i].shard)
	//	}

	//	log.Info("persist kafka offsets done")

	quit <- true
	log.Debug("finished shutting down")}

//parse config setting
func parseConfig() *TaskConfig {
	log.Debug("start parsing taskConfig")
	taskConfig := new(TaskConfig)
	taskConfig.LinkUrlExtractRegex = regexp.MustCompile(
	config.GetStringConfig("CrawlerRule", "LinkUrlExtractRegex", "(\\s+(src2|src|href|HREF|SRC))\\s*=\\s*[\"']?(.*?)[\"']"))

	taskConfig.SplitByUrlParameter = config.GetStringConfig("CrawlerRule", "SplitByUrlParameter", "p")


	taskConfig.LinkUrlExtractRegexGroupIndex = config.GetIntConfig("CrawlerRule", "LinkUrlExtractRegexGroupIndex", 3)
	taskConfig.Name = config.GetStringConfig("CrawlerRule", "Name", "GopaTask")


	taskConfig.FollowSameDomain = config.GetBoolConfig("CrawlerRule", "FollowSameDomain", true)
	taskConfig.FollowSubDomain = config.GetBoolConfig("CrawlerRule", "FollowSubDomain", true)
	taskConfig.LinkUrlMustContain = config.GetStringConfig("CrawlerRule", "LinkUrlMustContain", "")
	taskConfig.LinkUrlMustNotContain = config.GetStringConfig("CrawlerRule", "LinkUrlMustNotContain", "")

	taskConfig.SkipPageParsePattern = regexp.MustCompile(config.GetStringConfig("CrawlerRule", "SkipPageParsePattern", ".*?\\.((js)|(css)|(rar)|(gz)|(zip)|(exe)|(bmp)|(jpeg)|(gif)|(png)|(jpg)|(apk))\\b")) //end with js,css,apk,zip,ignore

	taskConfig.FetchUrlPattern = regexp.MustCompile(config.GetStringConfig("CrawlerRule", "FetchUrlPattern", ".*"))
	taskConfig.FetchUrlMustContain = config.GetStringConfig("CrawlerRule", "FetchUrlMustContain", "")
	taskConfig.FetchUrlMustNotContain = config.GetStringConfig("CrawlerRule", "FetchUrlMustNotContain", "")

	taskConfig.SavingUrlPattern = regexp.MustCompile(config.GetStringConfig("CrawlerRule", "SavingUrlPattern", ".*"))
	taskConfig.SavingUrlMustContain = config.GetStringConfig("CrawlerRule", "SavingUrlMustContain", "")
	taskConfig.SavingUrlMustNotContain = config.GetStringConfig("CrawlerRule", "SavingUrlMustNotContain", "")


	taskConfig.TaskDataPath = config.GetStringConfig("CrawlerRule", "TaskData", runtimeConfig.PathConfig.TaskData + "/" + taskConfig.Name + "/")

	defaultWebDataPath := runtimeConfig.PathConfig.WebData + "/" + taskConfig.Name + "/"
	if (runtimeConfig.StoreWebPageTogether) {
		defaultWebDataPath = runtimeConfig.PathConfig.WebData
	}

	taskConfig.WebDataPath = config.GetStringConfig("CrawlerRule", "WebData", defaultWebDataPath)


	log.Debug("finished parsing taskConfig")
	return taskConfig
}


func main() {
	flag.StringVar(&seedUrl, "seed", "http://example.com", "the seed url,where everything starts")
	flag.StringVar(&logLevel, "log", "info", "setting log level,options:trace,debug,info,warn,error")

	flag.Parse()

	defer log.Flush()

	runtimeConfig = RuntimeConfig{}

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetInitLogging(logLevel)


	runtimeConfig.PathConfig = new(PathConfig)
	runtimeConfig.ClusterConfig = new(ClusterConfig)

	runtimeConfig.ClusterConfig.Name = config.GetStringConfig("cluster", "name", "gopa")

	// per cluster:data/gopa/
	runtimeConfig.PathConfig.Home = config.GetStringConfig("path", "home", "cluster/"+runtimeConfig.ClusterConfig.Name + "/")

	runtimeConfig.PathConfig.Data = config.GetStringConfig("path", "data", "")
	if (runtimeConfig.PathConfig.Data == "") {
		runtimeConfig.PathConfig.Data = runtimeConfig.PathConfig.Home + "/" + "data/"
	}

	runtimeConfig.PathConfig.Log = config.GetStringConfig("path", "log", "")
	if (runtimeConfig.PathConfig.Log == "") {
		runtimeConfig.PathConfig.Log = runtimeConfig.PathConfig.Home + "/" + "log/"
	}

	runtimeConfig.PathConfig.WebData = config.GetStringConfig("path", "webdata", "")
	if (runtimeConfig.PathConfig.WebData == "") {
		runtimeConfig.PathConfig.WebData = runtimeConfig.PathConfig.Data + "/" + "webdata/"
	}

	runtimeConfig.PathConfig.TaskData = config.GetStringConfig("path", "taskdata", "")
	if (runtimeConfig.PathConfig.TaskData == "") {
		runtimeConfig.PathConfig.TaskData = runtimeConfig.PathConfig.Data + "/" + "taskdata/"
	}

	runtimeConfig.StoreWebPageTogether = config.GetBoolConfig("Global", "StoreWebPageTogether", true)


	runtimeConfig.TaskConfig = parseConfig()


	//set default logging
	logPath := runtimeConfig.PathConfig.Log + "/" + runtimeConfig.TaskConfig.Name + "/gopa.log";
	log.SetLogging(logLevel, logPath)


	runtimeConfig.ParseUrlsFromSavedPage = config.GetBoolConfig("Switch", "ParseUrlsFromSavedPage", true)
	runtimeConfig.LoadTemplatedFetchJob = config.GetBoolConfig("Switch", "LoadTemplatedFetchJob", true)
	runtimeConfig.FetchUrlsFromSavedPage = config.GetBoolConfig("Switch", "FetchUrlsFromSavedPage", true)
	runtimeConfig.ParseUrlsFromPreviousSavedPage = config.GetBoolConfig("Switch", "ParseUrlsFromPreviousSavedPage", true)
	runtimeConfig.ArrayStringSplitter = config.GetStringConfig("CrawlerRule", "ArrayStringSplitter", ",")

	runtimeConfig.GoProfEnabled = config.GetBoolConfig("CrawlerRule", "GoProfEnabled", false)

	runtimeConfig.WalkBloomFilterFileName = config.GetStringConfig("BloomFilter", "WalkBloomFilterFileName", runtimeConfig.TaskConfig.TaskDataPath +   "/filters/walk.bloomfilter")
	runtimeConfig.FetchBloomFilterFileName = config.GetStringConfig("BloomFilter", "FetchBloomFilterFileName", runtimeConfig.TaskConfig.TaskDataPath + "/filters/fetch.bloomfilter")
	runtimeConfig.ParseBloomFilterFileName = config.GetStringConfig("BloomFilter", "ParseBloomFilterFileName", runtimeConfig.TaskConfig.TaskDataPath + "/filters/parse.bloomfilter")
	runtimeConfig.PendingFetchBloomFilterFileName = config.GetStringConfig("BloomFilter", "PendingFetchBloomFilterFileName", runtimeConfig.TaskConfig.TaskDataPath + "/filters/pending_fetch.bloomfilter")

	runtimeConfig.PathConfig.SavedFileLog=runtimeConfig.TaskConfig.TaskDataPath+"/tasks/pending_parse.files"
	runtimeConfig.PathConfig.PendingFetchLog=runtimeConfig.TaskConfig.TaskDataPath+"/tasks/pending_fetch.urls"
	runtimeConfig.PathConfig.FetchFailedLog=runtimeConfig.TaskConfig.TaskDataPath+"/tasks/failed_fetch.urls"

	runtimeConfig.MaxGoRoutine = config.GetIntConfig("Global", "MaxGoRoutine", 1)
	if runtimeConfig.MaxGoRoutine < 0 {
		runtimeConfig.MaxGoRoutine = 1
	}

	log.Debug("maxGoRoutine:", runtimeConfig.MaxGoRoutine)
	log.Debug("path.home:", runtimeConfig.PathConfig.Home)


	os.MkdirAll(runtimeConfig.PathConfig.Home, 0777)
	os.MkdirAll(runtimeConfig.PathConfig.Data, 0777)
	os.MkdirAll(runtimeConfig.PathConfig.Log, 0777)
	os.MkdirAll(runtimeConfig.PathConfig.WebData, 0777)
	os.MkdirAll(runtimeConfig.PathConfig.TaskData, 0777)

	os.MkdirAll(runtimeConfig.TaskConfig.TaskDataPath, 0777)
	os.MkdirAll(runtimeConfig.TaskConfig.TaskDataPath + "/tasks/", 0777)
	os.MkdirAll(runtimeConfig.TaskConfig.TaskDataPath + "/filters/", 0777)
	os.MkdirAll(runtimeConfig.TaskConfig.TaskDataPath + "/urls/", 0777)
	os.MkdirAll(runtimeConfig.TaskConfig.WebDataPath, 0777)


	if seedUrl == "" || seedUrl == "http://example.com" {
		log.Error("no seed was given. type:\"gopa -h\" for help.")
		os.Exit(1)
	}


	log.Info("[gopa] is on.")



	runtimeConfig.Storage = &fsstore.FsStore{}

	runtimeConfig.Storage.InitWalkBloomFilter(runtimeConfig.WalkBloomFilterFileName);
	runtimeConfig.Storage.InitFetchBloomFilter(runtimeConfig.FetchBloomFilterFileName);
	runtimeConfig.Storage.InitParseBloomFilter(runtimeConfig.ParseBloomFilterFileName);
	runtimeConfig.Storage.InitPendingFetchBloomFilter(runtimeConfig.PendingFetchBloomFilterFileName);

	//	atr:="AZaz"
	//	btr:=[]byte(atr)
	//	fmt.Println(btr)
	//
	//	id:= getSeqStr([]byte("AA"),[]byte("ZZ"),false)
	//	fmt.Println(id)

	if runtimeConfig.GoProfEnabled {
		//pprof serves
		go func() {
			log.Info(http.ListenAndServe("localhost:6060", nil))
			log.Info("pprof server is up,http://localhost:6060/debug/pprof")
		}()
	}

	//adding default http protocol
	if !strings.HasPrefix(seedUrl, "http") {
		seedUrl = "http://" + seedUrl
	}

	maxGoRoutine := runtimeConfig.MaxGoRoutine
	fetchQuitChannels := make([]*chan bool, maxGoRoutine) //shutdownSignal signals for each go routing
	fetchTaskChannels := make([]*chan []byte, maxGoRoutine) //fetchTask channels
	fetchOffsets := make([]*RoutingOffset, maxGoRoutine)  //kafka fetchOffsets

	parseQuitChannels := make([]*chan bool, 2) //shutdownSignal signals for each go routing
	//	parseQuitChannels := make([]*chan bool, MaxGoRoutine) //shutdownSignal signals for each go routing
	parseOffsets := make([]*RoutingOffset, maxGoRoutine) //kafka fetchOffsets

	shutdownSignal := make(chan bool, 1)
	finalQuitSignal := make(chan bool, 1)


	//handle exit event
	exitEventChannel := make(chan os.Signal, 1)
	signal.Notify(exitEventChannel, syscall.SIGINT)
	signal.Notify(exitEventChannel, os.Interrupt)
	go func() {
		s := <-exitEventChannel
		log.Debug("got signal:", s)
		if s == os.Interrupt || s.(os.Signal) == syscall.SIGINT {
			log.Warn("got signal:os.Interrupt,saving data and exit")
			//			defer  os.Exit(0)

			runtimeConfig.Storage.PersistBloomFilter()

			//wait workers to exit
			log.Info("waiting workers exit")
			go shutdown(fetchOffsets, fetchQuitChannels, parseOffsets, parseQuitChannels, shutdownSignal)
			<-shutdownSignal
			log.Info("workers shutdown")
			finalQuitSignal <- true
		}
	}()

	//start fetcher
	for i := 0; i < maxGoRoutine; i++ {
		quitC := make(chan bool, 1)
		taskC := make(chan []byte)

		fetchQuitChannels[i] = &quitC
		fetchTaskChannels[i] = &taskC
		offset := new(RoutingOffset)
		offset.Offset = initOffset(runtimeConfig, "fetch", i)
		offset.Shard = i
		fetchOffsets[i] = offset

		go task.FetchGo(runtimeConfig, &taskC, &quitC, offset, i)
	}


	c2 := make(chan bool, 1)
	parseQuitChannels[0] = &c2
	offset2 := new(RoutingOffset)
	offset2.Offset = initOffset(runtimeConfig, "parse", 0)
	offset2.Shard = 0
	parseOffsets[0] = offset2
	pendingFetchUrls := make(chan []byte)

	//fetch rule:all urls -> persisted to sotre -> fetched from store -> pushed to pendingFetchUrls -> redistributed to sharded goroutines -> fetch -> save webpage to store -> done
	//parse rule:url saved to store -> local path persisted to store -> fetched to pendingParseFiles -> redistributed to sharded goroutines -> parse -> clean urls -> enqueue to url store ->done

	//sending feed to task queue
	go func() {
		//notice seed will not been persisted
		log.Debug("sending feed to fetch queue,", seedUrl)
		pendingFetchUrls <- []byte(seedUrl)
	}()

	//start local saved file parser
	if runtimeConfig.ParseUrlsFromSavedPage{
		go task.ParseGo(pendingFetchUrls, runtimeConfig, &c2, offset2)
	}


	//redistribute pendingFetchUrls to sharded workers
	go func() {
		for {
			url := <-pendingFetchUrls
			if !runtimeConfig.Storage.CheckWalkedUrl(url) {

				if (runtimeConfig.Storage.CheckFetchedUrl(url)) {
					log.Warn("dont hit walk bloomfilter but hit fetch bloomfilter,also ignore,", string(url))
					runtimeConfig.Storage.AddWalkedUrl(url)
					continue
				}

				randomShard := 0
				if maxGoRoutine > 1 {
					randomShard = rand.Intn(maxGoRoutine - 1)
				}
				log.Debug("publish:", string(url), ",shard:", randomShard)
				runtimeConfig.Storage.AddWalkedUrl(url)
				*fetchTaskChannels[randomShard] <- url
			}else {
				log.Trace("hit walk or fetch bloomfilter,just ignore,", string(url))
			}
		}
	}()

	//load predefined fetch jobs
	if runtimeConfig.LoadTemplatedFetchJob{
		go func() {

			if (util.CheckFileExists(runtimeConfig.TaskConfig.TaskDataPath + "/urls/template.txt")) {

				templates := util.ReadAllLines(runtimeConfig.TaskConfig.TaskDataPath + "/urls/template.txt")
				ids := util.ReadAllLines(runtimeConfig.TaskConfig.TaskDataPath + "/urls/id.txt")


				for _, id := range ids {
					for _, template := range templates {
						log.Trace("id:", id)
						log.Trace("template:", template)
						url := strings.Replace(template, "{id}", id, -1)
						log.Debug("new task from template:", url)
						pendingFetchUrls <- []byte(url)
					}
				}
				log.Info("templated download is done.")

			}

		}()
	}


	//fetch urls from saved pages
	if runtimeConfig.FetchUrlsFromSavedPage{
		c3 := make(chan bool, 1)
		parseQuitChannels[1] = &c3
		offset3 := new(RoutingOffset)
		offset3.Offset = initOffset(runtimeConfig, "fetch_from_saved", 0)
		offset3.Shard = 0
		parseOffsets[1] = offset3
		go task.LoadTaskFromLocalFile(pendingFetchUrls, runtimeConfig, &c3, offset3)
	}


	//parse fetch failed jobs,and will ignore the walk-filter
	//TODO

	<-finalQuitSignal
	log.Info("[gopa] is down")
}


