/**
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午4:45
 */
package main

import (
	"flag"
	log "github.com/cihub/seelog"
	"os"
	. "webhunter"
	. "github.com/zeebo/sbloom"
	"hash/fnv"
	"io/ioutil"
	"util"
	"os/signal"
	"strings"
 _ "net/http/pprof"
	"net/http"
	"regexp"
	config "config"
	"fmt"
	"kafka"
	"syscall"
//	"time"
	"strconv"
)

var seedUrl string
var logLevel string
var siteConfig *TaskConfig
var kafkaConfig *config.KafkaConfig
var bloomFilter *Filter
var maxGoRouting int

func persistBloomFilter(bloomFilterPersistFileName string){
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

func getSeqStr(start []byte,end []byte,mix bool) []byte{
	if(len(start)) == len(end){
		for i:=range start {
			fmt.Println(start[i])
		}
//		if(start>64 && end < 123){

//		}
	}

	return nil
}

func init(){
}


func initOffset( partition int) uint64{
	log.Info("start init offsets,",partition)

	path:="offset_"+strconv.FormatInt(int64(partition),10)
	if util.CheckFileExists(path){
		log.Debug("found offset file,start loading")
		n,err := ioutil.ReadFile(path)
		if err != nil {
			log.Error("offset",err)
			return  0
		}
		ret,err:=strconv.ParseInt(string(n),10,64)
		if err != nil {
			log.Error("offset",err)
			return  0
		}
		log.Info("init offsets successfully,",partition,":",ret)
		return uint64(ret)
	}

	return 0
}


func closeKafkaConsumer(offsets []*RoutingOffset,quitChannels []*chan bool,quit chan bool){

	for i:=range quitChannels{
		log.Debug("send exit signal to channel,",i)
		*quitChannels[i] <- true
	}
	log.Info("sent quit signal to go routings done")

	for i:=range offsets{
		//TODO
		log.Info("persist offset,",i,":",offsets[i].Offset,",",offsets[i].Partition)
	}

	log.Info("persist kafka offsets done")

	quit <- true
}


func main() {


	flag.StringVar(&seedUrl,"seed", "http://example.com", "the seed url,where everything begins")
	flag.StringVar(&logLevel,"log", "info", "setting log level,options:trace,debug,info,warn,error")

	flag.Parse()

	defer log.Flush()
	setLogging()

	log.Info("[gopa] is on.")


	if seedUrl == "" || seedUrl =="http://example.com" {
		log.Error("no seed was given. type:\"gopa -h\" for help.")
		os.Exit(1)
	}

	//urls need to be fetch
	curl := make(chan []byte)
	//tasks fetched,need to be parse
	success := make(chan Task)
	//urls failure
	failure := make(chan string)

	// Setting siteConfig
	maxGoRouting = config.GetIntConfig("Global", "maxGoRouting",1)

	if(maxGoRouting<0){
		maxGoRouting=1
	}

	//loading or initializing bloom filter
	bloomFilterPersistFileName:=config.GetStringConfig("BloomFilter", "FileName","bloomfilter.bin")

	if util.CheckFileExists(bloomFilterPersistFileName){
		log.Debug("found bloomFilter,start reload")
		n,err := ioutil.ReadFile(bloomFilterPersistFileName)
		if err != nil {
			log.Error("bloomFilter",err)
			return
		}
		bloomFilter= new(Filter)
		if err := bloomFilter.GobDecode(n); err != nil {
			log.Error("bloomFilter",err)
			return
		}

		log.Info("bloomFilter successfully reloaded")
	}else{
		probItems:=config.GetIntConfig("BloomFilter", "ItemSize",100000)
		log.Debug("initializing bloom-filter,virual size is,",probItems)
		bloomFilter = NewFilter(fnv.New64(), probItems)
		log.Info("bloomFilter successfully initialized")
	}




//	atr:="AZaz"
//	btr:=[]byte(atr)
//	fmt.Println(btr)
//
//	id:= getSeqStr([]byte("AA"),[]byte("ZZ"),false)
//	fmt.Println(id)

	//pprof serves
	go func() {
		log.Info(http.ListenAndServe("localhost:6060", nil))
		log.Info("pprof server is up,http://localhost:6060/debug/pprof")
	}()




	//setting siteConfig
	siteConfig=new (TaskConfig)
	siteConfig.LinkUrlExtractRegex = regexp.MustCompile(
	config.GetStringConfig("CrawlerRule","LinkUrlExtractRegex","(src2|src|href|HREF|SRC)\\s*=\\s*[\"']?(.*?)[\"']"))

	siteConfig.Name = config.GetStringConfig("CrawlerRule","Name","GopaTask")

	siteConfig.FollowSameDomain = config.GetBoolConfig("CrawlerRule","FollowSameDomain",true)
	siteConfig.FollowSubDomain  =  config.GetBoolConfig("CrawlerRule","FollowSubDomain",true)
	siteConfig.LinkUrlMustContain =config.GetStringConfig("CrawlerRule","LinkUrlMustContain","")
	siteConfig.LinkUrlMustNotContain = config.GetStringConfig("CrawlerRule","LinkUrlMustNotContain","")

	siteConfig.SkipPageParsePattern = regexp.MustCompile(config.GetStringConfig("CrawlerRule","SkipPageParsePattern",".*?\\.((js)|(css)|(rar)|(gz)|(zip)|(exe)|(bmp)|(jpeg)|(gif)|(png)|(jpg)|(apk))\\b"))//end with js,css,apk,zip,ignore

	siteConfig.DownloadUrlPattern= regexp.MustCompile(config.GetStringConfig("CrawlerRule","DownloadUrlPattern",".*"))
	siteConfig.DownloadUrlMustContain=config.GetStringConfig("CrawlerRule","DownloadUrlMustContain","")
	siteConfig.DownloadUrlMustNotContain=config.GetStringConfig("CrawlerRule","DownloadUrlMustNotContain","")


	kafkaConfig=new (config.KafkaConfig)
	kafkaConfig.Hostname =config.GetStringConfig("Kafka","Hostname","localhost:9092")
	kafkaConfig.MaxSize =uint32(config.GetIntConfig("Kafka","MaxSize",1048576))


	broker := kafka.NewBrokerPublisher(kafkaConfig.Hostname, siteConfig.Name, 0)
	log.Info("kafka publisher is up, connect at: ",kafkaConfig.Hostname,"")


	//adding default http protocol
	if !strings.HasPrefix(seedUrl,"http"){
		seedUrl="http://"+seedUrl
	}

	log.Debug("init KafkaChannel signal channel,size:",maxGoRouting)
	quitChannels:=make([]*chan bool,maxGoRouting)  //quit signals for each go routing
	offsets:=make([]*RoutingOffset,maxGoRouting) //kafka offsets

	for i := 0; i < maxGoRouting; i++ {
		c:=make(chan bool, 1)
		quitChannels[i]=&c
		offset:=new (RoutingOffset)
		offset.Offset=initOffset(i)

		offset.Partition=i
		offsets[i]=offset
	}


	quit := make(chan bool, 1)

	//handle exit event
	exitEventChannel := make(chan os.Signal, 1)
	signal.Notify(exitEventChannel, syscall.SIGINT)
	signal.Notify(exitEventChannel, os.Interrupt)
	go func(){
		s := <-exitEventChannel
		log.Debug("got signal:", s)
		if(s == os.Interrupt || s.(os.Signal) == syscall.SIGINT){
			log.Warn("got signal:os.Interrupt,saving data and exit")
			defer os.Exit(0)

			persistBloomFilter(bloomFilterPersistFileName)

			//wait kafka to exit
			log.Info("waiting kafka exit")
			closeKafkaConsumer(offsets,quitChannels,quit)

			<-quit
			log.Info("[gopa] is down")
		}
	}()


	// Start the crawling gorouting,listening kafka.
	go ThrottledCrawl(bloomFilter,siteConfig,kafkaConfig,curl, maxGoRouting, success, failure,quitChannels,offsets)


	// Giving a seed to gopa
//	go Seed(curl, *seedUrl)

	broker.Publish(kafka.NewMessage([]byte(seedUrl)))

	// Main loop that never exits and blocks on the data of a page.
	for {
		taskItem := <-success     //TODO 处理Kafka关闭异常
		ExtractLinksFromTaskResponse(bloomFilter,broker, taskItem, siteConfig)
	}

}



func setLogging() {
	testConfig := `
	<seelog type="sync" minlevel="`
	testConfig =testConfig + logLevel
	testConfig =testConfig +`">
		<outputs formatid="main">
			<filter levels="error">
				<file path="./log/filter.log"/>
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
