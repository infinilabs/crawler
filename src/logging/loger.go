
package logging

import (
	log "github.com/cihub/seelog"
)


func Flush(){
	log.Flush()
}

func Trace(v ...interface{}) {
	log.Trace(v)
}

func Debug(v ...interface{}) {
	log.Debug(v)
}

func Info(v ...interface{}) {
	log.Info(v)
}

func Warn(v ...interface{}) {
	log.Warn(v)
}
func Error(v ...interface{}) {
	log.Error(v)
}


func SetInitLogging(logLevel string) {
	testConfig := `
	<seelog  type="sync" minlevel="`
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
	<seelog  type="sync" minlevel="`
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
