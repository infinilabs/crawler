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
	."github.com/medcl/gopa/core/pipeline"
	"strings"
	"regexp"
	"github.com/medcl/gopa/core/util"
	log "github.com/cihub/seelog"
)

type HtmlToTextJoint struct {
	MergeWhitespace bool
}

func (this HtmlToTextJoint) Name() string {
	return "html2text"
}

func (this HtmlToTextJoint) Process(context *Context) (*Context, error) {

	//TODO all configable
	body:=context.MustGetBytes(CONTEXT_PAGE_BODY_BYTES)
	src := string(body)
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除META
	re, _ = regexp.Compile("\\<meta[\\S\\s]+?\\</meta\\>")
	src = re.ReplaceAllString(src, "")

	//去除注释
	re, _ = regexp.Compile("<!--[\\S\\s]*?-->")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT,NOSCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("\\<noscript[\\S\\s]+?\\</noscript\\>")
	src = re.ReplaceAllString(src, "")

	//去除iframe,frame
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

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	if(this.MergeWhitespace){
		src = strings.TrimSpace(util.MergeSpace(src))
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

	log.Trace("get text: ",src)

	context.Set(CONTEXT_PAGE_BODY_PLAIN_TEXT,src)
	return context, nil
}
