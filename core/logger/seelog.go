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
	"strings"
)

var config string
var level string
var file string


func SetLogging(env *Env, logLevel string, logFile string) {
	if(env!=nil){
		envLevel := strings.ToLower(env.LoggingLevel)
		if(env.SystemConfig!=nil){
			envLogFile := env.SystemConfig.Log + "/gopa.log"
			if(len(envLogFile)>0){
				file=envLogFile
			}
		}
		if(len(envLevel)>0){
			level=envLevel
		}
	}

	if len(logLevel) > 0 {
		level = strings.ToLower(logLevel)
	}
	if len(level) <= 0 {
		level = "info"
	}

	if len(file) <= 0 {
		logFile = "./log/gopa.log"
	}

	consoleWriter, _ := NewConsoleWriter()
	websocketWriter, _ := NewWebsocketWriter()

	formatter, _ := log.NewFormatter("[%Date(01-02) %Time] [%LEV] [%File:%Line] %Msg%n")

	NewRollingFileWriterSize(logFile, rollingArchiveNone, "", 10000000000, 5, rollingNameModePostfix)

	root, _ := log.NewSplitDispatcher(formatter, []interface{}{websocketWriter,consoleWriter})

	l,_:=log.LogLevelFromString(level)
	constraints, _ := log.NewMinMaxConstraints(l, log.CriticalLvl)

	specificConstraints, _ := log.NewListConstraints([]log.LogLevel{log.TraceLvl, log.ErrorLvl})

	ex, _ := log.NewLogLevelException("*", "*crawler.go", specificConstraints)

	exceptions := []*log.LogLevelException{ex}

	logger := log.NewAsyncLoopLogger(log.NewLoggerConfig(constraints, exceptions, root))

	log.ReplaceLogger(logger)
}

func GetLoggingConfig(env *Env) string {
	return config
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
