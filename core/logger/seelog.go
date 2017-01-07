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

package logger

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/config"
	"strings"
	"fmt"
	"sync"
)

var file string
var loggingConfig *config.LoggingConfig
var l sync.Mutex
var e *Env
func SetLogging(env *Env, logLevel string, logFile string) {

	e=env

	l.Lock()
	  if(loggingConfig==nil){
		  loggingConfig=&config.LoggingConfig{}
		  loggingConfig.LogLevel="info"
		  loggingConfig.FileFilterPattern="*"
		  loggingConfig.FuncFilterPattern="*"
	  }
	l.Unlock()

	if(env!=nil){
		envLevel := strings.ToLower(env.LoggingLevel)
		if(env.SystemConfig!=nil){
			envLogFile := env.SystemConfig.Log + "/gopa.log"
			if(len(envLogFile)>0){
				file=envLogFile
			}
		}
		if(len(envLevel)>0){
			loggingConfig.LogLevel=envLevel
		}
	}

	//overwrite env config
	if len(logLevel) > 0 {
		loggingConfig.LogLevel = strings.ToLower(logLevel)
	}

	if len(logFile) > 0 {
		file = logFile
	}

	//finally check filename
	if len(file) <= 0 {
		file = "./log/gopa.log"
	}

	consoleWriter, _ := NewConsoleWriter()
	websocketWriter, _ := NewWebsocketWriter()

	formatter, _ := log.NewFormatter("[%Date(01-02) %Time] [%LEV] [%File:%Line] %Msg%n")

	rollingFileWriter,_:=NewRollingFileWriterSize(file, rollingArchiveNone, "", 10000000000, 5, rollingNameModePostfix)
	bufferedWriter,_:=NewBufferedWriter(rollingFileWriter,10000,1000)

	l,_:=log.LogLevelFromString(loggingConfig.LogLevel)

	//logging receivers
	receivers:=[]interface{}{consoleWriter,bufferedWriter}
	if(loggingConfig.RealtimePushEnabled){
		receivers=append(receivers,websocketWriter)
	}

	root, _ := log.NewFilterDispatcher(formatter,receivers ,l)

	constraints, _ := log.NewMinMaxConstraints(l, log.CriticalLvl)

	specificConstraints, _ := log.NewListConstraints([]log.LogLevel{l, log.CriticalLvl})

	ex, _ := log.NewLogLevelException(loggingConfig.FuncFilterPattern, loggingConfig.FileFilterPattern, specificConstraints)

	exceptions := []*log.LogLevelException{ex}

	logger := log.NewAsyncLoopLogger(log.NewLoggerConfig(constraints, exceptions, root))

	err:=log.ReplaceLogger(logger)
	if(err!=nil){
		fmt.Println(err)
	}
}

func GetLoggingConfig() *config.LoggingConfig {
	return loggingConfig
}

func UpdateLoggingConfig(config *config.LoggingConfig)  {
	l.Lock()
	loggingConfig=config
	l.Unlock()
	SetLogging(e,"","")
}

func Flush() {
	log.Flush()
}

var websocketHandler func(message string, level log.LogLevel, context log.LogContextInterface)

func RegisterWebsocketHandler(func1 func(message string, level log.LogLevel, context log.LogContextInterface)) {

	websocketHandler = func1
	if func1 != nil {
		log.Debug("websocket logging ready")
	}
}

type WebsocketReceiver struct {
}

func NewWebsocketWriter() (writer *WebsocketReceiver, err error) {
	newWriter := new(WebsocketReceiver)
	return newWriter, nil
}

func (console *WebsocketReceiver) Write(bytes []byte) (int, error) {
	if websocketHandler != nil {
		websocketHandler(string(bytes), log.DebugLvl, nil)
	}
	return 0,nil
}

func (console *WebsocketReceiver) String() string {
	return "Websocket writer"
}
