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
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/core/stats"
	"github.com/rs/xid"
	"runtime"
	"strings"
	"sync"
	"time"
)

var l sync.RWMutex

type ContextKey string

type Context struct {
	Phrase    model.TaskPhrase
	data      map[ContextKey]interface{}
	breakFlag bool
	exitFlag  bool
	Payload   interface{}
}

/**
break all pipelines, but the end phrase not included
*/
func (this *Context) Break(msg interface{}) {
	log.Trace("break,", msg)
	this.breakFlag = true
	this.Payload = msg
}

func (this *Context) Init() {
	l.Lock()
	if this.data == nil {
		this.data = map[ContextKey]interface{}{}
	}
	l.Unlock()
}

func (this *Context) IsBreak() bool {
	return this.breakFlag
}

/**
break all pipelines, without execute the end phrase
*/
func (this *Context) IsExit() bool {
	return this.exitFlag
}

/**
tell pipeline to exit all
*/
func (this *Context) Exit(msg interface{}) {
	this.exitFlag = true
	this.Payload = msg
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
	l.RLock()
	v := this.data[key]
	l.RUnlock()
	return v
}

func (this *Context) Set(key ContextKey, value interface{}) {
	l.Lock()
	this.data[key] = value
	l.Unlock()
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
	pipe.context = &Context{}
	pipe.context.Init()
	return pipe
}

func (this *Pipeline) Context(s *Context) *Pipeline {
	this.context = s
	this.context.Init()
	return this
}

func (this *Pipeline) GetID() string {
	return this.id
}

func (this *Pipeline) Start(s Joint) *Pipeline {
	this.joints = []Joint{s}
	return this
}

func (this *Pipeline) Join(s Joint) *Pipeline {
	this.joints = append(this.joints, s)
	return this
}

func (this *Pipeline) End(s Joint) *Pipeline {
	this.endJoint = s
	return this
}

func (this *Pipeline) Run() *Context {

	stats.Increment(this.name+".pipeline", "total")

	//final phrase
	defer func() {
		if !global.Env().IsDebug {
			if r := recover(); r != nil {
				if _, ok := r.(runtime.Error); ok {
					err := r.(error)
					log.Errorf("%s: %v", this.name, err)
					//this.context.Break(err.Error())
				}
				log.Trace("error in pipeline, ", this.name)
				stats.Increment(this.name+".pipeline", "error")
			}
		}

		this.endPipeline()
		stats.Increment(this.name+".pipeline", "finished")
	}()

	var err error
	for _, v := range this.joints {
		log.Trace("pipe, ", this.name, ", start joint,", v.Name())
		if this.context.IsBreak() {
			log.Trace("break joint,", v.Name())
			stats.Increment(this.name+".pipeline", "break")
			break
		}
		if this.context.IsExit() {
			log.Trace("exit joint,", v.Name())
			stats.Increment(this.name+".pipeline", "exit")
			break
		}
		startTime := time.Now()
		this.context, err = v.Process(this.context)
		elapsedTime := time.Now().Sub(startTime)
		stats.Timing(this.name+".pipeline", v.Name(), elapsedTime.Nanoseconds())
		if err != nil {
			stats.Increment(this.name+".pipeline", "error")
			log.Errorf("%s-%s: %v", this.name, v.Name(), err)
			this.context.Break(err.Error())
			panic(err)
		}
		log.Trace(this.name, ", end joint,", v.Name())
	}

	return this.context
}

func (this *Pipeline) endPipeline() {
	if(this.context.IsExit()){
		log.Debug("exit pipeline, ", this.name,", ",this.context.Payload)
		return
	}

	log.Trace("start finish pipeline, ", this.name)
	if this.endJoint != nil {
		this.endJoint.Process(this.context)
	}
	log.Trace("end finish pipeline, ", this.name)
}
