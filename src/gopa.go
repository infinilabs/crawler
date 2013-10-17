/**
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午4:45
 */
package main

import (
    config "config"
    "flag"
    "fmt"
    log "github.com/cihub/seelog"
    . "github.com/zeebo/sbloom"
    "hash/fnv"
    "io/ioutil"
//    "kafka"
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
)

var seedUrl string
var logLevel string
//var taskConfig *TaskConfig
//var kafkaConfig *config.KafkaConfig
//var bloomFilter *Filter
//var MaxGoRoutine int
var runtimeConfig config.RuntimeConfig


func persistBloomFilter(bloomFilterPersistFileName string) {

    //save bloom-filter
    m, err := runtimeConfig.BloomFilter.GobEncode()
    if err != nil {
        log.Error(err)
        return
    }
    err = ioutil.WriteFile(bloomFilterPersistFileName, m, 0600)
    if err != nil {
        panic(err)
        return
    }
    log.Info("bloomFilter safety persisted.")
}

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

func initOffset(runtimeConfig config.RuntimeConfig,typeName string, partition int) uint64 {
    log.Info("start init offsets,partition:", partition)

    path := runtimeConfig.TaskConfig.BaseStoragePath+"task/"+typeName + "_offset_" + strconv.FormatInt(int64(partition), 10)
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
        log.Info("init offsets successfully,", partition, ":", ret)
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
    //		log.Info("persist offset,",i,":",offsets[i].Offset,",",offsets[i].Partition)
    //	}

    //	log.Info("persist kafka offsets done")

    quit <- true
	log.Debug("finished shutting down")}

//parse config setting
func parseConfig() *TaskConfig{
    log.Debug("start parsing taskConfig")
	taskConfig := new(TaskConfig)
    taskConfig.LinkUrlExtractRegex = regexp.MustCompile(
        config.GetStringConfig("CrawlerRule", "LinkUrlExtractRegex", "(\\s+(src2|src|href|HREF|SRC))\\s*=\\s*[\"']?(.*?)[\"']"))

	taskConfig.ArrayStringSplitter=config.GetStringConfig("CrawlerRule","ArrayStringSplitter",",")
	taskConfig.SplitByUrlParameter=config.GetStringConfig("CrawlerRule","SplitByUrlParameter","p")


	taskConfig.GoProfEnabled=config.GetBoolConfig("CrawlerRule","GoProfEnabled",false)

	taskConfig.LinkUrlExtractRegexGroupIndex=config.GetIntConfig("CrawlerRule", "LinkUrlExtractRegexGroupIndex", 3)
    taskConfig.Name = config.GetStringConfig("CrawlerRule", "Name", "GopaTask")

	taskConfig.BaseStoragePath="data/"+taskConfig.Name+"/";

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



//    kafkaConfig = new(config.KafkaConfig)
//    kafkaConfig.Hostname = config.GetStringConfig("Kafka", "Hostname", "localhost:9092")
//    kafkaConfig.MaxSize = uint32(config.GetIntConfig("Kafka", "MaxSize", 1048576))

    // Setting taskConfig
    taskConfig.MaxGoRoutine = config.GetIntConfig("Global", "MaxGoRoutine", 1)

    if taskConfig.MaxGoRoutine  < 0 {
		taskConfig.MaxGoRoutine  = 1
    }

	log.Debug("MaxGoRoutine:",taskConfig.MaxGoRoutine)

	log.Debug("finished parsing taskConfig")
	return taskConfig
}

func initBloomFilter(bloomFilterPersistFileName string) *Filter {
	var bloomFilter *Filter
    //loading or initializing bloom filter
    if util.CheckFileExists(bloomFilterPersistFileName) {
        log.Debug("found bloomFilter,start reload")
        n, err := ioutil.ReadFile(bloomFilterPersistFileName)
        if err != nil {
            log.Error("bloomFilter", err)
        }
        if err := bloomFilter.GobDecode(n); err != nil {
            log.Error("bloomFilter", err)
        }
        log.Info("bloomFilter successfully reloaded")
    } else {
        probItems := config.GetIntConfig("BloomFilter", "ItemSize", 100000)
        log.Debug("initializing bloom-filter,virual size is,", probItems)
        bloomFilter = NewFilter(fnv.New64(), probItems)
        log.Info("bloomFilter successfully initialized")
    }
	return bloomFilter
}

