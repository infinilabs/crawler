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

package util

import (
	"bytes"
	"encoding/json"
	. "strings"
	"unicode"
	"unicode/utf16"
)

func ContainStr(s, substr string) bool {
	return Index(s, substr) != -1
}

func StringToUTF16(s string) []uint16 {
	return utf16.Encode([]rune(s + "\x00"))
}

func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

// Removes all unnecessary whitespaces
func MergeSpace(in string) (out string) {
	var buffer bytes.Buffer
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				buffer.WriteString(" ")
			}
			white = true
		} else {
			buffer.WriteRune(c)
			white = false
		}
	}
	return TrimSpace(buffer.String())
}

func ToJson(in interface{}, indent bool) string {
	var b []byte
	if indent {
		b, _ = json.MarshalIndent(in, " ", " ")
	} else {
		b, _ = json.Marshal(in)
	}
	return string(b)
}
