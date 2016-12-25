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
)

type Logger interface {
	Trace(v ...interface{})

	Debug(v ...interface{})

	Info(v ...interface{})

	Warn(v ...interface{}) error

	Error(v ...interface{}) error

	Critical(v ...interface{}) error
}

var logger Logger

func Trace(v ...interface{}) {
	logger.Trace(v)
}

func Debug(v ...interface{}) {
	logger.Debug(v)
}

func Info(v ...interface{}) {
	logger.Info(v)
}

func Warn(v ...interface{}) error {
	return logger.Warn(v)
}

func Error(v ...interface{}) error {
	return logger.Error(v)
}

func Critical(v ...interface{}) error {
	return logger.Critical(v)
}

func Register(l Logger) {
	logger = l
}
