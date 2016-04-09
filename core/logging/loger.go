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

package logging

import (
	log "github.com/cihub/seelog"
)

type Logger interface {
	Trace(v ...interface{})

	Debug(v ...interface{})

	Info(v ...interface{})

	Warn(v ...interface{}) error

	Error(v ...interface{}) error

	Critical(v ...interface{}) error
}

type RealLogger struct{}

func (rl RealLogger) Trace(v ...interface{}) {
	log.Trace(v)
}

func (rl RealLogger) Debug(v ...interface{}) {
	log.Debug(v)
}

func (rl RealLogger) Info(v ...interface{}) {
	log.Info(v)
}

func (rl RealLogger) Warn(v ...interface{}) error {
	return log.Warn(v)
}

func (rl RealLogger) Error(v ...interface{}) error {
	return log.Error(v)
}

func (rl RealLogger) Critical(v ...interface{}) error {
	return log.Critical(v)
}

type NullLogger struct{}

func (nl NullLogger) Trace(v ...interface{}) {

}

func (nl NullLogger) Debug(v ...interface{}) {

}

func (nl NullLogger) Info(v ...interface{}) {

}

func (nl NullLogger) Warn(v ...interface{}) error {
	return nil
}

func (nl NullLogger) Error(v ...interface{}) error {
	return nil
}

func (nl NullLogger) Critical(v ...interface{}) error {
	return nil
}
