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

import (
	"fmt"
	"io"
)

// ErrorCode is enum object of errors
type ErrorCode int

// ErrorPayload is detail payload of errors
type ErrorPayload interface{}

var (
	// Default unknow error
	Default          ErrorCode = 1
	InvalidParameter ErrorCode = 2
	// JSONIsEmpty error when json is empty
	JSONIsEmpty ErrorCode = 100
	// BodyEmpty error when body is empty
	BodyEmpty ErrorCode = 101
	// URLRedirected error when url redirected
	URLRedirected ErrorCode = 102
)

// NewWithCode create a error with error code and message
func NewWithCode(err error, code ErrorCode, msg string) error {
	if err == nil {
		return nil
	}
	return wrapper{
		cause: cause{
			cause: err,
			msg:   msg,
			code:  code,
		},
		stack: callers(),
	}
}

// NewWithPayload create error with error code and payload and message
func NewWithPayload(err error, code ErrorCode, payload interface{}, msg string) error {
	if err == nil {
		return nil
	}
	return wrapper{
		cause: cause{
			cause:   err,
			msg:     msg,
			code:    code,
			payload: payload,
		},
		stack: callers(),
	}
}

// _error is an error implementation returned by New and Errorf
// that implements its own fmt.Formatter.
type _error struct {
	msg string
	*stack
}

func (e _error) Error() string { return e.msg }

func (e _error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.msg)
			fmt.Fprintf(s, "%+v", e.StackTrace())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	}
}

// New returns an error with the supplied message.
func New(message string) error {
	return _error{
		message,
		callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Errorf(format string, args ...interface{}) error {
	return _error{
		fmt.Sprintf(format, args...),
		callers(),
	}
}

type cause struct {
	code    ErrorCode
	cause   error
	msg     string
	payload interface{}
}

func (c cause) Error() string        { return fmt.Sprintf("%s: %v", c.msg, c.Cause()) }
func (c cause) Cause() error         { return c.cause }
func (c cause) Code() ErrorCode      { return c.code }
func (c cause) Payload() interface{} { return c.payload }

// wrapper is an error implementation returned by Wrap and Wrapf
// that implements its own fmt.Formatter.
type wrapper struct {
	cause
	*stack
}

func (w wrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			fmt.Fprintf(s, "%+v: %s", w.StackTrace()[0], w.msg)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	}
}

// Wrap returns an error annotating err with message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return wrapper{
		cause: cause{
			cause: err,
			msg:   message,
		},
		stack: callers(),
	}
}

// Wrapf returns an error annotating err with the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return wrapper{
		cause: cause{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		stack: callers(),
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type Causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// Code return error code
func Code(err error) ErrorCode {
	code := Default
	type causer interface {
		Code() ErrorCode
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		code = cause.Code()
	}
	return code
}

// CodeWithPayload return error code and payload
func CodeWithPayload(err error) (ErrorCode, interface{}) {
	type causer interface {
		Code() ErrorCode
		Payload() interface{}
	}

	if err != nil {
		cause, ok := err.(causer)
		if ok {
			code := cause.Code()
			payload := cause.Payload()
			return code, payload
		}
	}
	return Default, nil
}
