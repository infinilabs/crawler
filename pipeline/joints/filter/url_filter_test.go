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

package filter

import (
	"github.com/infinitbyte/framework/core/config"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/gopa/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlFilter(t *testing.T) {

	context := &pipeline.Context{}
	task := model.Task{}
	task.Url = "http://elasticsearch.cn/"
	task.OriginalUrl = "http://elasticsearch.cn/"
	task.Depth = 1
	task.Breadth = 1
	task.Host = "elasticsearch.cn"

	snapshot := model.Snapshot{}
	snapshot.Path = "/"
	snapshot.File = "default.html"

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_OriginalUrl, task.OriginalUrl)
	context.Set(model.CONTEXT_TASK_Depth, task.Depth)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse := UrlFilterJoint{}
	parse.Process(context)
	assert.Equal(t, false, context.IsEnd())

	task = model.Task{}
	task.Url = "mailto:g"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	parse.Process(context)
	assert.Equal(t, true, context.IsExit())

	task = model.Task{}
	task.Url = "asfasdf.gif"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	parse.Process(context)
	assert.Equal(t, true, context.IsExit())

}

func TestRuleCheck(t *testing.T) {
	rule := config.Rules{}
	rule.MustNot = &config.Rule{}
	rule.MustNot.Prefix = append(rule.MustNot.Prefix, "prefix")
	rule.MustNot.Contain = append(rule.MustNot.Contain, "contain")
	rule.MustNot.Suffix = append(rule.MustNot.Suffix, "suffix")

	assert.Equal(t, false, checkRule(rule, "prefix.google.com"))
	assert.Equal(t, false, checkRule(rule, "www.contain.com"))
	assert.Equal(t, false, checkRule(rule, "www.google.suffix"))

	rule = config.Rules{}
	rule.Must = &config.Rule{}
	rule.Must.Prefix = append(rule.Must.Prefix, "prefix")
	rule.Must.Contain = append(rule.Must.Contain, "contain")
	rule.Must.Suffix = append(rule.Must.Suffix, "suffix")

	assert.Equal(t, false, checkRule(rule, "prefix.google.com"))
	assert.Equal(t, false, checkRule(rule, "www.contain.com"))
	assert.Equal(t, false, checkRule(rule, "www.google.suffix"))
	assert.Equal(t, true, checkRule(rule, "prefix.contain.suffix"))

	rule = config.Rules{}
	rule.Must = &config.Rule{}
	rule.MustNot = &config.Rule{}

	rule.Should = &config.Rule{}
	rule.Should.Prefix = append(rule.Should.Prefix, "prefix")
	rule.Should.Contain = append(rule.Should.Contain, "contain")
	rule.Should.Suffix = append(rule.Should.Suffix, "suffix")

	assert.Equal(t, true, checkRule(rule, "prefix.google.com"))
	assert.Equal(t, true, checkRule(rule, "www.contain.com"))
	assert.Equal(t, true, checkRule(rule, "www.google.suffix"))
	assert.Equal(t, true, checkRule(rule, "prefix.contain.suffix"))
	assert.Equal(t, false, checkRule(rule, "www.baidu.com"))

	rule = config.Rules{}
	rule.Must = &config.Rule{}
	rule.MustNot = &config.Rule{}
	rule.MustNot.Prefix = append(rule.MustNot.Prefix, "non-exists")

	rule.Should = &config.Rule{}
	rule.Should.Prefix = append(rule.Should.Prefix, "prefix")
	rule.Should.Contain = append(rule.Should.Contain, "contain")
	rule.Should.Suffix = append(rule.Should.Suffix, "suffix")

	assert.Equal(t, true, checkRule(rule, "prefix.google.com"))
	assert.Equal(t, true, checkRule(rule, "www.contain.com"))
	assert.Equal(t, true, checkRule(rule, "www.google.suffix"))
	assert.Equal(t, true, checkRule(rule, "prefix.contain.suffix"))
	assert.Equal(t, false, checkRule(rule, "www.baidu.com"))

	rule = config.Rules{}
	rule.Must = &config.Rule{}
	rule.Must.Contain = append(rule.Must.Contain, ".")
	rule.MustNot = &config.Rule{}
	rule.MustNot.Prefix = append(rule.MustNot.Prefix, "non-exists")

	rule.Should = &config.Rule{}
	rule.Should.Prefix = append(rule.Should.Prefix, "prefix")
	rule.Should.Contain = append(rule.Should.Contain, "contain")
	rule.Should.Suffix = append(rule.Should.Suffix, "suffix")

	assert.Equal(t, true, checkRule(rule, "prefix.google.com"))
	assert.Equal(t, true, checkRule(rule, "www.contain.com"))
	assert.Equal(t, true, checkRule(rule, "www.google.suffix"))
	assert.Equal(t, true, checkRule(rule, "prefix.contain.suffix"))
	assert.Equal(t, false, checkRule(rule, "www.baidu.com"))

}
