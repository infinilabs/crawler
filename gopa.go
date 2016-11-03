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
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/logging"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/core/types"
	modules "github.com/medcl/gopa/modules"
	"github.com/medcl/gopa/core/daemon"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
	"github.com/medcl/gopa/core/global"
)

var (
	env             *Env
	startTime       time.Time
	components      *modules.Modules
	finalQuitSignal chan bool
)

func onStart() {
	fmt.Println(GetWelcomeMessage())
	startTime = time.Now()
}

func onShutdown() {
	log.Debug(string(stats.StatsAll()))
	fmt.Println("                         |    |                ")
	fmt.Println("   _` |   _ \\   _ \\   _` |     _ \\  |  |   -_) ")
	fmt.Println(" \\__, | \\___/ \\___/ \\__,_|   _.__/ \\_, | \\___| ")
	fmt.Println(" ____/                             ___/        ")
	fmt.Println("[gopa] "+VERSION+", uptime:", time.Now().Sub(startTime))
	fmt.Println(" ")
}

func main() {

	onStart()

	defer logging.Flush()

	var seedUrl = flag.String("seed", "", "the seed url, where everything starts")
	var logLevel = flag.String("log", "info", "the log level,options:trace,debug,info,warn,error, default: info")
	var configFile = flag.String("config", "gopa.yml", "the location of config file, default: gopa.yml")
	var isDaemon = flag.Bool("daemon", false, "run in background as daemon")
	var pidfile = flag.String("pidfile", "", "pidfile path (only for daemon)")

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to this file")
	var memprofile = flag.String("memprofile", "", "write memory profile to this file")
	var startPprof = flag.Bool("pprof", false, "start pprof service, endpoint: http://localhost:6060/debug/pprof/")

	flag.Parse()

	if *startPprof {
		go func() {
			http.ListenAndServe("localhost:6060", nil)
		}()
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Error(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		if *memprofile != "" {
			f, err := os.Create(*memprofile)
			if err != nil {
				log.Error(err)
			}
			pprof.WriteHeapProfile(f)
			defer f.Close()
		}
	}

	//daemon
	if *isDaemon {

		if(runtime.GOOS=="darwin"||runtime.GOOS=="linux"){
			runtime.LockOSThread()
			context := new(daemon.Context)
			if *pidfile != "" {
				context.PidFileName = *pidfile
				context.PidFilePerm = 0644
			}

			child, _ := context.Reborn()

			if child != nil {
				return
			}
			defer context.Release()

			runtime.UnlockOSThread()
		}else{
			fmt.Println("daemon only available in linux and darwin")
		}

	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	logging.SetInitLogging(EmptyEnv(), *logLevel)

	sysConfig := SystemConfig{Version: VERSION, ConfigFile: *configFile, LogLevel: *logLevel}

	env = Environment(sysConfig)

	//put env into global registrar
	global.Register(global.REGISTER_ENV,&env)

	logging.SetLogging(env)

	components = modules.New(env)
	components.Start()

	finalQuitSignal = make(chan bool)

	//handle exit event
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)

	go func() {
		s := <-sigc
		log.Info("got signal:", s)
		if s == os.Interrupt || s.(os.Signal) == syscall.SIGINT || s.(os.Signal) == syscall.SIGTERM ||
			s.(os.Signal) == syscall.SIGKILL || s.(os.Signal) == syscall.SIGQUIT {
			log.Infof("got signal:%s ,start shutting down", s.String())
			//wait workers to exit
			components.Stop()
			env.Channels.Close()
			finalQuitSignal <- true
		}
	}()

	//sending feed to task queue
	if len(*seedUrl) > 0 {
		log.Debug("sending feed to fetch queue,", *seedUrl)
		env.Channels.PushUrlToCheck(types.NewPageTask(*seedUrl, "", 0))
	}

	<-finalQuitSignal

	log.Debug("finished shutting down")

	onShutdown()
}
