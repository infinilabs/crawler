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
	"bytes"
	"fmt"
	"strconv"
	"strings"
	. "time"
)

type KV struct {
	Key   string
	Value []string
}

type PageItem struct {
	Proto      string
	Domain     string              // elastic.co
	UrlPath    string              // /index.html
	Headers    map[string][]string // key:value
	Parameters []KV                // key:value
	Meta       map[string]interface{}
	Images     []KV // images within this site, img:desc
	ExtImages  []KV // images outside, img:desc
	Links      []KV // link:desc
	Body       []byte
	StatusCode int
	Title      string
	Size       int
	SimHash    string
	H1         []string
	H2         []string
	H3         []string
}

type PageLink struct {
	Url   string `json:"url"`
	Label string `json:"label"`
}

type TaskSeed struct {
	ID         int    `storm:"id,increment" json:"id,omitempty"`
	Url        string `storm:"index" json:"url,omitempty"`
	Reference  string `json:"reference,omitempty"`
	Depth      int    `storm:"index" json:"depth,omitempty"`
	CreateTime *Time  `storm:"index" json:"created,omitempty"`
}

type CrawlerTask struct {
	ID            string    `storm:"id,unique" json:"id"`
	Seed          *TaskSeed `storm:"inline" json:"seed,omitempty"`
	Page          *PageItem `storm:"inline" json:"page,omitempty"`
	CreateTime    *Time     `storm:"index" json:"created,omitempty"`
	UpdateTime    *Time     `storm:"index" json:"updated,omitempty"`
	LastCheckTime *Time     `storm:"index" json:"checked,omitempty"`
	Snapshot      string    `json:"snapshot,omitempty"` //Snapshot storage info
}

func (this TaskSeed) Get(url string) TaskSeed {
	task := TaskSeed{}
	task.Url = url
	task.Reference = ""
	task.Depth = 0
	return task
}

func (this TaskSeed) MustGetBytes() []byte {

	bytes, err := this.GetBytes()
	if err != nil {
		panic(err)
	}
	return bytes
}

var delimiter = "|#|"

func (this TaskSeed) GetBytes() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprint(this.Depth))
	buf.WriteString(delimiter)
	buf.WriteString(this.Reference)
	buf.WriteString(delimiter)
	buf.WriteString(this.Url)

	return buf.Bytes(), nil
}

func PageTaskFromBytes(b []byte) TaskSeed {
	task, err := fromBytes(b)
	if err != nil {
		panic(err)
	}
	return task
}

func fromBytes(b []byte) (TaskSeed, error) {

	str := string(b)
	array := strings.Split(str, delimiter)
	task := TaskSeed{}
	i, _ := strconv.Atoi(array[0])
	task.Depth = i
	task.Reference = array[1]
	task.Url = array[2]

	return task, nil
}

func NewPageTask(url, ref string, depth int) TaskSeed {
	task := TaskSeed{}
	task.Url = url
	task.Reference = ref
	task.Depth = depth
	return task
}
