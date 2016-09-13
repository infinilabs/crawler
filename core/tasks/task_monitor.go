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

package tasks

import (
	"bufio"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/util"
	"os"
	"time"
)

func LoadTaskFromLocalFile(pendingFetchUrls chan []byte, runtimeConfig *RuntimeConfig) {

	log.Trace("LoadTaskFromLocalFile task started.")
	path := runtimeConfig.PathConfig.PendingFetchLog
	//touch local's file
	//read all of line
	//if hit the EOF,will wait 2s,and then reopen the file,and try again,may be check the time of last modified

waitFile:
	if !util.FileExists(path) {
		//log.Trace("waiting file create:", path)
		time.Sleep(100 * time.Millisecond)
		goto waitFile
	}
	//var storage = runtimeConfig.Storage

	//var offset int64 = storage.LoadOffset(runtimeConfig.PathConfig.PendingFetchLog + ".offset")
	//FetchFileWithOffset2(*runtimeConfig, pendingFetchUrls, path, offset)

}

func FetchFileWithOffset2(runtimeConfig RuntimeConfig, pendingFetchUrls chan []byte, path string, skipOffset int64) {

	var offset int64
	offset = 0
	time1, _ := util.FileMTime(path)
	log.Trace("start touch time:", time1)

	f, err := os.Open(path)
	if err != nil {
		log.Trace("error opening file,", path, " ", err)
		return
	}
	//var storage = runtimeConfig.Storage

	r := bufio.NewReader(f)
	s, e := util.Readln(r)
	offset = 0
	log.Trace("new offset:", offset)

	for e == nil {
		offset = offset + 1
		//TODO use byte offset instead of lines
		if offset > skipOffset {
			ParsedSavedFileLog2(runtimeConfig, pendingFetchUrls, s)
		}

		//storage.PersistOffset(runtimeConfig.PathConfig.PendingFetchLog+".offset", offset)

		s, e = util.Readln(r)
		//todo store offset
	}
	log.Trace("end offset:", offset, "vs ", skipOffset)

waitUpdate:
	time2, _ := util.FileMTime(path)

	//log.Trace("2nd touch time:", time2)

	if time2 > time1 {
		log.Debug("file has been changed,restart parse")
		FetchFileWithOffset2(runtimeConfig, pendingFetchUrls, path, offset)
	} else {
		//log.Trace("waiting file update,", path)
		time.Sleep(10 * time.Millisecond)
		goto waitUpdate
	}
}

func ParsedSavedFileLog2(runtimeConfig RuntimeConfig, pendingFetchUrls chan []byte, url string) {
	if url != "" {
		log.Trace("start parse filelog:", url)

		//var storage = runtimeConfig.Storage

		//if storage.UrlHasFetched([]byte(url)) {
		//	log.Debug("hit fetch filter ignore,", url)
		//	return
		//}
		log.Debug("new task extracted from saved page:", url)
		pendingFetchUrls <- []byte(url)
	}
}
