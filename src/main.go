/**
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午4:45
 */
package main

import (

	"log"
	"regexp"
    . "hunter"
)

func main() {
    log.Print("gopa is on.")

    curl := make(chan []byte)
    success := make(chan Task)
    failure := make(chan string)

    visited := make(map[string]int)

    regex := regexp.MustCompile("<a.*?href=[\"'](http.*?)[\"']")

    // Give our crawler a place to start.
    go Seed(curl,"http://www.baidu.com")

    // Start the throttled crawling.
    go ThrottledCrawl(curl, success, failure, visited)

    // Main loop that never exits and blocks on the data of a page.
    for {
        site := <-success
        go GetUrls(curl, site, regex)
    }

}


