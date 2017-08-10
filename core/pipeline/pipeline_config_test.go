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
	"github.com/infinitbyte/gopa/core/env"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPipelineConfig(t *testing.T) {

	util.SetIDPersistencePath("/tmp")

	global.RegisterEnv(env.EmptyEnv())
	global.Env().IsDebug = true

	config := PipelineConfig{}
	config.Name = "test_pipe_line"

	config.Context = &Context{}
	config.Context.Init()
	config.Context.Data["url"] = "gogol.com"
	config.Context.Data["webpage"] = "hello world gogo "

	Register(crawlerJoint{})
	Register(parserJoint{})
	Register(saveJoint{})
	Register(publishJoint{})

	config.StartJoint = &JointConfig{Enabled: true, JointName: "crawler", Parameters: map[string]interface{}{"url": "http://baidu12.com"}}
	joints := []*JointConfig{}
	joints = append(joints, &JointConfig{Enabled: true, JointName: "parser", Parameters: map[string]interface{}{}})
	joints = append(joints, &JointConfig{Enabled: true, JointName: "save", Parameters: map[string]interface{}{}})
	joints = append(joints, &JointConfig{Enabled: true, JointName: "publish", Parameters: map[string]interface{}{}})

	config.ProcessJoints = joints

	//fmt.Println(util.ToJson(config,true))

	pipe := NewPipelineFromConfig(&config)
	context := pipe.Run()

	fmt.Println("pipeline context")
	fmt.Println(context)
	fmt.Println(config.Context)
	fmt.Println(context.Data["received_url"])

	assert.Equal(t, "http://baidu12.com", context.Data["received_url"])
	assert.Equal(t, "true", context.Data["published"])
	assert.Equal(t, "true", context.Data["saved"])
	assert.Equal(t, true, context.Data["status"])
	assert.Equal(t, "http://gogo.com", context.Data["domain"])

}

type do interface {
	Do() string
}

type base struct {
	Para map[string]interface{}
}

type foo struct {
	base
	Id   int
	Name string
}

func (joint foo) Do() string {
	fmt.Println("foo do,", joint.Id, ",", joint.Name, ",", joint.Para)
	return joint.Name
}

func (joint bar) Do() string {
	fmt.Println("foo do")
	return ""
}

type bar struct {
}

func TestPipelineConfigReflection(t3 *testing.T) {
	var regStruct map[string]interface{}
	regStruct = make(map[string]interface{})
	regStruct["Foo"] = foo{Id: 1, Name: "medcl"}
	regStruct["Bar"] = bar{}

	str := "Bar"
	if regStruct[str] != nil {
		t := reflect.ValueOf(regStruct[str]).Type()
		v := reflect.New(t).Elem()
		fmt.Println(v)
		//v.MethodByName("Do").Call(nil)
	}

	//get another instance again
	str = "Foo"
	if regStruct[str] != nil {
		t := reflect.ValueOf(regStruct[str]).Type()
		v := reflect.New(t).Elem()
		fmt.Println(v)
		v1 := v.Interface().(do)
		v1.Do()
		assert.Equal(t3, "", v1.Do())
	}

	str = "Foo"
	if regStruct[str] != nil {
		t := reflect.ValueOf(regStruct[str]).Type()
		v := reflect.New(t).Elem()
		fmt.Println(v)

		f := v.FieldByName("Name")
		if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
			f.SetString("tom")
		}

		f = v.FieldByName("Id")
		if f.IsValid() && f.CanSet() && f.Kind() == reflect.Int {
			f.SetInt(55)
		}
		f = v.FieldByName("Para")
		fmt.Println(f.Kind())
		if f.IsValid() && f.CanSet() && f.Kind() == reflect.Map {
			para := map[string]interface{}{}
			para["key"] = "value123"
			f.Set(reflect.ValueOf(para))
		}

		fmt.Println(v)
		v1 := v.Interface().(do)
		v1.Do()
		assert.Equal(t3, "tom", v1.Do())

	}

	//get another instance again
	str = "Foo"
	if regStruct[str] != nil {
		t := reflect.ValueOf(regStruct[str]).Type()
		v := reflect.New(t).Elem()
		fmt.Println(v)
		v1 := v.Interface().(do)
		v1.Do()
		assert.Equal(t3, "", v1.Do())

	}

}
