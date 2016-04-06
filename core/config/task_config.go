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
	Name string

	//store page separately,such as url with paging,ie:http://baidu.com/?pn=10 http://baidu.com/?pn=20 ,here we can set value to "pn"
	SplitByUrlParameter string

	//follow page link,and walk around
	FollowLink bool

	//walking around pattern
	LinkUrlExtractRegex           *regexp.Regexp
	LinkUrlExtractRegexGroupIndex int
	LinkUrlMustContain            string
	LinkUrlMustNotContain         string

	//parsing url pattern,when url match this pattern,gopa will not parse urls from response of this url
	SkipPageParsePattern *regexp.Regexp

	//fetch url pattern
	FetchUrlPattern        *regexp.Regexp
	FetchUrlMustContain    string
	FetchUrlMustNotContain string

	//saving pattern
	SavingUrlPattern        *regexp.Regexp
	SavingUrlMustContain    string
	SavingUrlMustNotContain string

	//Crawling within domain
	FollowSameDomain bool
	FollowSubDomain  bool

	TaskDataPath string
	WebDataPath  string

	//User Cookie
	Cookie string

	//Fetch Speed Control
	FetchDelayThreshold int
}

type Task struct {
	Url, Request, Response []byte
}

type RoutingParameter struct {
	Shard int
}
