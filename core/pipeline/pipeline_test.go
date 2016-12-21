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
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)


type crawlerJoint struct {
	Url string
}

func (this crawlerJoint) Process(s *Context) (*Context, error) {
	s.data = map[ContextKey]interface{}{}
	s.data[("webpage")] = "hello world gogo "
	s.data[("status")] = true
	fmt.Println("start to crawlling url:"+this.Url)
	return s, nil
}

type parserJoint struct {
}

func (this parserJoint) Process(s *Context) (*Context, error) {
	s.data[("urls")] = "gogo"
	s.data[("domain")] = "http://gogo.com"
	//pub urls to channel
	fmt.Println("start to parse web content")
	return s, nil
}

type saveJoint struct {
}

func (this saveJoint) Process(s *Context) (*Context, error) {
	s.Set("saved","true")
	//pub urls to channel
	fmt.Println("start to save web content")
	return s, nil
}

type publishJoint struct {
}

func (this publishJoint) Process(s *Context) (*Context, error) {
	fmt.Println("start to end pipeline")
	s.Set("published","true")
	return s, nil
}


func TestPipeline(t *testing.T)  {

	pipeline:=NewPipeline("crawler_test")
	stream:=&Context{}
	stream.data =map[ContextKey]interface{}{}
	stream.data["url"]="gogol.com"
	stream.data["webpage"]="hello world gogo "

	stream= pipeline.Context(stream).
		Start(crawlerJoint{Url:"http://baidu.com"}).
		Join(parserJoint{}).
		Join(saveJoint{}).
		Join(publishJoint{}).
		End().
		Run()

	fmt.Println(stream.data)
	assert.Equal(t,stream.data["saved"],"true")
	assert.Equal(t,stream.data["status"],true)
	assert.Equal(t,stream.data["domain"],"http://gogo.com")
}
