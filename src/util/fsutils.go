package util

import (
	"os"
	"io"
log	"github.com/cihub/seelog"
)

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

func CopyFile(src,dst string)(w int64,err error){
	srcFile,err := os.Open(src)
	if err!=nil{
		log.Error(err.Error())
		return
	}
	defer srcFile.Close()

	dstFile,err := os.Create(dst)

	if err!=nil{
		log.Error(err.Error())
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile,srcFile)
}
