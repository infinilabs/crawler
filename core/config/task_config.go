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

package config

import (
	"regexp"
)

type TaskConfig struct {

	//name of this task
	//Name string

	//store page separately,such as url with paging,ie:http://baidu.com/?pn=10 http://baidu.com/?pn=20 ,here we can set value to "pn"
	//SplitByUrlParameter string `split_by_url_parameter`

	//follow page link,and walk around
	//FollowLink bool

	//walking around pattern
	LinkUrlExtractRegexStr           string `link_extract_pattern`
	LinkUrlExtractRegex           *regexp.Regexp
	LinkUrlExtractRegexGroupIndex int `link_extract_group`
	LinkUrlMustContain            string
	LinkUrlMustNotContain         string

	//parsing url pattern,when url match this pattern,gopa will not parse urls from response of this url
	SkipPageParsePatternStr string `skip_page_parse_pattern`
	SkipPageParsePattern *regexp.Regexp

	//fetch url pattern
	FetchUrlPatternStr        string `fetch_url_pattern`
	FetchUrlPattern         *regexp.Regexp
	FetchUrlMustContain     string
	FetchUrlMustNotContain  string

	//saving pattern
	SavingUrlPatternStr     string `save_url_pattern`
	SavingUrlPattern        *regexp.Regexp
	SavingUrlMustContain    string
	SavingUrlMustNotContain string

	//Crawling within domain
	FollowSameDomain        bool `follow_same_domain`
	FollowSubDomain         bool `follow_sub_domain`

	TaskDataPath            string
	//WebDataPath  string

	//User Cookie
	Cookie                  string

	//Fetch Speed Control
	FetchDelayThreshold     int
	TaskDBFilename          string `task_db_filename`
}

func (this *TaskConfig)Init() *TaskConfig  {
	this.TaskDBFilename ="taskdb"
	return this
}

type Task struct {
	Url, Request, Response []byte
}

type RoutingParameter struct {
	Shard int
}
