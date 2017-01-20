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
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/core/stats"
	"github.com/rs/xid"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

type ParaKey string

type Context struct {
	DryRun bool
	Parameters
	Phrase    model.TaskPhrase `json:"phrase"`
	breakFlag bool             `json:"-"`
	exitFlag  bool             `json:"-"`
	Payload   interface{}      `json:"-"`
}

/**
break all pipelines, but the end phrase not included
*/
func (this *Context) Break(msg interface{}) {
	log.Trace("break,", msg)
	this.breakFlag = true
	this.Payload = msg
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

type Parameters struct {
	Data   map[string]interface{} `json:"data"`
	l      sync.RWMutex
	inited bool
}

func (this *Parameters) Init() {
	if this.inited {
		return
	}
	this.l.Lock()
	if this.Data == nil {
		this.Data = map[string]interface{}{}
	}
	this.inited = true
	this.l.Unlock()
}

func (this *Parameters) GetString(key ParaKey) (string, bool) {
	v := this.Get(key)
	s, ok := v.(string)
	if ok {
		return s, ok
	}
	return s, ok
}

func (this *Parameters) Has(key ParaKey) bool {
	this.Init()
	_, ok := this.Data[string(key)]
	return ok
}

func (this *Parameters) GetInt(key ParaKey) (int, bool) {
	v := this.Get(key)
	s, ok := v.(int)
	if ok {
		return s, ok
	}
	return s, ok
}

func (this *Parameters) GetMap(key ParaKey) (map[string]interface{}, bool) {
	v := this.Get(key)
	s, ok := v.(map[string]interface{})
	if ok {
		return s, ok
	}
	return s, ok
}

func (this *Parameters) Get(key ParaKey) interface{} {
	this.Init()
	this.l.RLock()
	s := string(key)
	v := this.Data[s]
	this.l.RUnlock()
	return v
}

func (this *Parameters) Set(key ParaKey, value interface{}) {
	this.Init()
	this.l.Lock()
	s := string(key)
	this.Data[s] = value
	this.l.Unlock()
}

func (this *Parameters) MustGetString(key ParaKey) string {
	s, ok := this.GetString(key)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

func (this *Parameters) MustGetBytes(key ParaKey) []byte {
	s, ok := this.Get(key).([]byte)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

func (this *Parameters) MustGetInt(key ParaKey) int {
	s, ok := this.GetInt(key)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

func (this *Parameters) MustGetMap(key ParaKey) map[string]interface{} {
	s, ok := this.GetMap(key)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

type Joint interface {
	Name() string
	//Input()map[string]bool
	//Output()map[string]bool
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
	pipe.context.Parameters.Init()
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
	if this.context.IsExit() {
		log.Debug("exit pipeline, ", this.name, ", ", this.context.Payload)
		return
	}

	log.Trace("start finish pipeline, ", this.name)
	if this.endJoint != nil {
		this.endJoint.Process(this.context)
	}
	log.Trace("end finish pipeline, ", this.name)
}

func NewPipelineFromConfig(config *PipelineConfig) *Pipeline {
	pipe := &Pipeline{}
	pipe.id = xid.New().String()
	pipe.name = strings.TrimSpace(config.Name)

	pipe.Context(config.Context)

	if config.InputJoint != nil {
		input := GetJointInstance(config.InputJoint)
		pipe.Start(input)
	}

	for _, cfg := range config.ProcessJoints {
		j := GetJointInstance(cfg)
		pipe.Join(j)
	}

	if config.OutputJoint != nil {
		output := GetJointInstance(config.OutputJoint)
		pipe.End(output)
	}

	return pipe
}

var typeRegistry = make(map[string]interface{})

func GetAllRegisteredJoints() map[string]interface{} {
	return typeRegistry
}

func GetJointInstance(cfg *JointConfig) Joint {
	if typeRegistry[cfg.JointName] != nil {
		t := reflect.ValueOf(typeRegistry[cfg.JointName]).Type()
		v := reflect.New(t).Elem()

		f := v.FieldByName("Data")
		if f.IsValid() && f.CanSet() && f.Kind() == reflect.Map {
			f.Set(reflect.ValueOf(cfg.Parameters))
		}
		v1 := v.Interface().(Joint)
		return v1
	}
	panic(errors.New(cfg.JointName + " not found"))
}

type JointKey string

func Register(jointName JointKey, joint Joint) {
	k := string(jointName)
	RegisterByName(k, joint)
}

func RegisterByName(jointName string, joint Joint) {
	typeRegistry[jointName] = joint
}