func main() {
	flag.StringVar(&seedUrl, "seed", "http://example.com", "the seed url,where everything starts")
    flag.StringVar(&logLevel, "log", "info", "setting log level,options:trace,debug,info,warn,error")

    flag.Parse()

    defer log.Flush()

	runtimeConfig = config.RuntimeConfig{}

	runtime.GOMAXPROCS(runtime.NumCPU())

	util.SetInitLogging(logLevel)

	runtimeConfig.TaskConfig=parseConfig()

	logPath:=runtimeConfig.TaskConfig.BaseStoragePath+"log/gopa.log";
	util.SetLogging(logLevel,logPath)
	log.Info("[gopa] is on.")

	log.Debug("ArrayStringSplitter:",runtimeConfig.TaskConfig.ArrayStringSplitter)



	os.MkdirAll(runtimeConfig.TaskConfig.BaseStoragePath+     "task/",0777)
	os.MkdirAll(runtimeConfig.TaskConfig.BaseStoragePath+     "store/",0777)
	os.MkdirAll(runtimeConfig.TaskConfig.BaseStoragePath+      "log/",0777)

    bloomFilterPersistFileName := config.GetStringConfig("BloomFilter", "FileName", runtimeConfig.TaskConfig.BaseStoragePath+"task/bloomfilter.bin")

    if seedUrl == "" || seedUrl == "http://example.com" {
        log.Error("no seed was given. type:\"gopa -h\" for help.")
        os.Exit(1)
    }
	runtimeConfig.Storage=&fsstore.FsStore{}

	runtimeConfig.BloomFilter=initBloomFilter(bloomFilterPersistFileName)

    //	atr:="AZaz"
    //	btr:=[]byte(atr)
    //	fmt.Println(btr)
    //
    //	id:= getSeqStr([]byte("AA"),[]byte("ZZ"),false)
    //	fmt.Println(id)

	if runtimeConfig.TaskConfig.GoProfEnabled {
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

	maxGoRoutine := runtimeConfig.TaskConfig.MaxGoRoutine
	fetchQuitChannels := make([]*chan bool, maxGoRoutine) //kafkaQuitSignal signals for each go routing
	fetchOffsets := make([]*RoutingOffset, maxGoRoutine)  //kafka fetchOffsets

	parseQuitChannels := make([]*chan bool, 1) //kafkaQuitSignal signals for each go routing
	//	parseQuitChannels := make([]*chan bool, MaxGoRoutine) //kafkaQuitSignal signals for each go routing
	parseOffsets := make([]*RoutingOffset, maxGoRoutine) //kafka fetchOffsets

	kafkaQuitSignal := make(chan bool, 1)
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

			persistBloomFilter(bloomFilterPersistFileName)

			//wait kafka to exit
			log.Info("waiting kafka exit")
			go shutdown(fetchOffsets, fetchQuitChannels, parseOffsets, parseQuitChannels, kafkaQuitSignal)
			<-kafkaQuitSignal
			log.Info("kafka worker is down")
			finalQuitSignal <- true
		}
	}()

	//start fetcher
	for i := 0; i < maxGoRoutine; i++ {
		c := make(chan bool, 1)
		fetchQuitChannels[i] = &c
		offset := new(RoutingOffset)
		offset.Offset = initOffset(runtimeConfig, "fetch", i)
		offset.Partition = i
		fetchOffsets[i] = offset

		go task.FetchGo(runtimeConfig, &c, offset, i)
	}


	c2 := make(chan bool, 1)
	parseQuitChannels[0] = &c2
	offset2 := new(RoutingOffset)
	offset2.Offset = initOffset(runtimeConfig, "parse", 0)
	offset2.Partition = 0
	parseOffsets[0] = offset2
	pendingUrls := make(chan []byte)

	//sending feed to task queue
	go func() {
		//notice seed will not persisted to task queue
		log.Debug("sending feed to fetch queue,", seedUrl)
		//		broker1 := kafka.NewBrokerPublisher(kafkaConfig.Hostname, taskConfig.Name+"_fetch", 0)
		//		broker1.Publish(kafka.NewMessage([]byte(seedUrl)))
		//TODO replace to interface
		pendingUrls <- []byte(seedUrl)
	}()



	go task.ParseGo(pendingUrls, runtimeConfig, &c2, offset2)


	go func() {
		for {
			url := <-pendingUrls
			if !runtimeConfig.BloomFilter.Lookup(url) {
				randomPartition := 0
				if maxGoRoutine > 1 {
					randomPartition = rand.Intn(maxGoRoutine - 1)
				}
				log.Debug("publish:",string(url),",partition:",randomPartition)
//				publisher := kafka.NewBrokerPublisher(kafkaConfig.Hostname, taskConfig.Name+"_fetch", randomPartition)
//				publisher.Publish(kafka.NewMessage(url))
				runtimeConfig.Storage.TaskEnqueue(url)
				//TODO sharding
				runtimeConfig.BloomFilter.Add(url)
			}else{
				log.Trace("hit bloomfilter,ignore,",string(url))
			}

        }
    }()

    <-finalQuitSignal
    log.Info("[gopa] is down")
}


