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
	"github.com/infinitbyte/gopa/core/env"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

type crawlerJoint struct {
	Parameters
}

func (joint crawlerJoint) Name() string {
	return "crawler"
}

func (joint crawlerJoint) Process(s *Context) error {
	s.Data[("webpage")] = "hello world gogo "
	s.Data["received_url"] = joint.Data["url"]
	s.Data[("status")] = true
	fmt.Println("start to crawlling url: ", joint.Get("url"))
	return nil
}

type parserJoint struct {
}

func (joint parserJoint) Name() string {
	return "parser"
}

func (joint parserJoint) Process(s *Context) error {
	s.Data[("urls")] = "gogo"
	s.Data[("domain")] = "http://gogo.com"
	//pub urls to channel
	fmt.Println("start to parse web content")
	return nil
}

type saveJoint struct {
}

func (joint saveJoint) Name() string {
	return "save"
}

func (joint saveJoint) Process(s *Context) error {
	s.Set("saved", "true")
	//pub urls to channel
	fmt.Println("start to save web content")
	return nil
}

type publishJoint struct {
}

func (joint publishJoint) Name() string {
	return "publish"
}

func (joint publishJoint) Process(s *Context) error {
	fmt.Println("start to end pipeline")
	s.Set("published", "true")
	return nil
}

func TestPipeline(t *testing.T) {

	global.RegisterEnv(env.EmptyEnv())

	pipeline := NewPipeline("crawler_test")
	context := &Context{}
	context.Init()
	context.Data["url"] = "gogol.com"
	context.Data["webpage"] = "hello world gogo "

	crawler := crawlerJoint{}

	pipeline.Context(context).
		Start(crawler).
		Join(parserJoint{}).
		Join(saveJoint{}).
		Join(publishJoint{}).
		Run()

	fmt.Println(context.Parameters.Data)
	fmt.Println(context.Data)
	assert.Equal(t, "true", context.Parameters.Data["published"])
	assert.Equal(t, "true", context.Parameters.Data["saved"])
	assert.Equal(t, true, context.Parameters.Data["status"])
	assert.Equal(t, "http://gogo.com", context.Parameters.Data["domain"])
}

const key1 ParaKey = "DEPTH"
const key2 ParaKey = "DEPTH2"

func TestContext(t *testing.T) {
	global.RegisterEnv(env.EmptyEnv())
	context := &Context{}
	context.Parameters.Init()
	context.Parameters.Set(key1, 23)
	fmt.Println(util.ToJson(context, true))
	v := context.MustGetInt(key1)
	assert.Equal(t, 23, v)
	v, _ = context.GetInt(key2, 0)
	assert.Equal(t, 0, v)

}
