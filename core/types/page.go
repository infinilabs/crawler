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
	Key   string `json:"key,omitempty"`
	Value []string `storm:"inline" json:"value,omitempty"`
}

type PageItem struct {
	Proto      string `storm:"index" json:"proto,omitempty"`
	Domain     string  `storm:"index" json:"domain,omitempty"`            // elastic.co
	UrlPath    string  `storm:"index" json:"path,omitempty"`            // /index.html
	Headers    map[string][]string`storm:"inline" json:"headers,omitempty"` // key:value
	Parameters []KV   `storm:"inline" json:"parameters,omitempty"`             // key:value
	Images     []KV `storm:"inline" json:"images,omitempty"`// images within this site, img:desc
	ExtImages  []KV `storm:"inline" json:"ext_images,omitempty"`// images outside, img:desc
	Links      []KV `storm:"inline" json:"links,omitempty"`// link:desc
	Body       []byte`json:"-"`
	StatusCode int`storm:"index" json:"status_code,omitempty"`
	Title      string`json:"title,omitempty"`
	Size       int`json:"size,omitempty"`
	SimHash    string`storm:"index" json:"sim_hash,omitempty"`
	H1         []string`json:"h1,omitempty"`
	H2         []string`json:"h2,omitempty"`
	H3         []string`json:"h3,omitempty"`
	Metadata      *map[string]interface{} `storm:"inline" json:"metadata,omitempty"`

}

type PageLink struct {
	Url   string `json:"url"`
	Label string `json:"label"`
}

type Seed struct {
	ID         int    `storm:"id,increment" json:"id,omitempty"`
	Url        string `storm:"index" json:"url,omitempty"`
	Reference  string `json:"reference,omitempty"`
	Depth      int    `storm:"index" json:"depth,omitempty"`
	CreateTime *Time  `storm:"index" json:"created,omitempty"`
}

type TaskStatus int

const TaskCreated  =0
const TaskFetchStarted  =1
const TaskFetchFailed  =2
const TaskFetchSuccess  =3

type Task struct {
	ID            string    `storm:"id,unique" json:"id"`
	Url           string `storm:"index" json:"url,omitempty"`
	Status        TaskStatus `storm:"index" json:"status,omitempty"`
	Seed          *Seed `storm:"inline" json:"seed,omitempty"`
	Page          *PageItem `storm:"inline" json:"page,omitempty"`
	Message       interface{} `storm:"inline" json:"message,omitempty"`
	CreateTime    *Time     `storm:"index" json:"created,omitempty"`
	UpdateTime    *Time     `storm:"index" json:"updated,omitempty"`
	LastCheckTime *Time     `storm:"index" json:"checked,omitempty"`
	Snapshot      string    `json:"snapshot,omitempty"` //Snapshot storage info
}

func (this Seed) Get(url string) Seed {
	task := Seed{}
	task.Url = url
	task.Reference = ""
	task.Depth = 0
	return task
}

func (this Seed) MustGetBytes() []byte {

	bytes, err := this.GetBytes()
	if err != nil {
		panic(err)
	}
	return bytes
}

var delimiter = "|#|"

func (this Seed) GetBytes() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprint(this.Depth))
	buf.WriteString(delimiter)
	buf.WriteString(this.Reference)
	buf.WriteString(delimiter)
	buf.WriteString(this.Url)

	return buf.Bytes(), nil
}

func TaskSeedFromBytes(b []byte) Seed {
	task, err := fromBytes(b)
	if err != nil {
		panic(err)
	}
	return task
}

func fromBytes(b []byte) (Seed, error) {

	str := string(b)
	array := strings.Split(str, delimiter)
	task := Seed{}
	i, _ := strconv.Atoi(array[0])
	task.Depth = i
	task.Reference = array[1]
	task.Url = array[2]

	return task, nil
}

func NewTaskSeed(url, ref string, depth int) Seed {
	task := Seed{}
	task.Url = url
	task.Reference = ref
	task.Depth = depth
	return task
}
