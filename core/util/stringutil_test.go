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

package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeSpacke(t *testing.T) {
	str:="hello world     !"
	str=MergeSpace(str)
	assert.Equal(t, "hello world !",str)

	str=" hello world  \n \r  !"
	str=MergeSpace(str)
	assert.Equal(t, "hello world !",str)
}
