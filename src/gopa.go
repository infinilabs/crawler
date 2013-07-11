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
	"util/bloom"
)

var seed_url = flag.String("seed", "", "Seed URL")
var siteConfig SiteConfig

func main() {
	log.Print("[gopa] is on.")

	// Normal bloom filter

	// Create a bloom filter which will contain an expected 100,000 items, and which
	// allows a false positive rate of 1%.
	f := bloom.New64(1000000, 0.01)

	// Add an item to the filter
	f.Add([]byte("foo"))

	// Check if an item has been added to the filter (if true, subject to the
	// false positive chance; if false, then the item definitely has not been
	// added to the filter.)
	log.Printf("%v", bool(f.Test([]byte("foo"))))

	flag.Parse()
	if *seed_url == "" {
		log.Fatal("no seed was given.")
	}

	curl := make(chan []byte)
	success := make(chan Task)
	failure := make(chan string)

	visited := make(map[string]int)

	reg := regexp.MustCompile("<a.*?href=[\"'](http.*?)[\"']")
	siteConfig.LinkUrlExtractRegex = reg

	// Give our crawler a place to start.
	go Seed(curl, *seed_url)

	// Start the throttled crawling.
	go ThrottledCrawl(curl, success, failure, visited)

	// Main loop that never exits and blocks on the data of a page.
	for {
		site := <-success
		go GetUrls(curl, site, siteConfig)
	}

}
