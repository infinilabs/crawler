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
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/env"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/ryanuber/go-glob"
	"strings"
	"sync"
)

var file string
var loggingConfig *config.LoggingConfig
var l sync.Mutex
var e *env.Env

// SetLogging init set logging
func SetLogging(env *env.Env, logLevel string, logFile string) {

	e = env

	l.Lock()
	if loggingConfig == nil {
		loggingConfig = &config.LoggingConfig{}
		loggingConfig.LogLevel = "info"
		loggingConfig.PushLogLevel = "info"
		loggingConfig.FileFilterPattern = "*"
		loggingConfig.FuncFilterPattern = "*"
	}
	l.Unlock()

	if env != nil {
		envLevel := strings.ToLower(env.LoggingLevel)
		if env.SystemConfig != nil {
			file = env.SystemConfig.PathConfig.Log + "/gopa.log"
		}
		if len(envLevel) > 0 {
			loggingConfig.LogLevel = envLevel
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
	if file == "" {
		file = "./log/gopa.log"
	}

	if loggingConfig.FuncFilterPattern == "" {
		loggingConfig.FuncFilterPattern = "*"
	}
	if loggingConfig.FileFilterPattern == "" {
		loggingConfig.FileFilterPattern = "*"
	}
	if loggingConfig.LogLevel == "" {
		loggingConfig.LogLevel = "info"
	}
	if loggingConfig.PushLogLevel == "" {
		loggingConfig.PushLogLevel = "info"
	}

	consoleWriter, _ := NewConsoleWriter()

	format := "[%Date(01-02) %Time] [%LEV] [%File:%Line] %Msg%n"
	formatter, err := log.NewFormatter(format)
	if err != nil {
		fmt.Println(err)
	}

	rollingFileWriter, _ := NewRollingFileWriterSize(file, rollingArchiveNone, "", 10000000000, 5, rollingNameModePostfix)
	bufferedWriter, _ := NewBufferedWriter(rollingFileWriter, 10000, 1000)

	l, _ := log.LogLevelFromString(strings.ToLower(loggingConfig.LogLevel))
	pushl, _ := log.LogLevelFromString(strings.ToLower(loggingConfig.PushLogLevel))

	//logging receivers
	receivers := []interface{}{consoleWriter, bufferedWriter}

	root, err := log.NewSplitDispatcher(formatter, receivers)
	if err != nil {
		fmt.Println(err)
	}

	golbalConstraints, err := log.NewMinMaxConstraints(l, log.CriticalLvl)
	if err != nil {
		fmt.Println(err)
	}

	exceptions := []*log.LogLevelException{}

	if loggingConfig.RealtimePushEnabled {

		logger, err := log.LoggerFromCustomReceiver(&CustomReceiver{config: loggingConfig, minLogLevel: l, pushminLogLevel: pushl})
		err = log.ReplaceLogger(logger)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		logger := log.NewAsyncLoopLogger(log.NewLoggerConfig(golbalConstraints, exceptions, root))
		err = log.ReplaceLogger(logger)
		if err != nil {
			fmt.Println(err)
		}
	}

}

// GetLoggingConfig return logging configs
func GetLoggingConfig() *config.LoggingConfig {
	return loggingConfig
}

// UpdateLoggingConfig update logging config
func UpdateLoggingConfig(config *config.LoggingConfig) {
	l.Lock()
	loggingConfig = config
	l.Unlock()
	SetLogging(e, "", "")
}

// Flush is flush logs to output
func Flush() {
	log.Flush()
}

var websocketHandler func(message string, level log.LogLevel, context log.LogContextInterface)

// RegisterWebsocketHandler used to register websocket handler
func RegisterWebsocketHandler(func1 func(message string, level log.LogLevel, context log.LogContextInterface)) {

	websocketHandler = func1
	if func1 != nil {
		log.Debug("websocket logging ready")
	}
}

// CustomReceiver is a struct of custom log receiver, which implements seelog.CustomReceiver
type CustomReceiver struct {
	config          *config.LoggingConfig
	minLogLevel     log.LogLevel
	pushminLogLevel log.LogLevel
}

// ReceiveMessage impl how to receive log message
func (ar *CustomReceiver) ReceiveMessage(message string, level log.LogLevel, context log.LogContextInterface) error {

	//truncate huge message
	if len(message) > 300 {
		message = util.SubString(message, 0, 300) + "..."
	}

	f := context.Func()
	spl := strings.Split(f, ".")
	funcName := spl[len(spl)-1]

	preparedMessage := fmt.Sprintf("[%s] [%s] [%s:%d] [%s] %s\n",
		context.CallTime().Format("15:04:05"),
		strings.ToUpper(level.String()),
		context.FileName(),
		context.Line(),
		funcName,
		message,
	)

	//console output
	if level >= ar.minLogLevel {
		fmt.Printf(preparedMessage)
	}

	if ar.config != nil {
		if level < ar.pushminLogLevel {
			return nil
		}

		if len(ar.config.FileFilterPattern) > 0 && ar.config.FileFilterPattern != "*" {
			if !glob.Glob(ar.config.FileFilterPattern, context.FileName()) {
				return nil
			}
		}
		if len(ar.config.FuncFilterPattern) > 0 && ar.config.FuncFilterPattern != "*" {
			if !glob.Glob(ar.config.FuncFilterPattern, funcName) {
				return nil
			}
		}
		if len(ar.config.MessageFilterPattern) > 0 && ar.config.MessageFilterPattern != "*" {
			if !glob.Glob(ar.config.MessageFilterPattern, message) {
				return nil
			}
		}
	}

	//push message to websocket
	if websocketHandler != nil {

		websocketHandler(preparedMessage, log.DebugLvl, nil)
	}

	return nil
}

// AfterParse nothing to do here
func (ar *CustomReceiver) AfterParse(initArgs log.CustomReceiverInitArgs) error {
	return nil
}

// Flush logs
func (ar *CustomReceiver) Flush() {

}

// Close logs
func (ar *CustomReceiver) Close() error {
	return nil
}
