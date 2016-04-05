/**
 * User: Medcl
 * Date: 13-7-11
 * Time: 下午9:51
 */
package util

import (
	. "strings"
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
