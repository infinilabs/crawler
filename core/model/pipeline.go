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

package model

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
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

//  Common pipeline context keys
const (
	CONTEXT_SNAPSHOT   ParaKey = "SNAPSHOT"
	CONTEXT_PAGE_LINKS ParaKey = "PAGE_LINKS"
)

type Context struct {
	Parameters

	SequenceID       int64       `json:"sequence"`
	PipelineConfigID string      `json:"pipeline_config_id"`
	Simulate         bool        `json:"simulate"`
	IgnoreBroken     bool        `json:"ignore_broken"`
	Payload          interface{} `json:"-"`

	//private parameters
	breakFlag bool
	exitFlag  bool
}

// End break all pipelines, but the end phrase not included
func (context *Context) End(msg interface{}) {
	log.Trace("break,", context, ",", msg)
	if context == nil {
		panic(errors.New("context is nil"))
	}
	context.breakFlag = true
	context.Payload = msg
}

// IsEnd indicates whether the pipe process is end, end means no more processes will be execute
func (context *Context) IsEnd() bool {
	return context.breakFlag
}

// IsExit means all pipelines will be broke and jump to outside, even the end phrase will not be executed as well
func (context *Context) IsExit() bool {
	return context.exitFlag
}

// Exit tells pipeline to exit
func (context *Context) Exit(msg interface{}) {
	context.exitFlag = true
	context.Payload = msg
}

type Parameters struct {
	Data   map[string]interface{} `json:"data"`
	l      *sync.RWMutex
	inited bool
}

func (para *Parameters) Init() {
	if para.inited {
		return
	}
	if para.l == nil {
		para.l = &sync.RWMutex{}
	}
	para.l.Lock()
	if para.Data == nil {
		para.Data = map[string]interface{}{}
	}
	para.inited = true
	para.l.Unlock()
}

func (para *Parameters) MustGetTime(key ParaKey) time.Time {
	v, ok := para.GetTime(key)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return v
}

func (para *Parameters) GetTime(key ParaKey) (time.Time, bool) {
	v := para.Get(key)
	s, ok := v.(time.Time)
	if ok {
		return s, ok
	}
	return s, ok
}

func (para *Parameters) GetString(key ParaKey) (string, bool) {
	v := para.Get(key)
	s, ok := v.(string)
	if ok {
		return s, ok
	}
	return s, ok
}

func (para *Parameters) GetBool(key ParaKey, defaultV bool) bool {
	v := para.Get(key)
	s, ok := v.(bool)
	if ok {
		return s
	}
	return defaultV
}

func (para *Parameters) Has(key ParaKey) bool {
	para.Init()
	_, ok := para.Data[string(key)]
	return ok
}

func (para *Parameters) GetIntOrDefault(key ParaKey, defaultV int) int {
	v, ok := para.GetInt(key, defaultV)
	if ok {
		return v
	}
	return defaultV
}

func (para *Parameters) GetInt(key ParaKey, defaultV int) (int, bool) {
	v, ok := para.GetInt64(key, 0)
	if ok {
		return int(v), ok
	}
	return defaultV, ok
}

func (para *Parameters) GetInt64OrDefault(key ParaKey, defaultV int64) int64 {
	v, ok := para.GetInt64(key, defaultV)
	if ok {
		return v
	}
	return defaultV
}

func (para *Parameters) GetInt64(key ParaKey, defaultV int64) (int64, bool) {
	v := para.Get(key)

	s, ok := v.(int64)
	if ok {
		return s, ok
	}

	s1, ok := v.(uint64)
	if ok {
		return int64(s1), ok
	}

	s2, ok := v.(int)
	if ok {
		return int64(s2), ok
	}

	s3, ok := v.(uint)
	if ok {
		return int64(s3), ok
	}

	return defaultV, ok
}

func (para *Parameters) MustGet(key ParaKey) interface{} {
	para.Init()

	s := string(key)

	para.l.RLock()
	v, ok := para.Data[s]
	para.l.RUnlock()

	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}

	return v
}

func (para *Parameters) GetMap(key ParaKey) (map[string]interface{}, bool) {
	v := para.Get(key)
	s, ok := v.(map[string]interface{})
	return s, ok
}

// GetStringArray will return a array which type of the items are string
func (para *Parameters) GetStringArray(key ParaKey) ([]string, bool) {
	v := para.Get(key)
	s, ok := v.([]string)
	return s, ok
}

func (para *Parameters) Get(key ParaKey) interface{} {
	para.Init()
	para.l.RLock()
	s := string(key)
	v := para.Data[s]
	para.l.RUnlock()
	return v
}

func (para *Parameters) GetOrDefault(key ParaKey, val interface{}) interface{} {
	para.Init()
	para.l.RLock()
	s := string(key)
	v := para.Data[s]
	para.l.RUnlock()
	if v == nil {
		return val
	}
	return v
}

func (para *Parameters) Set(key ParaKey, value interface{}) {
	para.Init()
	para.l.Lock()
	s := string(key)
	para.Data[s] = value
	para.l.Unlock()
}

