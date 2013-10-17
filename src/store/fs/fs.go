/** 
 * User: Medcl
 * Date: 13-7-10
 * Time: 下午10:55 
 */
package fs

import (
	util "util"
	log "github.com/cihub/seelog"
)


type FsStore struct{}

func (this *FsStore) Store(url string, data []byte){
	util.FilePutContentWithByte(url,data)
}

func (this *FsStore)  Get(key string) []byte {
	file,error:= util.FileGetContent(key)
	if(error!=nil){
		log.Error("get file:",key,error)
	}
	return file
}

func (this *FsStore)  List(from int, size int) [][]byte{
	  return nil
}

func (this *FsStore) TaskEnqueue(url []byte){
	 log.Info("task enqueue:",string(url))
}

