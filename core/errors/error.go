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

package errors

type ErrorCode int

var (
	Default ErrorCode = 1
	JSONIsEmpty ErrorCode = 100
	BodyEmpty  ErrorCode =101
	URLRedirected  ErrorCode =102
)


type Error struct {
	Code ErrorCode
	Message string
	InnerError error
	Payload interface{}
}

func (this *Error) Error()string  {
	if(this.InnerError!=nil){
		return this.InnerError.Error()
	}
	return this.Message
}

func NewWithCode(code ErrorCode,msg string) error {
	return &Error{Code:code,Message:msg}
}

func NewWithPayload(code ErrorCode,msg string,payload interface{}) error {
	return &Error{Code:code,Message:msg,Payload:payload}
}

func New(text string) error {
	return &Error{Code:Default,Message:text}
}
