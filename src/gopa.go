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
	"regexp"
	. "webhunter"
	. "github.com/zeebo/sbloom"
	"hash/fnv"
	"io/ioutil"
	"util"
	"os/signal"
	"strings"
)

var seedUrl = flag.String("seed", "http://example.com", "the seed url,where everything begins")
var logLevel = flag.String("log", "info", "setting log level,options:trace,debug,info,warn,error")
var siteConfig SiteConfig
var bloomFilter *Filter

func persistBloomFilter(bloomFilterPersistFileName string){
	//save bloom-filter
	m, err := bloomFilter.GobEncode()
	if err != nil {
		log.Error(err)
	}
	err = ioutil.WriteFile(bloomFilterPersistFileName, m, 0600)
	if err != nil {
		panic(err)
	}
	log.Info("bloomFilter safety persisted.")
}

func main() {
	defer log.Flush()

	flag.Parse()

	setLogging()


	log.Info("[gopa] is on.")

	if *seedUrl == "" || *seedUrl =="http://example.com" {
		log.Error("no seed was given. type:\"gopa -h\" for help.")
		os.Exit(1)
	}

	curl := make(chan []byte)
	success := make(chan Task)
	failure := make(chan string)

	// Setting siteConfig
	reg := regexp.MustCompile("(src2|src|href|HREF|SRC)\\s*=\\s*[\"']?(.*?)[\"']")

	MaxGoRouting := 1

	//loading or initializing bloom filter
	bloomFilterPersistFileName:="bloomfilter.bin"
	if util.CheckFileExists(bloomFilterPersistFileName){
		log.Debug("found bloomFilter,start reload")
		n,err := ioutil.ReadFile(bloomFilterPersistFileName)
		if err != nil {
			log.Error(err)
		}
		bloomFilter= new(Filter)
		if err := bloomFilter.GobDecode(n); err != nil {
			log.Error(err)
		}

		log.Info("bloomFilter successfully reloaded")
	}else{
		probItems:=1000000
		log.Debug("initializing bloom-filter,virual size is,",probItems)
		bloomFilter = NewFilter(fnv.New64(), probItems)
		log.Info("bloomFilter successfully initialized")
	}


	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func(){
//		for sig := range c {
			// sig is a ^C, handle it
			s := <-c
			log.Debug("got signal:", s)
		    if(s == os.Interrupt){
				log.Warn("got signal:os.Interrupt,saving data and exit")
				defer os.Exit(0)
				persistBloomFilter(bloomFilterPersistFileName)
				log.Info("[gopa] is down")

			}
//			persistBloomFilter(bloomFilterPersistFileName)
//		}
	}()


	siteConfig.LinkUrlExtractRegex = reg
	siteConfig.FollowSameDomain = true
	siteConfig.FollowSubDomain = true
	siteConfig.LinkUrlMustContain = "moko.cc"
	//	siteConfig.LinkUrlMustNotContain = "wenku"

	if !strings.HasPrefix(*seedUrl,"http"){
		*seedUrl="http://"+*seedUrl
	}

	// Giving a seed to gopa
	go Seed(curl, *seedUrl)

	// Start the throttled crawling.
	go ThrottledCrawl(bloomFilter,curl, MaxGoRouting, success, failure)





	// Main loop that never exits and blocks on the data of a page.
	for {
		site := <-success
		go GetUrls(bloomFilter,curl, site, siteConfig)
	}

}



func setLogging() {
	testConfig := `
	<seelog type="sync" minlevel="`
	testConfig =testConfig + *logLevel
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
