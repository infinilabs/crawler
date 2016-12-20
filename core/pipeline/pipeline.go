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

package pipeline

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/stats"
	"github.com/rs/xid"
	"runtime"
	"strings"
	"time"
)

type ContextKey string

type Context struct {
	Data      map[ContextKey]interface{}
	Env       *env.Env
	breakFlag bool
	Payload   interface{}
}

/**
break all pipeline
*/
func (this *Context) Break(msg interface{}) {
	this.breakFlag = true
	this.Payload = msg
}

func (this *Context) IsBreak() bool {
	return this.breakFlag
}

func (this *Context) GetString(key ContextKey) (string, bool) {
	v := this.Get(key)
	s, ok := v.(string)
	if ok {
		return s, ok
	}
	return s, ok
}

func (this *Context) GetInt(key ContextKey) (int, bool) {
	v := this.Get(key)
	s, ok := v.(int)
	if ok {
		return s, ok
	}
	return s, ok
}

func (this *Context) MustGetString(key ContextKey) string {
	s, ok := this.GetString(key)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

func (this *Context) MustGetBytes(key ContextKey) []byte {
	s, ok := this.Get(key).([]byte)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

func (this *Context) MustGetInt(key ContextKey) int {
	s, ok := this.GetInt(key)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

func (this *Context) MustGetMap(key ContextKey) map[string]interface{} {
	s, ok := this.GetMap(key)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

func (this *Context) GetMap(key ContextKey) (map[string]interface{}, bool) {
	v := this.Get(key)
	s, ok := v.(map[string]interface{})
	if ok {
		return s, ok
	}
	return s, ok
}

func (this *Context) Get(key ContextKey) interface{} {
	return this.Data[key]
}

func (this *Context) Set(key ContextKey, value interface{}) {
	this.Data[key] = value
}

type Joint interface {
	Name() string
	Process(s *Context) (*Context, error)
}

type Pipeline struct {
	id       string
	name     string
	joints   []Joint
	context  *Context
	endJoint Joint
}

func NewPipeline(name string) *Pipeline {
	pipe := &Pipeline{}
	pipe.id = xid.New().String()
	pipe.name = strings.TrimSpace(name)
	return pipe
}

func (this *Pipeline) Context(s *Context) *Pipeline {
	this.context = s
	return this
}

func (this *Pipeline) Start(s Joint) *Pipeline {

	if this.context == nil {
		this.context = &Context{}
	}
	if this.context.Data == nil {
		this.context.Data = map[ContextKey]interface{}{}
	}
	this.joints = []Joint{s}
	return this
}

func (this *Pipeline) Join(s Joint) *Pipeline {
	this.joints = append(this.joints, s)
	return this
}

//func (this *Pipeline) Err(err error) *Pipeline {
//	//TODO handle error, persist error log
//	log.Error("error in pipeline: ", this.id)
//	this.context.Payload =err
//	return this
//}

func (this *Pipeline) End(s Joint) *Pipeline {
	this.endJoint = s
	return this
}

func (this *Pipeline) Run() *Context {

	stats.Increment(this.name+".pipeline", "total")

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				err := r.(error)
				stats.Increment(this.name+".pipeline", "error")
				log.Errorf("%s: %v",this.name, err)
			}
			this.endPipeline()
			log.Trace("error in pipe")
		}
	}()

	var err error
	for _, v := range this.joints {
		log.Trace("start joint,", v.Name())
		if this.context.IsBreak() {
			log.Trace("break joint,", v.Name())
			stats.Increment(this.name+".pipeline", "break")
			break
		}
		startTime := time.Now()
		this.context, err = v.Process(this.context)
		elapsedTime := time.Now().Sub(startTime)
		stats.Timing(this.name+".pipeline", v.Name(), elapsedTime.Nanoseconds())
		if err != nil {
			stats.Increment(this.name+".pipeline", "error")
			log.Errorf("%s-%s: %v",this.name,v.Name(), err)
			panic(err)
		}
		log.Trace("end joint,", v.Name())
	}

	this.endPipeline()
	stats.Increment(this.name+".pipeline", "finished")
	return this.context
}

func (this *Pipeline) endPipeline() {
	log.Trace("start finish piplne")
	if this.endJoint != nil {
		this.endJoint.Process(this.context)
	}
	log.Trace("end finish piplne")
}