func (para *Parameters) MustGetString(key ParaKey) string {
	s, ok := para.GetString(key)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

func (para *Parameters) GetStringOrDefault(key ParaKey, val string) string {
	s, ok := para.GetString(key)
	if (!ok) || len(s) == 0 {
		return val
	}
	return s
}

func (para *Parameters) MustGetBytes(key ParaKey) []byte {
	s, ok := para.Get(key).([]byte)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

// MustGetInt return 0 if not key was found
func (para *Parameters) MustGetInt(key ParaKey) int {
	v, ok := para.GetInt(key, 0)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return v
}

func (para *Parameters) MustGetInt64(key ParaKey) int64 {
	s, ok := para.GetInt64(key, 0)
	if !ok {
		panic(fmt.Errorf("%s not found in context", key))
	}
	return s
}

func (para *Parameters) MustGetMap(key ParaKey) map[string]interface{} {
	s, ok := para.GetMap(key)
	if !ok {
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

	currentJointName string
}

func NewPipeline(name string) *Pipeline {
	pipe := &Pipeline{}
	pipe.id = util.GetIncrementID("pipe")
	pipe.name = strings.TrimSpace(name)
	pipe.context = &Context{}
	pipe.context.Parameters.Init()
	return pipe
}

func (pipe *Pipeline) Context(s *Context) *Pipeline {
	if s != nil {
		pipe.context = s
		pipe.context.Init()
	}

	return pipe
}

func (pipe *Pipeline) GetID() string {
	return pipe.id
}

func (pipe *Pipeline) GetContext() *Context {
	return pipe.context
}

func (pipe *Pipeline) Start(s Joint) *Pipeline {
	pipe.joints = []Joint{s}
	return pipe
}

func (pipe *Pipeline) Join(s Joint) *Pipeline {
	pipe.joints = append(pipe.joints, s)
	return pipe
}

func (pipe *Pipeline) End(s Joint) *Pipeline {
	pipe.endJoint = s
	return pipe
}

// setCurrentJoint set current joint's name, used for debugging
func (context *Pipeline) setCurrentJoint(name string) {
	context.currentJointName = name
}
func (pipe *Pipeline) GetCurrentJoint() string {
	return pipe.currentJointName
}

func (pipe *Pipeline) Run() *Context {

	stats.Increment(pipe.name+".pipeline", "total")

	//final phrase
	defer func() {
		if !global.Env().IsDebug {
			if r := recover(); r != nil {

				if r == nil {
					return
				}
				var v string
				switch r.(type) {
				case error:
					v = r.(error).Error()
				case runtime.Error:
					v = r.(runtime.Error).Error()
				case string:
					v = r.(string)
				}
				log.Error("error in pipeline, ", pipe.name, ", ", pipe.id, ", ", pipe.currentJointName, ", ", v)
				stats.Increment(pipe.name+".pipeline", "error")
			}
		}
		if !pipe.context.IsExit() && (!(pipe.context.IgnoreBroken && pipe.context.IsEnd())) {
			pipe.endPipeline()
		}

		stats.Increment(pipe.name+".pipeline", "finished")
	}()

	var err error
	for _, v := range pipe.joints {
		log.Trace("pipe, ", pipe.name, ", start joint,", v.Name())
		if pipe.context.IsEnd() {
			log.Trace("break joint,", v.Name())
			stats.Increment(pipe.name+".pipeline", "break")
			break
		}
		if pipe.context.IsExit() {
			if global.Env().IsDebug {
				log.Debug(util.ToJson(pipe.id, true))
				log.Debug(util.ToJson(pipe.name, true))
				log.Debug(util.ToJson(pipe.context, true))
			}
			log.Trace("exit joint,", v.Name())
			stats.Increment(pipe.name+".pipeline", "exit")
			break
		}
		pipe.setCurrentJoint(v.Name())
		startTime := time.Now()
		err = v.Process(pipe.context)

		elapsedTime := time.Now().Sub(startTime)
		stats.Timing(pipe.name+".pipeline", v.Name(), elapsedTime.Nanoseconds())
		if err != nil {
			stats.Increment(pipe.name+".pipeline", "error")
			log.Debug("%s-%s: %v", pipe.name, v.Name(), err)
			break
		}
		log.Trace(pipe.name, ", end joint,", v.Name())
	}

	return pipe.context
}

func (pipe *Pipeline) endPipeline() {
	if pipe.context.IsExit() {
		log.Debug("exit pipeline, ", pipe.name, ", ", pipe.context.Payload)
		return
	}

	log.Trace("start finish pipeline, ", pipe.name)
	if pipe.endJoint != nil {
		pipe.setCurrentJoint(pipe.endJoint.Name())
		pipe.endJoint.Process(pipe.context)
	}
	log.Trace("end finish pipeline, ", pipe.name)
}

func NewPipelineFromConfig(name string, config *PipelineConfig, context *Context) *Pipeline {
	if global.Env().IsDebug {
		log.Debugf("pipeline config: %v", util.ToJson(config, true))
	}

	pipe := &Pipeline{}
	pipe.id = util.GetIncrementID("pipe")
	pipe.name = strings.TrimSpace(name)

	pipe.Context(context)

	if config.StartJoint != nil && config.StartJoint.Enabled {
		input := GetJointInstance(config.StartJoint)
		pipe.Start(input)
	}

	for _, cfg := range config.ProcessJoints {
		if cfg.Enabled {
			j := GetJointInstance(cfg)
			pipe.Join(j)
		}
	}

	if config.EndJoint != nil && config.EndJoint.Enabled {
		output := GetJointInstance(config.EndJoint)
		pipe.End(output)
	}

	if global.Env().IsDebug {
		log.Debugf("get pipeline: %v", util.ToJson(pipe, true))
	}

	return pipe
}

var typeRegistry = make(map[string]interface{})

func GetAllRegisteredJoints() map[string]interface{} {
	return typeRegistry
}

func GetJointInstance(cfg *JointConfig) Joint {
	log.Trace("get joint instances, ", cfg.JointName)
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

func RegisterPipeJoint(joint Joint) {
	k := string(joint.Name())
	RegisterPipeJointWithName(k, joint)
}

func RegisterPipeJointWithName(jointName string, joint Joint) {
	if typeRegistry[jointName] != nil {
		panic(errors.Errorf("joint with same name already registered, %s", jointName))
	}
	typeRegistry[jointName] = joint
}
