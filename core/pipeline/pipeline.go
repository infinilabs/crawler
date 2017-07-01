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
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/util"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

type ParaKey string

type Phrase int

type Context struct {
	Simulate bool `json:"simulate"`
	Parameters
	Phrase       Phrase      `json:"phrase"`
	IgnoreBroken bool        `json:"ignore_broken"`
	breakFlag    bool        `json:"-"`
	exitFlag     bool        `json:"-"`
	Payload      interface{} `json:"-"`
}

/**
break all pipelines, but the end phrase not included
*/
func (this *Context) Break(msg interface{}) {
	log.Trace("break,", this, ",", msg)
	if this == nil {
		panic(errors.New("context is nil"))
	}
	this.breakFlag = true
	this.Payload = msg
}

func (this *Context) IsBreak() bool {
	return this.breakFlag
}

/**
break all pipelines, without execute the end phrase
*/
func (this *Context) IsErrorExit() bool {
	return this.exitFlag
}

/**
tell pipeline to exit all
*/
func (this *Context) ErrorExit(msg interface{}) {
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

func (this *Parameters) MustGet(key ParaKey) interface{} {
	this.Init()

	s := string(key)

	this.l.RLock()
	v, ok := this.Data[s]
	this.l.RUnlock()

	if !ok {
		log.Debug(util.ToJson(this.Data, true))
		panic(fmt.Errorf("%s not found in context", key))
	}

	return v
}

func (this *Parameters) GetMap(key ParaKey) (map[string]interface{}, bool) {
	v := this.Get(key)
	s, ok := v.(map[string]interface{})
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

func (this *Parameters) GetOrDefault(key ParaKey, val interface{}) interface{} {
	this.Init()
	this.l.RLock()
	s := string(key)
	v := this.Data[s]
	this.l.RUnlock()
	if v == nil {
		return val
	}
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
		log.Debug(util.ToJson(this.Data, true))
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

func (this *Parameters) GetStringOrDefault(key ParaKey, val string) string {
	s, ok := this.GetString(key)
	if (!ok) || len(s) == 0 {
		return val
	}
	return s
}

func (this *Parameters) MustGetBytes(key ParaKey) []byte {
	s, ok := this.Get(key).([]byte)
	if !ok {
		log.Debug(util.ToJson(this.Data, true))
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

/*
return 0 if not key was found
*/
func (this *Parameters) MustGetInt(key ParaKey) int {
	s, _ := this.GetInt(key)
	return s
}

func (this *Parameters) MustGetMap(key ParaKey) map[string]interface{} {
	s, ok := this.GetMap(key)
	if !ok {
		log.Debug(util.ToJson(this.Data, true))
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

type Joint interface {
	Name() string
	Process(s *Context) error
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
	pipe.id = util.GetIncrementID("pipe")
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

func (this *Pipeline) GetContext() *Context {
	return this.context
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
				if e, ok := r.(runtime.Error); ok {
					log.Errorf("%v", r)
					this.context.Break(util.GetRuntimeErrorMessage(e))
				}
				log.Debug("error in pipeline, ", this.name)
				stats.Increment(this.name+".pipeline", "error")
			}
		}
		if !this.context.IsErrorExit() && (!(this.context.IgnoreBroken && this.context.IsBreak())) {
			this.endPipeline()
		}

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
		if this.context.IsErrorExit() {
			if global.Env().IsDebug {
				log.Debug(util.ToJson(this.id, true))
				log.Debug(util.ToJson(this.name, true))
				log.Debug(util.ToJson(this.context, true))
			}
			log.Trace("exit joint,", v.Name())
			stats.Increment(this.name+".pipeline", "exit")
			break
		}
		startTime := time.Now()
		err = v.Process(this.context)
		elapsedTime := time.Now().Sub(startTime)
		stats.Timing(this.name+".pipeline", v.Name(), elapsedTime.Nanoseconds())
		if err != nil {
			stats.Increment(this.name+".pipeline", "error")
			log.Debug("%s-%s: %v", this.name, v.Name(), err)
			this.context.Break(err)
			panic(err)
		}
		log.Trace(this.name, ", end joint,", v.Name())
	}

	return this.context
}
func (this *Pipeline) endPipeline() {
	if this.context.IsErrorExit() {
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
	log.Tracef("pipeline config: %v", util.ToJson(config, true))

	pipe := &Pipeline{}
	pipe.id = util.GetIncrementID("pipe")
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

	log.Tracef("get pipeline: %v", util.ToJson(pipe, true))

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
