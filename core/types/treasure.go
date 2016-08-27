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

import ."time"

type KV struct {
	Key   string
	Value []string
}

type Treasure struct {
	Proto        string
	Domain        string // elastic.co
	Path          string // /index.html
	Headers       map[string][]string   // key:value
	Parameters    []KV   // key:value
	Meta          map[string]interface{}
	Images        []KV //images within this site, img:desc
	ExtImages     []KV //images outside, img:desc
	Links         []KV //link:desc
	Body          []byte
	StatusCode    int
	RefUrl        string //the parent url to enter this url
	Url           string //full url
	Title         string
	Size          int
	SimHash       string
	Snapshot      string //Snapshot storage info
	CreateTime    Time
	UpdateTime    Time
	LastCheckTime Time
	H1            []string
	H2            []string
	H3            []string
}
