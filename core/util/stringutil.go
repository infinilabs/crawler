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
	"errors"
	"runtime"
	"strconv"
	. "strings"
	"time"
	"unicode"
	"unicode/utf16"
)

func ContainStr(s, substr string) bool {
	return Index(s, substr) != -1
}

func StringToUTF16(s string) []uint16 {
	return utf16.Encode([]rune(s + "\x00"))
}

func SubStringWithSuffix(str string, length int, suffix string) string {
	if len(str) > length {
		str = SubString(str, 0, length) + suffix
	}
	return str
}

func UnicodeIndex(str, substr string) int {
	// 子串在字符串的字节位置
	result := Index(str, substr)
	if result >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:result]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		result = len(rs)
	}

	return result
}

func SubString(str string, begin, length int) (substr string) {
	rs := []rune(str)
	lth := len(rs)

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

	return string(rs[begin:end])
}

func NoWordBreak(in string) string {
	return Replace(in, "\n", " ", -1)
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

func IntToString(num int) string {
	return strconv.Itoa(num)
}

func ToInt(str string) (int, error) {
	if IndexAny(str, ".") > 0 {
		nonFractionalPart := Split(str, ".")
		return strconv.Atoi(nonFractionalPart[0])
	} else {
		return strconv.Atoi(str)
	}

}

func FormatTime(date *time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

func GetRuntimeErrorMessage(r runtime.Error) string {
	if r != nil {
		return r.Error()
	}
	panic(errors.New("nil runtime error"))
	return "nil runtime error"
}

func XSSHandle(src string) string {
	src = Replace(src, ">", "&lt; ", -1)
	src = Replace(src, ">", "&gt; ", -1)
	return src
}
