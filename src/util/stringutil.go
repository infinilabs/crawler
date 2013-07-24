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
