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
	"expvar"
	_ "expvar"
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/daemon"
	"github.com/infinitbyte/gopa/core/env"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/logger"
	"github.com/infinitbyte/gopa/core/module"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules"
	"github.com/infinitbyte/gopa/plugins"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

var (
	environment     *env.Env
	startTime       time.Time
	finalQuitSignal chan bool
)

func onStart() {
	fmt.Println(env.GetWelcomeMessage())
	startTime = time.Now()
}

func onShutdown(isDaemon bool) {
	if environment.IsDebug {
		fmt.Println(string(*stats.StatsAll()))
	}

	//force flush all logs
	log.Flush()

	if isDaemon {
		fmt.Println("[gopa] started.")
		return
	}
	fmt.Println("                         |    |                ")
	fmt.Println("   _` |   _ \\   _ \\   _` |     _ \\  |  |   -_) ")
	fmt.Println(" \\__, | \\___/ \\___/ \\__,_|   _.__/ \\_, | \\___| ")
	fmt.Println(" ____/                             ___/        ")
	fmt.Println("[gopa] "+env.VERSION+", uptime:", time.Since(startTime))
	fmt.Println(" ")
}

// report expvar and all metrics
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	first := true
	report := func(key string, value interface{}) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		if str, ok := value.(string); ok {
			fmt.Fprintf(w, "%q: %q", key, str)
		} else {
			fmt.Fprintf(w, "%q: %v", key, value)
		}
	}

	fmt.Fprintf(w, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		report(kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	onStart()

	var logLevel = flag.String("log", "info", "the log level,options:trace,debug,info,warn,error")
	var configFile = flag.String("config", "gopa.yml", "the location of config file")
	var isDaemon = flag.Bool("daemon", false, "run in background as daemon")
	var pidfile = flag.String("pidfile", "", "pidfile path (only for daemon)")

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to this file")
	var memprofile = flag.String("memprofile", "", "write memory profile to this file")
	var httpprof = flag.String("pprof", "", "enable and setup pprof/expvar service, eg: localhost:6060 , the endpoint will be: http://localhost:6060/debug/pprof/ and http://localhost:6060/debug/vars")
	var isDebug = flag.Bool("debug", false, "run in debug mode, wi")

	var logDir = flag.String("log_path", "log", "the log path")

	flag.Parse()
	logger.SetLogging(env.EmptyEnv(), *logLevel, *logDir)

	environment = env.Environment(*configFile)
	environment.IsDebug = *isDebug
	//put env into global registrar
	global.RegisterEnv(environment)

	logger.SetLogging(environment, *logLevel, *logDir)

	//check instance lock
	util.CheckInstanceLock(environment.SystemConfig.GetWorkingDir())

	//set path to persist id
	util.RestorePersistID(environment.SystemConfig.GetWorkingDir())

	//cleanup
	defer func() {
		util.ClearInstanceLock()

		if r := recover(); r != nil {
			if r == nil {
				return
			}
			var v string
			switch r.(type) {
			case error:
				v = r.(error).Error()
			case runtime.Error:
				v = r.(runtime.Error).Error()
			case string:
				v = r.(string)
			}
			log.Error("main: ", v)
		}
		util.SnapshotPersistID()
		log.Flush()
		logger.Flush()

		//print goodbye message
		onShutdown(*isDaemon)
	}()

	//profile options
	if *httpprof != "" {
		go func() {
			log.Infof("pprof listen at: http://%s/debug/pprof/", *httpprof)
			mux := http.NewServeMux()

			// register pprof handler
			mux.HandleFunc("/debug/pprof/", func(w http.ResponseWriter, r *http.Request) {
				http.DefaultServeMux.ServeHTTP(w, r)
			})

			// register metrics handler
			mux.HandleFunc("/debug/vars", metricsHandler)

			endpoint := http.ListenAndServe(*httpprof, mux)
			log.Debug("stop pprof server: %v", endpoint)
		}()
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		if *memprofile != "" {
			f, err := os.Create(*memprofile)
			if err != nil {
				panic(err)
			}
			pprof.WriteHeapProfile(f)
			f.Close()
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

	//modules
	module.New()
	modules.Register()
	plugins.Register()
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
		if s == os.Interrupt || s.(os.Signal) == syscall.SIGINT || s.(os.Signal) == syscall.SIGTERM ||
			s.(os.Signal) == syscall.SIGKILL || s.(os.Signal) == syscall.SIGQUIT {
			fmt.Printf("\n[gopa] got signal:%s, start shutting down\n", s.String())
			//wait workers to exit
			module.Stop()
			finalQuitSignal <- true
		}
	}()

	<-finalQuitSignal

}
