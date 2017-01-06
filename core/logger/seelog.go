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
	. "github.com/medcl/gopa/core/env"
	"strings"
)

func SetInitLogging(env *Env, logLevel string) {

	setLogging(env, logLevel, "./log/gopa.log")
}

func SetLogging(env *Env) {
	logLevel := strings.ToLower(env.LoggingLevel)
	logFile := env.SystemConfig.Log + "/gopa.log"
	setLogging(env, logLevel, logFile)

}

var config string
var env *Env

func setLogging(env *Env, logLevel string, logFile string) {

	logLevel = strings.ToLower(logLevel)

	testConfig := `
	<seelog  type="sync" minlevel="`
	testConfig = testConfig + logLevel
	testConfig = testConfig + `">
		<outputs formatid="main">
			<console formatid="main"/>
			<filter levels="` + logLevel + `">
				<file path="` + logFile + `"/>
			</filter>
			 <rollingfile formatid="main" type="size" filename="` + logFile + `" maxsize="10000000000" maxrolls="5" />
			<custom name="websocket" formatid="main"/>
		</outputs>
		<formats>
			<format id="main" format="[%Date(01-02) %Time] [%LEV] [%File:%Line] %Msg%n"/>
		</formats>
	</seelog>`
	ReplaceConfig(env, testConfig)
}

func ReplaceConfig(e *Env, cfg string) {

	log.RegisterReceiver("websocket", &WebsocketReceiver{})

	logger, err := log.LoggerFromConfigAsBytes([]byte(cfg))
	if err != nil {
		log.Error("replace config error,", err)
		return
	}
	err = log.ReplaceLogger(logger)
	if err != nil {
		log.Error("replace config error,", err)
		return
	}
	config = cfg
	env = e
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
		log.Error("logging func registed")
	}
}

type WebsocketReceiver struct {
}

func (ar *WebsocketReceiver) ReceiveMessage(message string, level log.LogLevel, context log.LogContextInterface) error {
	fmt.Sprintln("custom logging func calling")
	//if websocketHandler != nil {
	//	websocketHandler(message, level, context)
	//	log.Error("logging func called")
	//}
	return nil
}
func (ar *WebsocketReceiver) AfterParse(initArgs log.CustomReceiverInitArgs) error {
	return nil
}
func (ar *WebsocketReceiver) Flush() {

}
func (ar *WebsocketReceiver) Close() error {
	return nil
}
