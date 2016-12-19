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

package types

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPageTask_GetBytes(t *testing.T) {
	task:= TaskSeed{Url:"baidu.com",Reference:"google.com",Depth:5}
	bytes,_:=task.GetBytes()
	task1,_:=fromBytes(bytes)

	assert.Equal(t,"baidu.com",task1.Url)
	assert.Equal(t,"google.com",task1.Reference)
	assert.Equal(t,5,task1.Depth)

}

