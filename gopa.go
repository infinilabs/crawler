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
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/logging"
	"github.com/medcl/gopa/core/store/leveldb"
	task "github.com/medcl/gopa/core/tasks"
	"github.com/medcl/gopa/core/util"
	apiModule "github.com/medcl/gopa/modules/api"
	crawlerModule "github.com/medcl/gopa/modules/crawler"
	parserModule "github.com/medcl/gopa/modules/parser"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var seedUrl string
var logLevel string
var env *Env
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

	logging.SetInitLogging(NullEnv(), logLevel)

	env = Environment(&Registrar{}, &SystemConfig{Version: version}, InitOrGetConfig())

	env.RuntimeConfig.LogLevel = logLevel

	logging.SetLogging(env, env.RuntimeConfig.LogLevel, env.RuntimeConfig.LogPath)

	store := leveldb.LeveldbStore{}

	store.Open()
	env.RuntimeConfig.Storage = &store

	//start modules
	apiModule.Start(env)
	crawlerModule.Start(env)
	parserModule.Start(env)

	//adding default http protocol
	if !strings.HasPrefix(seedUrl, "http") {
		seedUrl = "http://" + seedUrl
	}

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
			log.Warn("got signal:os.Interrupt,start shutting down")

			close(env.Channels.PendingFetchUrl)
			//wait workers to exit
			shutdownSignal <- true
			parserModule.Stop()
			crawlerModule.Stop()
			apiModule.Stop()
			env.RuntimeConfig.Storage.Close()
			<-shutdownSignal
			log.Info("all modules stopeed")
			finalQuitSignal <- true
		}
	}()

	//fetch rule:all urls -> persisted to sotre -> fetched from store -> pushed to pendingFetchUrls -> redistributed to sharded goroutines -> fetch -> save webpage to store -> done
	//parse rule:url saved to store -> local path persisted to store -> fetched to pendingParseFiles -> redistributed to sharded goroutines -> parse -> clean urls -> enqueue to url store ->done

	//sending feed to task queue
	go func() {
		//notice seed will not been persisted
		log.Debug("sending feed to fetch queue,", seedUrl)
		env.Channels.PendingFetchUrl <- []byte(seedUrl)
	}()

	//load predefined fetch jobs
	if env.RuntimeConfig.LoadTemplatedFetchJob {
		go func() {

			if util.CheckFileExists(env.RuntimeConfig.TaskConfig.TaskDataPath + "/urls/template.txt") {

				templates := util.ReadAllLines(env.RuntimeConfig.TaskConfig.TaskDataPath + "/urls/template.txt")
				ids := util.ReadAllLines(env.RuntimeConfig.TaskConfig.TaskDataPath + "/urls/id.txt")

				for _, id := range ids {
					for _, template := range templates {
						log.Trace("id:", id)
						log.Trace("template:", template)
						url := strings.Replace(template, "{id}", id, -1)
						log.Debug("new task from template:", url)
						env.Channels.PendingFetchUrl <- []byte(url)
					}
				}
				log.Info("templated download is done.")

			}

		}()
	}

	//fetch urls from saved pages
	if env.RuntimeConfig.LoadPendingFetchJobs {
		go task.LoadTaskFromLocalFile(env.Channels.PendingFetchUrl, env.RuntimeConfig)
	}

	//parse fetch failed jobs,and will ignore the walk-filter
	//TODO

	if env.RuntimeConfig.LoadRuledFetchJob {
		log.Debug("start ruled fetch")
		go func() {
			if env.RuntimeConfig.RuledFetchConfig.UrlTemplate != "" {
				for i := env.RuntimeConfig.RuledFetchConfig.From; i <= env.RuntimeConfig.RuledFetchConfig.To; i += env.RuntimeConfig.RuledFetchConfig.Step {
					url := strings.Replace(env.RuntimeConfig.RuledFetchConfig.UrlTemplate, "{id}", strconv.FormatInt(int64(i), 10), -1)
					log.Debug("add ruled url:", url)
					env.Channels.PendingFetchUrl <- []byte(url)
				}
			} else {
				log.Error("ruled template is empty,ignore")
			}
		}()

	}

	<-finalQuitSignal

	log.Debug("finished shutting down")

	printShutdownInfo()
}
