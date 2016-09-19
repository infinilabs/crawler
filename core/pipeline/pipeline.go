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
	"github.com/medcl/gopa/core/env"
	"fmt"
	"time"
	"github.com/medcl/gopa/core/stats"
)
type ContextKey string

//func (this ContextKey) String() string {
//	return string(this)
//}

type Context struct {
	Data map[ContextKey]interface{}
	Env *env.Env
	breakFlag bool
}

func (this *Context) Break(){
	this.breakFlag=true
}

func (this *Context) IsBreak()(bool){
	return this.breakFlag
}

func (this *Context) GetString(key ContextKey)(string,bool){
	v:=this.Get(key)
	s,ok:=v.(string)
	if(ok){
		return s,ok
	}
	return s,ok
}
func (this *Context) GetInt(key ContextKey)(int,bool){
	v:=this.Get(key)
	s,ok:=v.(int)
	if(ok){
		return s,ok
	}
	return s,ok
}

func (this *Context) MustGetString(key ContextKey)(string){
	s,ok:=this.GetString(key)
	if(!ok){
		panic(fmt.Errorf("%s not found in context",key))
	}
	return s
}

func (this *Context) MustGetInt(key ContextKey)(int){
	s,ok:=this.GetInt(key)
	if(!ok){
		panic(fmt.Errorf("%s not found in context",key))
	}
	return s
}

func (this *Context) MustGetMap(key ContextKey)(map[string]interface{}){
	s,ok:=this.GetMap(key)
	if(!ok){
		panic(fmt.Errorf("%s not found in context",key))
	}
	return s
}

func (this *Context) GetMap(key ContextKey)(map[string]interface{},bool){
	v:=this.Get(key)
	s,ok:=v.(map[string]interface{})
	if(ok){
		return s,ok
	}
	return s,ok
}

func (this *Context) Get(key ContextKey)interface{}{
	return this.Data[key]
}


func (this *Context) Set(key ContextKey,value interface{}){
	this.Data[key]=value
}

type Joint interface {
	Name()string
	Process(s *Context) (*Context, error)
}

type Pipeline struct {
	joints  []Joint
	context *Context
}

func (this *Pipeline) Context(s *Context) *Pipeline {
	this.context = s
	return this
}

func (this *Pipeline) Start(s Joint) *Pipeline {
	if(this.context ==nil){
		this.context =&Context{}
	}
	if(this.context.Data==nil){
		this.context.Data =map[ContextKey]interface{}{}
	}
	this.joints = []Joint{s}
	return this
}

func (this *Pipeline) Join(s Joint) *Pipeline {
	this.joints = append(this.joints, s)
	return this
}

func (this *Pipeline) End() *Pipeline {
	return this
}

func (this *Pipeline) Run()(*Context) {
	stats.Increment("crawler.pipeline","total")

	var err error
	for _, v := range this.joints {
		if(this.context.breakFlag){
			break
		}
		startTime := time.Now()
		this.context, err = v.Process(this.context)
		elapsedTime:=time.Now().Sub(startTime)
		stats.Timing("crawler.pipeline",v.Name(),elapsedTime.Nanoseconds())
		if err != nil {
			stats.Increment("crawler.pipeline","error")
			panic(err)
		}
	}
	return this.context
}
