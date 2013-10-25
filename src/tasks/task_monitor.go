/**
 * Created with IntelliJ IDEA.
 * User: medcl
 * Date: 13-10-25
 * Time: 下午1:33
 */
package tasks
import (
	log "logging"
//	"io/ioutil"
	//	"net/http"
//	. "net/url"
	"os"
	//	"regexp"
//	"strings"
	time    "time"
//	. "github.com/PuerkitoBio/purell"
	//	. "github.com/zeebo/sbloom"
	//	"kafka"
	//	"math/rand"
	//	"strconv"
	. "types"
	util "util"
	//	utils "util"
	//	bloom "github.com/zeebo/sbloom"
	//	"hash/fnv"
	"bufio"
)

func LoadTaskFromLocalFile(pendingFetchUrls chan []byte, runtimeConfig RuntimeConfig, quit *chan bool, offsets *RoutingOffset){

	log.Error("LoadTaskFromLocalFile task started.")
	path := runtimeConfig.PathConfig.PendingFetchLog
	//touch local's file
	//read all of line
	//if hit the EOF,will wait 2s,and then reopen the file,and try again,may be check the time of last modified

waitFile:
	if (!util.CheckFileExists(path)) {
		log.Error("waiting file crate",path)
		time.Sleep(2*time.Millisecond)
		goto waitFile
	}

	offset := 0
	FetchFileWithOffset2(runtimeConfig,pendingFetchUrls,path, offset)


}


func FetchFileWithOffset2(runtimeConfig RuntimeConfig,pendingFetchUrls chan []byte,path string, skipOffset int) {

	offset := 0
	time1, _ := util.FileMTime(path)
	log.Error("start touch time:", time1)

	f, err := os.Open(path)
	if err != nil {
		log.Error("error opening file,", path, " ", err)
		return
	}

	r := bufio.NewReader(f)
	s, e := util.Readln(r)
	offset = 0
	log.Trace("new offset:", offset)

	for e == nil {
		offset = offset + 1
		//TODO use byte offset instead of lines
		if (offset > skipOffset) {
			ParsedSavedFileLog2(runtimeConfig,pendingFetchUrls,s)
		}
		s, e = util.Readln(r)
		//todo store offset
	}
	log.Trace("end offset:", offset, "vs ", skipOffset)

waitUpdate:
	time2, _ := util.FileMTime(path)

	log.Trace("2nd touch time:", time2)

	if (time2 > time1) {
		log.Debug("file has been changed,restart parse")
		FetchFileWithOffset2(runtimeConfig,pendingFetchUrls,path, offset)
	}else {
		log.Trace("waiting file update",path)
		time.Sleep(3*time.Millisecond)
		goto waitUpdate
	}
}

func ParsedSavedFileLog2(runtimeConfig RuntimeConfig,pendingFetchUrls chan []byte,url string) {
	if (url != "") {
		log.Trace("start parse filelog:", url)

		if(runtimeConfig.Storage.CheckFetchedUrl([]byte(url))){
			log.Debug("hit fetch filter ignore,",url)
			return
		}
		log.Debug("new task extracted from saved page:", url)
		pendingFetchUrls <- []byte(url)
	}
}
