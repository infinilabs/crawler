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
)

type HtmlToTextJoint struct {
	Parameters
}

const mergeWhitespace ParaKey = "merge_whitespace" //merge whitespace and \n

func (joint HtmlToTextJoint) Name() string {
	return "html2text"
}

func (joint HtmlToTextJoint) Process(context *Context) error {

	//TODO all configable
	snapshot := context.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	body := snapshot.Payload
	src := string(body)
	//lowercase html tags
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//remove STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//remove META
	re, _ = regexp.Compile("\\<meta[\\S\\s]+?\\</meta\\>")
	src = re.ReplaceAllString(src, "")

	//remove comments
	re, _ = regexp.Compile("<!--[\\S\\s]*?-->")
	src = re.ReplaceAllString(src, "")

	//remove SCRIPT,NOSCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<noscript[\\S\\s]+?\\</noscript\\>")
	src = re.ReplaceAllString(src, "")

	//remove iframe,frame
	re, _ = regexp.Compile("\\<iframe[\\S\\s]+?\\</iframe\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<frame[\\S\\s]+?\\</frame\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<frameset[\\S\\s]+?\\</frameset\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<noframes[\\S\\s]+?\\</noframes\\>")
	src = re.ReplaceAllString(src, "")

	//remove embed objects
	re, _ = regexp.Compile("\\<noembed[\\S\\s]+?\\</noembed\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<embed[\\S\\s]+?\\</embed\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<applet[\\S\\s]+?\\</applet\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<object[\\S\\s]+?\\</object\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<base[\\S\\s]+?\\</base\\>")
	src = re.ReplaceAllString(src, "")

	//remove code blocks
	re, _ = regexp.Compile("\\<pre[\\S\\s]+?\\</pre\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<code[\\S\\s]+?\\</code\\>")
	src = re.ReplaceAllString(src, "")

	//remove all HTML tags and replaced with \n
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//remove continued break lines
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

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
