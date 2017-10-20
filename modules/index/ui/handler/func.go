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

package handler

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

func safeGetField(v interface{}, nullValue string) string {
	if v != nil {
		return v.(string)
	}
	return nullValue
}

func smartGetField(v1 []interface{}, v2 interface{}, nullValue string) string {

	if len(v1) > 0 {
		vv1 := safeGetField(v1[0], "")
		if vv1 != "" {
			return vv1
		}
	}

	return safeGetField(v2, nullValue)
}

func getBucketLabel(k string) string {
	if strings.Contains(k, "|") {
		return strings.Split(k, "|")[1]
	}
	return k
}
func getBucketKey(k string) string {
	if strings.Contains(k, "|") {
		return strings.Split(k, "|")[0]
	}
	return k
}

func getNavBlock(w io.Writer, r *http.Request) string {
	var buffer bytes.Buffer

	buffer.WriteString("<div style='clear:both;'>")
	buffer.WriteString("<P id=page><span style='paddint:2px;'>&nbsp;1&nbsp;</span>")
	buffer.WriteString("<span style='paddint:2px;'>")
	buffer.WriteString("<A style='padding:2px;' href=?q=search&amp;tab=1&amp;&amp;pn=2>[2]</A></span>")
	buffer.WriteString("<span style='paddint:2px;'>")
	buffer.WriteString("<A style='padding:2px;' href=?q=search&amp;tab=1&amp;&amp;pn=3>[3]</A>")
	buffer.WriteString("</span>")

	buffer.WriteString("<A class=n href=?q=search&amp;tab=1&amp;&amp;pn=2>Next</A></P>")
	buffer.WriteString("<script type=text/javascript> var maxpage = 95;var next_page='?q=<%==q%>&pn=2';")
	buffer.WriteString("</script>")
	buffer.WriteString("</div>")

	return buffer.String()
}
