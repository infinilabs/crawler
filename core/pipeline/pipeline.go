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
)

type Context struct {
	Data map[string]interface{}
	Env *env.Env
}

func (this *Context) GetString(key string)(string,bool){
	v:=this.Get(key)
	s,ok:=v.(string)
	if(ok){
		return s,ok
	}
	return s,ok
}

func (this *Context) MustGetString(key string)(string){
	s,ok:=this.GetString(key)
	if(!ok){
		panic(fmt.Errorf("%s not found in context",key))
	}
	return s
}

func (this *Context) MustGetMap(key string)(map[string]interface{}){
	s,ok:=this.GetMap(key)
	if(!ok){
		panic(fmt.Errorf("%s not found in context",key))
	}
	return s
}

func (this *Context) GetMap(key string)(map[string]interface{},bool){
	v:=this.Get(key)
	s,ok:=v.(map[string]interface{})
	if(ok){
		return s,ok
	}
	return s,ok
}

func (this *Context) Get(key string)interface{}{
	return this.Data[key]
}


func (this *Context) Set(key string,value interface{}){
	this.Data[key]=value
}

type JointInterface interface {
	Process(s *Context) (*Context, error)
}

type Pipeline struct {
	joints  []JointInterface
	context *Context
}

func (this *Pipeline) Context(s *Context) *Pipeline {
	this.context = s
	return this
}

func (this *Pipeline) Start(s JointInterface) *Pipeline {
	if(this.context ==nil){
		this.context =&Context{}
	}
	if(this.context.Data==nil){
		this.context.Data =map[string]interface{}{}
	}
	this.joints = []JointInterface{s}
	return this
}

func (this *Pipeline) Join(s JointInterface) *Pipeline {
	this.joints = append(this.joints, s)
	return this
}

func (this *Pipeline) End() *Pipeline {
	return this
}

func (this *Pipeline) Run()(*Context) {
	var err error
	for _, v := range this.joints {
		this.context, err = v.Process(this.context)
		if err != nil {
			panic(err)
		}
	}
	return this.context
}
