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

package api

import (
	"bytes"
	"fmt"
	"strconv"
)

// GetPagination return a pagination html code snippet
func GetPagination(from, size, total int, url string, param map[string]interface{}) string {

	if total > 10000 {
		total = 10000
	}

	if total <= size {
		return ""
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

	// init args start
	var moreArgs bytes.Buffer
	moreArgs.WriteString("var args='")
	if len(param) > 0 {
		for k, v := range param {
			hostStr := fmt.Sprintf("&%s=%v", k, v)
			moreArgs.WriteString(hostStr)
		}
	}

	moreArgs.WriteString("';")

	if moreArgs.Len() > 0 {
		buffer.Write(moreArgs.Bytes())
	}

	buffer.WriteString("var size=")
	buffer.WriteString(strconv.Itoa(size))
	buffer.WriteString(";")

	//init args end

	buffer.WriteString("    $(function() {")

	buffer.WriteString("$('[data-uk-pagination]').on('select.uk.pagination', function(e, pageIndex){")

	buffer.WriteString("var from=pageIndex*size;")

	buffer.WriteString("window.location='?from='+from+'&size='+size+args")

	buffer.WriteString("});")

	buffer.WriteString("   });")

	//init para for hot key  start
	buffer.WriteString(fmt.Sprintf("var maxpage = %v;", total))
	if from > 0 && from >= size {
		buffer.WriteString(fmt.Sprintf("var prev_page='?from=%v&size='+size+args;", from-size))

	}
	if from+size < total {
		buffer.WriteString(fmt.Sprintf("var next_page='?from=%v&size='+size+args;", from+size))
	}
	//init para for hot key end

	buffer.WriteString("</script>")

	return buffer.String()
}
