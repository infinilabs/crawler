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

package joint

import (
	"github.com/infinitbyte/framework/core/config"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/model"
)

// UrlFilterJoint used to validate urls, include host,path,file and file extension
type UrlFilterJoint struct {
	pipeline.Parameters
}

// Name is url_filter
func (joint UrlFilterJoint) Name() string {
	return "url_filter"
}

var urlMatchRule pipeline.ParaKey = "url_match_rule"
var hostMatchRule pipeline.ParaKey = "host_match_rule"
var pathMatchRule pipeline.ParaKey = "path_match_rule"
var fileMatchRule pipeline.ParaKey = "file_match_rule"
var fileExtensionMatchRule pipeline.ParaKey = "file_ext_match_rule"

// Process check all the url match rules
func (joint UrlFilterJoint) Process(context *pipeline.Context) error {

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	originalUrl := context.GetStringOrDefault(model.CONTEXT_TASK_OriginalUrl, "")
	url := context.MustGetString(model.CONTEXT_TASK_URL)
	host := context.MustGetString(model.CONTEXT_TASK_Host)
	if url == "" {
		context.Exit("nil url")
		return nil
	}

	if originalUrl != "" {
		if !joint.validRule(urlMatchRule, originalUrl) {
			context.Exit("invalid url (original), " + originalUrl)
			return nil
		}
	}

	if !joint.validRule(urlMatchRule, url) {
		context.Exit("invalid url, " + url)
		return nil
	}

	if !joint.validRule(hostMatchRule, host) {
		context.Exit("invalid host, " + host)
		return nil
	}

	if !joint.validRule(pathMatchRule, snapshot.Path) {
		context.Exit("invalid path, " + snapshot.Path)
		return nil
	}

	if !joint.validRule(fileMatchRule, snapshot.File) {
		context.Exit("invalid file, " + snapshot.File)
		return nil
	}

	if !joint.validRule(fileExtensionMatchRule, util.FileExtension(snapshot.File)) {
		context.Exit("invalid file extension, " + snapshot.File)
		return nil
	}

	return nil
}

func getDefaultUrlMatchConfig() config.Rules {
	rule := config.Rules{}
	rule.MustNot = &config.Rule{}
	rule.MustNot.Prefix = append(rule.MustNot.Prefix, "mailto:")
	rule.MustNot.Prefix = append(rule.MustNot.Prefix, "data:image/")
	rule.MustNot.Prefix = append(rule.MustNot.Prefix, "#")
	rule.MustNot.Prefix = append(rule.MustNot.Prefix, "javascript:")
	return rule
}

// validRule check if the url are satisfy the rule, default is true
func (joint UrlFilterJoint) validRule(key pipeline.ParaKey, target string) bool {

	if target == "" {
		return true
	}

	rule, ok := joint.GetMap(key)

	matchRule := config.Rules{}

	if !ok {
		if key == urlMatchRule {
			matchRule = getDefaultUrlMatchConfig()
		} else {
			return true
		}
	}

	cfg, err := config.NewConfigFrom(rule)
	if err != nil {
		panic(err)
	}
	cfg.Unpack(&matchRule)

	return checkRule(matchRule, target)
}

func checkRule(matchRule config.Rules, target string) bool {

	result := true
	if matchRule.Must != nil {
		for _, item := range matchRule.Must.Contain {
			if !(util.ContainStr(target, item)) {
				return false
			}
		}

		for _, item := range matchRule.Must.Prefix {
			if !(util.PrefixStr(target, item)) {
				return false
			}
		}
		for _, item := range matchRule.Must.Suffix {
			if !(util.SuffixStr(target, item)) {
				return false
			}
		}

		if len(matchRule.Must.Prefix) > 0 || len(matchRule.Must.Contain) > 0 || len(matchRule.Must.Suffix) > 0 {
			result = true
		}
	}

	if matchRule.MustNot != nil {
		for _, item := range matchRule.MustNot.Contain {
			if util.ContainStr(target, item) {
				return false
			}
		}

		for _, item := range matchRule.MustNot.Prefix {
			if util.PrefixStr(target, item) {
				return false
			}
		}
		for _, item := range matchRule.MustNot.Suffix {
			if util.SuffixStr(target, item) {
				return false
			}
		}
		if len(matchRule.MustNot.Prefix) > 0 || len(matchRule.MustNot.Contain) > 0 || len(matchRule.MustNot.Suffix) > 0 {
			result = true
		}
	}

	if matchRule.Should != nil {
		for _, item := range matchRule.Should.Contain {
			if util.ContainStr(target, item) {
				return true
			}
		}

		for _, item := range matchRule.Should.Prefix {
			if util.PrefixStr(target, item) {
				return true
			}
		}
		for _, item := range matchRule.Should.Suffix {
			if util.SuffixStr(target, item) {
				return true
			}
		}
		if len(matchRule.Should.Prefix) > 0 || len(matchRule.Should.Contain) > 0 || len(matchRule.Should.Suffix) > 0 {
			result = false
		}
	}
	return result
}
