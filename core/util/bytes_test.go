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
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToLowercase(t *testing.T) {
	str := []byte("AZazUPPERcase")

	printStr(str)
	ToLowercase(str)
	fmt.Println("lowercased:")
	assert.Equal(t, "azazuppercase", string(str))
	printStr(str)
	ToUppercase(str)
	fmt.Println("uppercased:")
	assert.Equal(t, "AZAZUPPERCASE", string(str))
	printStr(str)
}

func printStr(str []byte) {
	for i, s := range str {
		fmt.Println(i, "-", s, "-", string(s))
	}
}
