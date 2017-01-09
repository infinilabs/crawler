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
	"github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/global"
	"github.com/stretchr/testify/assert"
	"testing"
)


type crawlerJoint struct {
	Parameters
}

func (this crawlerJoint) Name() string {
	return "crawlerJoint"
}

func (this crawlerJoint) Process(s *Context) (*Context, error) {
	s.Data[("webpage")] = "hello world gogo "
	s.Data["received_url"] = this.Data["url"]
	s.Data[("status")] = true
	fmt.Println("start to crawlling url: ",this.Get("url"))// + this.GetParameter("url").(string))
	return s, nil
}

type parserJoint struct {

}

func (this parserJoint) Name() string {
	return "parserJoint"
}

func (this parserJoint) Process(s *Context) (*Context, error) {
	s.Parameters.Data[("urls")] = "gogo"
	s.Parameters.Data[("domain")] = "http://gogo.com"
	//pub urls to channel
	fmt.Println("start to parse web content")
	return s, nil
}

type saveJoint struct {
}

func (this saveJoint) Name() string {
	return "saveJoint"
}

func (this saveJoint) Process(s *Context) (*Context, error) {
	s.Parameters.Set("saved", "true")
	//pub urls to channel
	fmt.Println("start to save web content")
	return s, nil
}

type publishJoint struct {

}

func (this publishJoint) Name() string {
	return "publishJoint"
}

func (this publishJoint) Process(s *Context) (*Context, error) {
	fmt.Println("start to end pipeline")
	s.Parameters.Set("published", "true")
	return s, nil
}

func TestPipeline(t *testing.T) {

	global.RegisterEnv(env.EmptyEnv())

	pipeline := NewPipeline("crawler_test")
	context := &Context{}
	context.Parameters.Init()
	context.Parameters.Data["url"] = "gogol.com"
	context.Parameters.Data["webpage"] = "hello world gogo "

	crawler:=crawlerJoint{}

	context = pipeline.Context(context).
		Start(crawler).
		Join(parserJoint{}).
		Join(saveJoint{}).
		Join(publishJoint{}).
		Run()

	fmt.Println(context.Data)
	assert.Equal(t, context.Parameters.Data["published"], "true")
	assert.Equal(t, context.Parameters.Data["saved"], "true")
	assert.Equal(t, context.Parameters.Data["status"], true)
	assert.Equal(t, context.Parameters.Data["domain"], "http://gogo.com")
}
