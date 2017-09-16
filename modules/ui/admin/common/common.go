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

package common

import (
	"bytes"
	"strconv"
)

// NavCurrent return the current nav html code snippet
func NavCurrent(cur, nav string) string {
	if cur == nav {
		return " class=\"uk-active\" "
	}
	return ""
}

type navObj struct {
	name        string
	displayName string
	url         string
}

var navs []navObj

// RegisterNav register a custom nav link
func RegisterNav(name, displayName string, url string) {
	obj := navObj{name: name, displayName: displayName, url: url}
	navs = append(navs, obj)
}

// GetPagination return a pagination html code snippet
func GetPagination(domain string, from, size, total int, url string) string {

	if total > 10000 {
		total = 10000
	}

	var cur = from / size

	var buffer bytes.Buffer
	buffer.WriteString("<ul class=\"uk-pagination\" data-uk-pagination=\"{items:")
	buffer.WriteString(strconv.Itoa(total))
	buffer.WriteString(", itemsOnPage:")
	buffer.WriteString(strconv.Itoa(size))
	buffer.WriteString(",currentPage:")
	buffer.WriteString(strconv.Itoa(cur))
	buffer.WriteString("}\"></ul>")
	buffer.WriteString("<script type=\"text/javascript\">")
	buffer.WriteString("    $(function() {")

	buffer.WriteString("$('[data-uk-pagination]').on('select.uk.pagination', function(e, pageIndex){")
	buffer.WriteString("var size=")
	buffer.WriteString(strconv.Itoa(size))
	buffer.WriteString(";")
	buffer.WriteString("var from=pageIndex*size;")

	var moreArgs bytes.Buffer
	moreArgs.WriteString("var args='")
	if domain != "" {
		var domainStr = "&domain=" + domain
		moreArgs.WriteString(domainStr)
	}
	moreArgs.WriteString("';")

	if moreArgs.Len() > 0 {
		buffer.Write(moreArgs.Bytes())
	}

	buffer.WriteString("window.location='?from='+from+'&size='+size+args")

	buffer.WriteString("});")

	buffer.WriteString("   });")
	buffer.WriteString("</script>")

	return buffer.String()
}

// GetJSBlock return a JS wrapped code block
func GetJSBlock(buffer *bytes.Buffer, js string) {

	buffer.WriteString("<script type=\"text/javascript\">")
	buffer.WriteString("    $(function() {")
	buffer.WriteString(js)
	buffer.WriteString("   });")
	buffer.WriteString("</script>")

}
