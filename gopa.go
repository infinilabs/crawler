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
	"math/rand"
	"net/http"
	_ "net/http/pprof"
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
)

var seedUrl string
var logLevel string
var runtimeConfig *RuntimeConfig
var version string

func shutdown(offsets []*RoutingParameter, quitChannels []*chan bool, offsets2 []*RoutingParameter, quitChannels2 []*chan bool, quit chan bool) {
	log.Debug("start shutting down")
	for i := range quitChannels {
		*quitChannels[i] <- true
		log.Debug("send exit signal to quit channel-1,", i)
	}

	for i, item := range quitChannels2 {
		if item != nil {
			*item <- true
		}
		log.Debug("send exit signal to quit channel-2,", i)
	}

	log.Info("sent quit signal to go routings done")

	runtimeConfig.Storage.Close()
	log.Info("storage closed")

	quit <- true
	log.Debug("finished shutting down")
}

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

	runtimeConfig = InitOrGetConfig()

	runtimeConfig.Version = version

	runtimeConfig.LogLevel = logLevel

	logging.SetLogging(runtimeConfig.LogLevel, runtimeConfig.LogPath)

	if seedUrl == "" || seedUrl == "http://example.com" {
		log.Error("no seed was given. type:\"gopa -h\" for help.")
		os.Exit(1)
	}

	store := leveldb.LeveldbStore{}

	store.Open()
	runtimeConfig.Storage = &store

	//pprof server
	if runtimeConfig.GoProfEnabled {
		go func() {
			log.Info(http.ListenAndServe("localhost:6060", nil))
			log.Info("pprof server is up,http://localhost:6060/debug/pprof")
		}()
	}

	//start modules
	apiModule.Start(runtimeConfig)

	//adding default http protocol
	if !strings.HasPrefix(seedUrl, "http") {
		seedUrl = "http://" + seedUrl
	}

	maxGoRoutine := runtimeConfig.MaxGoRoutine
	fetchQuitChannels := make([]*chan bool, maxGoRoutine)   //shutdownSignal signals for each go routing
	fetchTaskChannels := make([]*chan []byte, maxGoRoutine) //fetchTask channels
	fetchOffsets := make([]*RoutingParameter, maxGoRoutine) //kafka fetchOffsets

	parseQuitChannels := make([]*chan bool, 2) //shutdownSignal signals for each go routing
	//	parseQuitChannels := make([]*chan bool, MaxGoRoutine) //shutdownSignal signals for each go routing
	parseOffsets := make([]*RoutingParameter, maxGoRoutine) //kafka fetchOffsets

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

			runtimeConfig.Storage.Close()

			//wait workers to exit
			log.Info("waiting workers exit")
			shutdown(fetchOffsets, fetchQuitChannels, parseOffsets, parseQuitChannels, shutdownSignal)
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
		parameter := new(RoutingParameter)
		parameter.Shard = i
		fetchOffsets[i] = parameter

		fetchTask := new(task.FetchTask)
		innerTaskConfig := new(task.InnerTaskConfig)
		innerTaskConfig.RuntimeConfig = runtimeConfig
		innerTaskConfig.MessageChan = &taskC
		innerTaskConfig.QuitChan = &quitC
		innerTaskConfig.Parameter = parameter

		fetchTask.Init(innerTaskConfig)
		go fetchTask.Start()

	}

	c2 := make(chan bool, 1)
	parseQuitChannels[0] = &c2
	offset2 := new(RoutingParameter)
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
	if runtimeConfig.ParseUrlsFromSavedFileLog {
		go task.ParseGo(pendingFetchUrls, runtimeConfig, &c2, offset2)
	}

	//redistribute pendingFetchUrls to sharded workers
	go func() {
		for {
			url := <-pendingFetchUrls
			if !runtimeConfig.Storage.UrlHasWalked(url) {

				if runtimeConfig.Storage.UrlHasFetched(url) {
					log.Warn("don't hit walk filter but hit fetch filter, also ignore,", string(url))
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
			} else {
				log.Trace("hit walk or fetch filter,just ignore,", string(url))
			}
		}
	}()

	//load predefined fetch jobs
	if runtimeConfig.LoadTemplatedFetchJob {
		go func() {

			if util.CheckFileExists(runtimeConfig.TaskConfig.TaskDataPath + "/urls/template.txt") {

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
	if runtimeConfig.LoadPendingFetchJobs {
		c3 := make(chan bool, 1)
		parseQuitChannels[1] = &c3
		offset3 := new(RoutingParameter)
		offset3.Shard = 0
		parseOffsets[1] = offset3
		go task.LoadTaskFromLocalFile(pendingFetchUrls, runtimeConfig, &c3, offset3)
	}

	//parse fetch failed jobs,and will ignore the walk-filter
	//TODO

	if runtimeConfig.LoadRuledFetchJob {
		log.Debug("start ruled fetch")
		go func() {
			if runtimeConfig.RuledFetchConfig.UrlTemplate != "" {
				for i := runtimeConfig.RuledFetchConfig.From; i <= runtimeConfig.RuledFetchConfig.To; i += runtimeConfig.RuledFetchConfig.Step {
					url := strings.Replace(runtimeConfig.RuledFetchConfig.UrlTemplate, "{id}", strconv.FormatInt(int64(i), 10), -1)
					log.Debug("add ruled url:", url)
					pendingFetchUrls <- []byte(url)
				}
			} else {
				log.Error("ruled template is empty,ignore")
			}
		}()

	}

	<-finalQuitSignal
	printShutdownInfo()
}
