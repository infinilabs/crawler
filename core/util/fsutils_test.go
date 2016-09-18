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
	"testing"
	"fmt"
	"path"
	"github.com/stretchr/testify/assert"
)

func TestJoinPath(t *testing.T) {
	path1:="wwww.baidu.com"
	path2:="/blog/"
	path3:="/comments/1.html"
	str:=path.Join(path1,path2,path3)
	fmt.Println(str)
	assert.Equal(t,"wwww.baidu.com/blog/comments/1.html",str)
}

