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
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/model"
	"regexp"
	"strings"
	"sync"
)

type HtmlToTextJoint struct {
	pipeline.Parameters
}

//merge whitespace and \n
const mergeWhitespace pipeline.ParaKey = "merge_whitespace"
const removeNonScript pipeline.ParaKey = "remove_nonscript"

func (joint HtmlToTextJoint) Name() string {
	return "html2text"
}

type cleanRule struct {
	l                   sync.RWMutex
	replaceRules        []*regexp.Regexp
	inited              bool
	removeTagsRule      *regexp.Regexp
	removeBreaksRule    *regexp.Regexp
	removeNonScriptRule *regexp.Regexp
	lowercase           bool
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

	//remove SCRIPT
	rules.replaceRules = append(rules.replaceRules, getRule(`\<script[\S\s]+?.*?\</script\>`))

	//remove NOSCRIPT
	rules.removeNonScriptRule = getRule(`\<noscript[\S\s]+?\</noscript\>`)

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

	//remove all HTML tags and replaced with \n
	rules.removeTagsRule, _ = regexp.Compile("\\<[\\S\\s]+?\\>")

	//remove continued break lines
	rules.removeBreaksRule, _ = regexp.Compile("\\s{2,}")

	//lowercase all the text
	rules.lowercase = true

	rules.inited = true

}

// should equal to regex("\\<[\\S\\s]+?\\>").ReplaceAllStringFunc(str, strings.ToLower)
func lowercaseTag(str []byte) {

	startLowercase := false
	startLowercaseIndex := -1
	endLowercase := false
	endLowercaseIndex := -1

	for i, s := range str {
		if s == 60 {
			startLowercase = true
			startLowercaseIndex = i
		}
		if s == 62 {
			endLowercase = true
			endLowercaseIndex = i
		}
		if startLowercase && endLowercase && endLowercaseIndex > startLowercaseIndex {
			for j := startLowercaseIndex; j < endLowercaseIndex; j++ {
				x := str[j]
				if x > 64 && x < 91 {
					str[j] = x + 32
				}
			}
			startLowercase = false
			endLowercase = false
			startLowercaseIndex = -1
			endLowercaseIndex = -1
		}
	}
}

var empty = []byte(" ")

func replaceAll(src []byte) []byte {

	if rules.lowercase {
		lowercaseTag(src)
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

func (joint HtmlToTextJoint) Process(context *pipeline.Context) error {
	initRules()

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	if !util.PrefixStr(snapshot.ContentType, "text/") {
		log.Debugf("snapshot is not text, %s, %s , %s", snapshot.ID, snapshot.Url, snapshot.ContentType)
		return nil
	}

	body := snapshot.Payload
	if joint.GetBool(removeNonScript, true) {
		body = rules.removeNonScriptRule.ReplaceAll(body, []byte(""))
	}

	body = replaceAll(body)

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

	snapshot.Text = util.XSSHandle(src)

	if global.Env().IsDebug {
		log.Trace("get text: ", src)
	}

	return nil
}
