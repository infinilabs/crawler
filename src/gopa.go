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
    "kafka"
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

)

var seedUrl string
var logLevel string
var taskConfig *TaskConfig
var kafkaConfig *config.KafkaConfig
var bloomFilter *Filter
var MaxGoRoutine int

func persistBloomFilter(bloomFilterPersistFileName string) {

    //save bloom-filter
    m, err := bloomFilter.GobEncode()
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

func initOffset(typeName string, partition int) uint64 {
    log.Info("start init offsets,partition:", partition)

    path := taskConfig.BaseStoragePath+"task/"+typeName + "_offset_" + strconv.FormatInt(int64(partition), 10)
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

func closeKafkaConsumer(offsets []*RoutingOffset, quitChannels []*chan bool, offsets2 []*RoutingOffset, quitChannels2 []*chan bool, quit chan bool) {

    for i := range quitChannels {
        log.Debug("send exit signal to channel,", i)
        *quitChannels[i] <- true
    }

    for i := range quitChannels2 {
        log.Debug("send exit signal to channel,", i)
        *quitChannels[i] <- true
    }

    log.Info("sent quit signal to go routings done")

    //	for i:=range offsets{
    //		//TODO
    //		log.Info("persist offset,",i,":",offsets[i].Offset,",",offsets[i].Partition)
    //	}

    //	log.Info("persist kafka offsets done")

    quit <- true
}

//parse config setting
func parseConfig() {
    taskConfig = new(TaskConfig)
    taskConfig.LinkUrlExtractRegex = regexp.MustCompile(
        config.GetStringConfig("CrawlerRule", "LinkUrlExtractRegex", "(src2|src|href|HREF|SRC)\\s*=\\s*[\"']?(.*?)[\"']"))

	taskConfig.ArrayStringSplitter=config.GetStringConfig("CrawlerRule","ArrayStringSplitter","##")
	taskConfig.SplitByUrlParameter=config.GetStringConfig("CrawlerRule","SplitByUrlParameter","p")


	taskConfig.GoProfEnabled=config.GetBoolConfig("CrawlerRule","GoProfEnabled",false)

	taskConfig.LinkUrlExtractRegexGroupIndex=config.GetIntConfig("CrawlerRule", "LinkUrlExtractRegexGroupIndex", 2)
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



    kafkaConfig = new(config.KafkaConfig)
    kafkaConfig.Hostname = config.GetStringConfig("Kafka", "Hostname", "localhost:9092")
    kafkaConfig.MaxSize = uint32(config.GetIntConfig("Kafka", "MaxSize", 1048576))

    // Setting taskConfig
    MaxGoRoutine = config.GetIntConfig("Global", "MaxGoRoutine", 1)

    if MaxGoRoutine < 0 {
        MaxGoRoutine = 1
    }
}

func initBloomFilter(bloomFilterPersistFileName string) {
    //loading or initializing bloom filter
    if util.CheckFileExists(bloomFilterPersistFileName) {
        log.Debug("found bloomFilter,start reload")
        n, err := ioutil.ReadFile(bloomFilterPersistFileName)
        if err != nil {
            log.Error("bloomFilter", err)
            return
        }
        bloomFilter = new(Filter)
        if err := bloomFilter.GobDecode(n); err != nil {
            log.Error("bloomFilter", err)
            return
        }

        log.Info("bloomFilter successfully reloaded")
    } else {
        probItems := config.GetIntConfig("BloomFilter", "ItemSize", 100000)
        log.Debug("initializing bloom-filter,virual size is,", probItems)
        bloomFilter = NewFilter(fnv.New64(), probItems)
        log.Info("bloomFilter successfully initialized")
    }

}

func main() {

    flag.StringVar(&seedUrl, "seed", "http://example.com", "the seed url,where everything begins")
    flag.StringVar(&logLevel, "log", "info", "setting log level,options:trace,debug,info,warn,error")

    flag.Parse()

    defer log.Flush()

    runtime.GOMAXPROCS(2)

    parseConfig()

	setLogging()
	log.Info("[gopa] is on.")

	log.Debug("ArrayStringSplitter:",taskConfig.ArrayStringSplitter)
	log.Debug("MaxGoRoutine:",MaxGoRoutine)



	os.MkdirAll(taskConfig.BaseStoragePath+     "task/",0777)
	os.MkdirAll(taskConfig.BaseStoragePath+     "store/",0777)
	os.MkdirAll(taskConfig.BaseStoragePath+      "log/",0777)

    bloomFilterPersistFileName := config.GetStringConfig("BloomFilter", "FileName", taskConfig.BaseStoragePath+"task/bloomfilter.bin")

    if seedUrl == "" || seedUrl == "http://example.com" {
        log.Error("no seed was given. type:\"gopa -h\" for help.")
        os.Exit(1)
    }

    initBloomFilter(bloomFilterPersistFileName)

    //	atr:="AZaz"
    //	btr:=[]byte(atr)
    //	fmt.Println(btr)
    //
    //	id:= getSeqStr([]byte("AA"),[]byte("ZZ"),false)
    //	fmt.Println(id)

	if taskConfig.GoProfEnabled {
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

    fetchQuitChannels := make([]*chan bool, MaxGoRoutine) //kafkaQuitSignal signals for each go routing
    fetchOffsets := make([]*RoutingOffset, MaxGoRoutine)  //kafka fetchOffsets

    parseQuitChannels := make([]*chan bool, 1) //kafkaQuitSignal signals for each go routing
    //	parseQuitChannels := make([]*chan bool, MaxGoRoutine) //kafkaQuitSignal signals for each go routing
    parseOffsets := make([]*RoutingOffset, MaxGoRoutine) //kafka fetchOffsets

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
            go closeKafkaConsumer(fetchOffsets, fetchQuitChannels, parseOffsets, parseQuitChannels, kafkaQuitSignal)
            <-kafkaQuitSignal
            log.Info("kafka worker is down")
            finalQuitSignal <- true
        }
    }()

    //start parse local files' task
    go func() {
        log.Error("sending feed to fetch queue,", seedUrl)
        broker1 := kafka.NewBrokerPublisher(kafkaConfig.Hostname, taskConfig.Name+"_fetch", 0)
        broker1.Publish(kafka.NewMessage([]byte(seedUrl)))
    }()

    for i := 0; i < MaxGoRoutine; i++ {
        c := make(chan bool, 1)
        fetchQuitChannels[i] = &c
        offset := new(RoutingOffset)
        offset.Offset = initOffset("fetch", i)
        offset.Partition = i
        fetchOffsets[i] = offset

        go task.Fetch(bloomFilter, taskConfig, kafkaConfig, &c, offset, i)
    }

    c2 := make(chan bool, 1)
    parseQuitChannels[0] = &c2
    offset2 := new(RoutingOffset)
    offset2.Offset = initOffset("parse", 0)
    offset2.Partition = 0
    parseOffsets[0] = offset2
    pendingUrls := make(chan []byte)
    go task.ParseLinks(pendingUrls, bloomFilter, taskConfig, kafkaConfig, &c2, offset2, MaxGoRoutine)

    go func() {
        for {
            url := <-pendingUrls
			if !bloomFilter.Lookup(url) {
				randomPartition := 0
				if MaxGoRoutine > 1 {
					randomPartition = rand.Intn(MaxGoRoutine - 1)
				}
				log.Debug("publish:",string(url),",partition:",randomPartition)
				publisher := kafka.NewBrokerPublisher(kafkaConfig.Hostname, taskConfig.Name+"_fetch", randomPartition)
				publisher.Publish(kafka.NewMessage(url))
				bloomFilter.Add(url)
			}else{
				log.Trace("hit bloomfilter,ignore,",string(url))
			}

        }
    }()

    <-finalQuitSignal
    log.Info("[gopa] is down")
}

func setLogging() {
	logPath:=taskConfig.BaseStoragePath+"log/filter.log";
    testConfig := `
	<seelog type="sync" minlevel="`
    testConfig = testConfig + logLevel
    testConfig = testConfig + `">
		<outputs formatid="main">
			<filter levels="error">
				<file path="`+logPath+`"/>
			</filter>
			<console />
		</outputs>
		<formats>
			<format id="main" format="[%LEV] %Msg%n"/>
		</formats>
	</seelog>`
    logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
    log.ReplaceLogger(logger)
}
