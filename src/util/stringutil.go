/** 
 * User: Medcl
 * Date: 13-7-11
 * Time: 下午9:51 
 */
package stringutil

import . "strings"

func ContainStr(s, substr string) bool {
	return Index(s, substr) != -1
}
