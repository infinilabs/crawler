/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	logging "github.com/medcl/gopa/core/logging"
	"github.com/medcl/gopa/core/store/leveldb"
	task "github.com/medcl/gopa/core/tasks"
	"github.com/medcl/gopa/core/util"
	apiModule "github.com/medcl/gopa/modules/api"
	crawlerModule "github.com/medcl/gopa/modules/crawler"
	profilerModule "github.com/medcl/gopa/modules/profiler"
)

var seedUrl string
var logLevel string
var gopaConfig *GopaConfig
var version string

var startTime time.Time

func printStartInfo() {
	fmt.Println("  __ _  ___  _ __   __ _ ")
	fmt.Println(" / _` |/ _ \\| '_ \\ / _` |")
	fmt.Println("| (_| | (_) | |_) | (_| |")
	fmt.Println(" \\__, |\\___/| .__/ \\__,_|")
	fmt.Println(" |___/      |_|          ")
	fmt.Println(" ")

	fmt.Println("[gopa] " + version + " is on")
	fmt.Println(" ")
	startTime = time.Now()
}

func printShutdownInfo() {
	fmt.Println("                         |    |                ")
	fmt.Println("   _` |   _ \\   _ \\   _` |     _ \\  |  |   -_) ")
	fmt.Println(" \\__, | \\___/ \\___/ \\__,_|   _.__/ \\_, | \\___| ")
	fmt.Println(" ____/                             ___/        ")
	fmt.Println("[gopa] "+version+" is down, uptime:", time.Now().Sub(startTime))
	fmt.Println(" ")
}

func main() {

	version = "0.7_SNAPSHOT"

	printStartInfo()
	defer logging.Flush()

	flag.StringVar(&seedUrl, "seed", "http://example.com", "the seed url,where everything starts")
	flag.StringVar(&logLevel, "log", "info", "setting log level,options:trace,debug,info,warn,error")

	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	logging.SetInitLogging(logLevel)

	gopaConfig = InitGopaConfig()

	gopaConfig.RuntimeConfig = InitOrGetConfig()

	gopaConfig.SystemConfig.Version = version

	gopaConfig.RuntimeConfig.LogLevel = logLevel

	logging.SetLogging(gopaConfig.RuntimeConfig.LogLevel, gopaConfig.RuntimeConfig.LogPath)

	store := leveldb.LeveldbStore{}

	store.Open()
	gopaConfig.RuntimeConfig.Storage = &store

	//start modules
	apiModule.Start(gopaConfig)
	profilerModule.Start(gopaConfig)
	crawlerModule.Start(gopaConfig)

	//adding default http protocol
	if !strings.HasPrefix(seedUrl, "http") {
		seedUrl = "http://" + seedUrl
	}

	parseQuitChannels := make([]*chan bool, 2) //shutdownSignal signals for each go routing
	parseOffsets := make([]*RoutingParameter, gopaConfig.RuntimeConfig.MaxGoRoutine)

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

			//wait workers to exit
			log.Info("waiting workers exit")
			shutdown(parseOffsets, parseQuitChannels, shutdownSignal)
			apiModule.Stop()
			crawlerModule.Stop()
			apiModule.Stop()
			<-shutdownSignal
			log.Info("workers shutdown")
			finalQuitSignal <- true
		}
	}()

	c2 := make(chan bool, 1)
	parseQuitChannels[0] = &c2
	offset2 := new(RoutingParameter)
	offset2.Shard = 0
	parseOffsets[0] = offset2

	//fetch rule:all urls -> persisted to sotre -> fetched from store -> pushed to pendingFetchUrls -> redistributed to sharded goroutines -> fetch -> save webpage to store -> done
	//parse rule:url saved to store -> local path persisted to store -> fetched to pendingParseFiles -> redistributed to sharded goroutines -> parse -> clean urls -> enqueue to url store ->done

	//sending feed to task queue
	go func() {
		//notice seed will not been persisted
		log.Debug("sending feed to fetch queue,", seedUrl)
		gopaConfig.Channels.PendingFetchUrl <- []byte(seedUrl)
	}()

	//start local saved file parser
	if gopaConfig.RuntimeConfig.ParseUrlsFromSavedFileLog {
		go task.ParseGo(gopaConfig.Channels.PendingFetchUrl, gopaConfig.RuntimeConfig, &c2, offset2)
	}

	//load predefined fetch jobs
	if gopaConfig.RuntimeConfig.LoadTemplatedFetchJob {
		go func() {

			if util.CheckFileExists(gopaConfig.RuntimeConfig.TaskConfig.TaskDataPath + "/urls/template.txt") {

				templates := util.ReadAllLines(gopaConfig.RuntimeConfig.TaskConfig.TaskDataPath + "/urls/template.txt")
				ids := util.ReadAllLines(gopaConfig.RuntimeConfig.TaskConfig.TaskDataPath + "/urls/id.txt")

				for _, id := range ids {
					for _, template := range templates {
						log.Trace("id:", id)
						log.Trace("template:", template)
						url := strings.Replace(template, "{id}", id, -1)
						log.Debug("new task from template:", url)
						gopaConfig.Channels.PendingFetchUrl <- []byte(url)
					}
				}
				log.Info("templated download is done.")

			}

		}()
	}

	//fetch urls from saved pages
	if gopaConfig.RuntimeConfig.LoadPendingFetchJobs {
		c3 := make(chan bool, 1)
		parseQuitChannels[1] = &c3
		offset3 := new(RoutingParameter)
		offset3.Shard = 0
		parseOffsets[1] = offset3
		go task.LoadTaskFromLocalFile(gopaConfig.Channels.PendingFetchUrl, gopaConfig.RuntimeConfig, &c3, offset3)
	}

	//parse fetch failed jobs,and will ignore the walk-filter
	//TODO

	if gopaConfig.RuntimeConfig.LoadRuledFetchJob {
		log.Debug("start ruled fetch")
		go func() {
			if gopaConfig.RuntimeConfig.RuledFetchConfig.UrlTemplate != "" {
				for i := gopaConfig.RuntimeConfig.RuledFetchConfig.From; i <= gopaConfig.RuntimeConfig.RuledFetchConfig.To; i += gopaConfig.RuntimeConfig.RuledFetchConfig.Step {
					url := strings.Replace(gopaConfig.RuntimeConfig.RuledFetchConfig.UrlTemplate, "{id}", strconv.FormatInt(int64(i), 10), -1)
					log.Debug("add ruled url:", url)
					gopaConfig.Channels.PendingFetchUrl <- []byte(url)
				}
			} else {
				log.Error("ruled template is empty,ignore")
			}
		}()

	}

	<-finalQuitSignal
	printShutdownInfo()
}

func shutdown(offsets2 []*RoutingParameter, quitChannels2 []*chan bool, quit chan bool) {
	log.Debug("start shutting down")

	for i, item := range quitChannels2 {
		if item != nil {
			*item <- true
		}
		log.Debug("send exit signal to quit channel-2,", i)
	}

	log.Info("sent quit signal to go routings done")

	gopaConfig.RuntimeConfig.Storage.Close()
	log.Info("storage closed")

	quit <- true
	log.Debug("finished shutting down")
}
