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
	"github.com/medcl/gopa/core/daemon"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/logger"
	"github.com/medcl/gopa/core/module"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/core/util"
	"github.com/medcl/gopa/modules"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

var (
	env             *Env
	startTime       time.Time
	finalQuitSignal chan bool
)

func onStart() {
	fmt.Println(GetWelcomeMessage())
	startTime = time.Now()
}

func onShutdown() {
	log.Debug(string(*stats.StatsAll()))
	fmt.Println("                         |    |                ")
	fmt.Println("   _` |   _ \\   _ \\   _` |     _ \\  |  |   -_) ")
	fmt.Println(" \\__, | \\___/ \\___/ \\__,_|   _.__/ \\_, | \\___| ")
	fmt.Println(" ____/                             ___/        ")
	fmt.Println("[gopa] "+VERSION+", uptime:", time.Now().Sub(startTime))
	fmt.Println(" ")
}

func main() {

	onStart()

	defer logger.Flush()

	var logLevel = flag.String("log", "info", "the log level,options:trace,debug,info,warn,error, default: info")
	var configFile = flag.String("config", "gopa.yml", "the location of config file, default: gopa.yml")
	var isDaemon = flag.Bool("daemon", false, "run in background as daemon")
	var pidfile = flag.String("pidfile", "", "pidfile path (only for daemon)")

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to this file")
	var memprofile = flag.String("memprofile", "", "write memory profile to this file")
	var startPprof = flag.Bool("pprof", false, "start pprof service, endpoint: http://localhost:6060/debug/pprof/")
	var isDebug = flag.Bool("debug", false, "enable debug")

	var httpBinding = flag.String("http_bind", "", "the http binding address, eg: 127.0.0.1:8001")
	var clusterBinding = flag.String("cluster_bind", "", "the cluster binding address, eg: 127.0.0.1:13001")
	var clusterSeed = flag.String("cluster_seeds", "", "the cluster address to start join in, seprated by comma, eg: 127.0.0.1:8001,127.0.0.1:8002,127.0.0.1:8003")
	var clusterName = flag.String("cluster_name", "gopa", "the cluster name, default: gopa")
	var dataDir = flag.String("data_path", "data", "the data path, default: data")
	var logDir = flag.String("log_path", "log", "the log path, default: log")

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

		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
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
		} else {
			fmt.Println("daemon only available in linux and darwin")
		}

	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.SetLogging(EmptyEnv(), *logLevel,*logDir)

	sysConfig := SystemConfig{ConfigFile: *configFile, LogLevel: *logLevel, HttpBinding: *httpBinding, ClusterBinding: *clusterBinding, ClusterSeeds: *clusterSeed, ClusterName: *clusterName, Data: *dataDir, Log: *logDir}
	sysConfig.Init()

	env = Environment(sysConfig)
	env.IsDebug = *isDebug
	//put env into global registrar
	global.RegisterEnv(env)
	logger.SetLogging(env,*logLevel,*logDir)

	//check instance lock
	util.CheckInstanceLock(env.SystemConfig.GetDataDir())

	module.New(env)
	modules.Register()
	module.Start()

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
			module.Stop()
			util.ClearInstanceLock()
			log.Flush()
			finalQuitSignal <- true
		}
	}()

	<-finalQuitSignal
	log.Debug("finished shutting down")

	onShutdown()
}
