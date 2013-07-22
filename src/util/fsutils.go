package util

import "os"

/**
 * User: Medcl
 * Date: 13-7-22
 * Time: 下午12:23
 */
func CheckFileExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return false
}
