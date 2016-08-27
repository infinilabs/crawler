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

type Stream struct {
	data map[string]interface{}
}

type JointInterface interface {
	Process(s *Stream) (*Stream, error)
}

type Pipeline struct {
	joints []JointInterface
	stream *Stream
}

func (this *Pipeline) Input(s *Stream) *Pipeline {
	this.stream = s
	return this
}

func (this *Pipeline) Start(s JointInterface) *Pipeline {
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

func (this *Pipeline) Run()(*Stream) {
	var err error
	for _, v := range this.joints {
		this.stream, err = v.Process(this.stream)
		if err != nil {
			panic(err)
		}
	}
	return this.stream
}
