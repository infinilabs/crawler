
package logging

import (
	log "github.com/cihub/seelog"
"strings"
)


func Flush(){
	log.Flush()
}


func SetInitLogging(logLevel string) {

	logLevel = strings.ToLower(logLevel)

	testConfig := `
	<seelog  type="sync" minlevel="`
	testConfig =testConfig + logLevel
	testConfig =testConfig +`">
		<outputs formatid="main">
			<filter levels="error">
				<file path="./log/gopa.log"/>
			</filter>
			<console formatid="main" />
		</outputs>
		<formats>
			<format id="main" format="[%Date(01-02) %Time] [%LEV] [%File:%Line,%FuncShort] %Msg%n"/>
		</formats>
	</seelog>`
	logger, _ := log.LoggerFromConfigAsString(testConfig)
	log.ReplaceLogger(logger)
}

func SetLogging(logLevel string,logFile string) {

	logLevel = strings.ToLower(logLevel)

	testConfig := `
	<seelog  type="sync" minlevel="`
	testConfig = testConfig + logLevel
	testConfig = testConfig + `">
		<outputs formatid="main">
			<console formatid="main"/>
			<filter levels="`+logLevel+`">
				<file path="`+logFile+`"/>
			</filter>
			 <rollingfile formatid="main" type="size" filename="`+logFile+`" maxsize="10000000000" maxrolls="5" />
		</outputs>
		<formats>
			<format id="main" format="[%Date(01-02) %Time] [%LEV] [%File:%Line,%FuncShort] %Msg%n"/>
		</formats>
	</seelog>`
	logger, _ := log.LoggerFromConfigAsString(testConfig)
	log.ReplaceLogger(logger)
}
