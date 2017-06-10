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
	"encoding/json"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	data, _ := get("http://es-guide-preview.elasticsearch.cn", "", "")

	data1, _ := json.Marshal(data)
	fmt.Println("", string(data1))
	//assert.Equal(t,data.StatusCode,301)
}

func TestGetHost(t *testing.T) {

	url := "/index.html"
	host := GetHost(url)
	fmt.Println("", host)

	url = "www.baidu.com/index.html"
	host = GetHost(url)
	fmt.Println("www.baidu.com", host)

	url = "//www.baidu.com/index.html"
	host = GetHost(url)
	fmt.Println("www.baidu.com", host)

	url = "http://www.baidu.com/index.html"
	host = GetHost(url)
	fmt.Println("www.baidu.com", host)

	url = "https://www.baidu.com/index.html"
	host = GetHost(url)
	fmt.Println("www.baidu.com", host)
}
