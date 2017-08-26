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

package pipe

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/model"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/util"
	"regexp"
	"strings"
	"sync"
)

type HtmlToTextJoint struct {
	Parameters
}

//merge whitespace and \n
const mergeWhitespace ParaKey = "merge_whitespace"

func (joint HtmlToTextJoint) Name() string {
	return "html2text"
}

type cleanRule struct {
	l                sync.RWMutex
	replaceRules     []*regexp.Regexp
	inited           bool
	lowerCaseRule    *regexp.Regexp
	removeTagsRule   *regexp.Regexp
	removeBreaksRule *regexp.Regexp
}

var rules = cleanRule{replaceRules: []*regexp.Regexp{}}

func getRule(str string) *regexp.Regexp {
	re, _ := regexp.Compile(str)
	return re
}

func initRules() {
	rules.l.Lock()
	defer rules.l.Unlock()
	if rules.inited {
		return
	}

	log.Trace("init html2text rule")

	//remove STYLE
	rules.replaceRules = append(rules.replaceRules, getRule(`<style[\S\s]+?\</style\>`))

	//remove META
	rules.replaceRules = append(rules.replaceRules, getRule(`\<meta[\S\s]+?\</meta\>`))

	//remove comments
	rules.replaceRules = append(rules.replaceRules, getRule(`<!--[\S\s]*?-->`))

	//remove SCRIPT,NOSCRIPT
	rules.replaceRules = append(rules.replaceRules, getRule(`\<script[\S\s]+?.*?\</script\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<noscript[\S\s]+?\</noscript\>`))

	//remove iframe,frame
	rules.replaceRules = append(rules.replaceRules, getRule(`\<iframe[\S\s]+?\</iframe\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<frame[\S\s]+?\</frame\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<frameset[\S\s]+?\</frameset\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<noframes[\S\s]+?\</noframes\>`))

	//remove embed objects
	rules.replaceRules = append(rules.replaceRules, getRule(`\<noembed[\S\s]+?\</noembed\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<embed[\S\s]+?\</embed\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<applet[\S\s]+?\</applet\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<object[\S\s]+?\</object\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<base[\S\s]+?\</base\>`))

	//remove code blocks
	rules.replaceRules = append(rules.replaceRules, getRule(`\<pre[\S\s]+?\</pre\>`))
	rules.replaceRules = append(rules.replaceRules, getRule(`\<code[\S\s]+?\</code\>`))

	//lowercase html tags
	rules.lowerCaseRule, _ = regexp.Compile("\\<[\\S\\s]+?\\>")

	//remove all HTML tags and replaced with \n
	rules.removeTagsRule, _ = regexp.Compile("\\<[\\S\\s]+?\\>")

	//remove continued break lines
	rules.removeBreaksRule, _ = regexp.Compile("\\s{2,}")

	rules.inited = true

}

func replaceAll(src []byte) []byte {
	initRules()
	empty := []byte(" ")

	str := string(src)
	if rules.lowerCaseRule != nil {
		str = rules.lowerCaseRule.ReplaceAllStringFunc(str, strings.ToLower)
		src = []byte(src)
	}

	if rules.replaceRules != nil {
		for _, rule := range rules.replaceRules {
			src = rule.ReplaceAll(src, empty)
		}
	}

	if rules.removeTagsRule != nil {
		src = rules.removeTagsRule.ReplaceAll(src, []byte("\n"))
	}

	if rules.removeBreaksRule != nil {
		src = rules.removeBreaksRule.ReplaceAll(src, []byte("\n"))
	}

	return src
}

func (joint HtmlToTextJoint) Process(context *Context) error {

	snapshot := context.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	body := replaceAll(snapshot.Payload)

	src := string(body)

	if joint.GetBool(mergeWhitespace, false) {
		src = util.MergeSpace(src)
	}

	src = strings.Replace(src, "&#8216;", "'", -1)
	src = strings.Replace(src, "&#8217;", "'", -1)
	src = strings.Replace(src, "&#8220;", "\"", -1)
	src = strings.Replace(src, "&#8221;", "\"", -1)
	src = strings.Replace(src, "&nbsp;", " ", -1)
	src = strings.Replace(src, "&quot;", "\"", -1)
	src = strings.Replace(src, "&apos;", "'", -1)
	src = strings.Replace(src, "&#34;", "\"", -1)
	src = strings.Replace(src, "&#39;", "'", -1)
	src = strings.Replace(src, "&amp; ", "& ", -1)
	src = strings.Replace(src, "&amp;amp; ", "& ", -1)

	log.Trace("get text: ", src)

	snapshot.Text = util.XSSHandle(src)
	return nil
}
