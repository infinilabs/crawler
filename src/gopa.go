/**
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午4:45
 */
package main

import (
	"flag"
	. "hunter"
	"log"
	"regexp"
)

var seed_url = flag.String("seed", "", "Seed URL")
var siteConfig SiteConfig

func main() {
	log.Print("[gopa] is on.")

	flag.Parse()
	if *seed_url == "" {
		log.Fatal("no seed was given.")
	}

	curl := make(chan []byte)
	success := make(chan Task)
	failure := make(chan string)

	visited := make(map[string]int)

	// Setting siteConfig
	reg := regexp.MustCompile("<a.*?href=[\"'](http.*?)[\"']")
	siteConfig.LinkUrlExtractRegex = reg
	siteConfig.LinkUrlMustContain = "baidu"
	siteConfig.LinkUrlMustNotContain = "wenku"

	// Giving a seed to gopa
	go Seed(curl, *seed_url)

	// Start the throttled crawling.
	go ThrottledCrawl(curl, success, failure, visited)

	// Main loop that never exits and blocks on the data of a page.
	for {
		site := <-success
		go GetUrls(curl, site, siteConfig)
	}

}
