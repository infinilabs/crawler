
package logging

import (
	log "github.com/cihub/seelog"
)


func Flush(){
	log.Flush()
}


func SetInitLogging(logLevel string) {
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
			<format id="main" format="[%Date %Time] [%LEV] [%File:%Line,%FuncShort] %Msg%n"/>
		</formats>
	</seelog>`
	logger, _ := log.LoggerFromConfigAsString(testConfig)
	log.ReplaceLogger(logger)
}

func SetLogging(logLevel string,logFile string) {
	testConfig := `
	<seelog  type="sync" minlevel="`
	testConfig = testConfig + logLevel
	testConfig = testConfig + `">
		<outputs formatid="main">
			<console formatid="main"/>
			<filter levels="`+logLevel+`">
				<file path="`+logFile+`"/>
			</filter>
			 <rollingfile formatid="main" type="size" filename="`+logFile+`" maxsize="100" maxrolls="5" />
		</outputs>
		<formats>
			<format id="main" format="[%Date %Time] [%LEV] [%File:%Line,%FuncShort] %Msg%n"/>
		</formats>
	</seelog>`
	logger, _ := log.LoggerFromConfigAsString(testConfig)
	log.ReplaceLogger(logger)
}
