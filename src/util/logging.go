/**
 * Created with IntelliJ IDEA.
 * User: medcl
 * Date: 13-10-17
 * Time: 下午12:47
 */
package util

import (
	log "github.com/cihub/seelog"

)

func SetInitLogging(logLevel string) {
	testConfig := `
	<seelog type="sync" minlevel="`
	testConfig =testConfig + logLevel
	testConfig =testConfig +`">
		<outputs formatid="main">
			<filter levels="error">
				<file path="./log/filter.log"/>
			</filter>
			<console />
		</outputs>
		<formats>
			<format id="main" format="[%LEV] %Msg%n"/>
		</formats>
	</seelog>`
	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)
}

func SetLogging(logLevel string,logFile string) {
	testConfig := `
	<seelog type="sync" minlevel="`
	testConfig = testConfig + logLevel
	testConfig = testConfig + `">
		<outputs formatid="main">
			<filter levels="`+logLevel+`">
				<file path="`+logFile+`"/>
			</filter>
			<console />
		</outputs>
		<formats>
			<format id="main" format="[%LEV] %Msg%n"/>
		</formats>
	</seelog>`
	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)
}
