/**
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午4:45
 */
package main

import (
	 "flag"
	. "hunter"
	  log "github.com/cihub/seelog"
	"regexp"
	"os"

)

var seed_url = flag.String("seed", "", "Seed URL")
var siteConfig SiteConfig

func main() {
	defer log.Flush()
	setLogging();


	flag.Parse()

	log.Info("[gopa] is on.")

	if *seed_url == "" {
		log.Error("no seed was given.")
		os.Exit(1)
	}

	curl := make(chan []byte)
	success := make(chan Task)
	failure := make(chan string)

//	visited := make(map[string]int)

	// Setting siteConfig
	reg := regexp.MustCompile("<a.*?href=[\"'](http.*?)[\"']")



	MaxGoRouting:= 10

	siteConfig.LinkUrlExtractRegex = reg
	siteConfig.FollowSameDomain=true
	siteConfig.FollowSubDomain=true
	siteConfig.LinkUrlMustContain = "moko.cc"
//	siteConfig.LinkUrlMustNotContain = "wenku"

	// Giving a seed to gopa
	go Seed(curl, *seed_url)

	// Start the throttled crawling.
//	go ThrottledCrawl(curl, success, failure, visited)
	go ThrottledCrawl(curl,MaxGoRouting, success, failure)

	// Main loop that never exits and blocks on the data of a page.
	for {
		site := <-success
		go GetUrls(curl, site, siteConfig)
	}

}

func setLogging() {

	testConfig := `
	<seelog type="sync" minlevel="debug">
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

