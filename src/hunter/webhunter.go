/** 
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午5:42 
 */
package hunter

import (
	"net/http"
	 log "github.com/cihub/seelog"
	"io/ioutil"
	"regexp"
	"util/bloom"
	"sync"
	"util/stringutil"
	"time"

)

type SiteConfig struct{

	//walking around pattern
	LinkUrlExtractRegex *regexp.Regexp
	LinkUrlMustContain string
	LinkUrlMustNotContain string

	//downloading pattern
	DownloadUrlPattern *regexp.Regexp
	DownloadUrlMustContain string
	DownloadUrlMustNotContain string
}

type  Task struct{
  Url,Request,Response []byte
}

func fetchUrl(url []byte,success chan Task,failure chan string,timeout time.Duration){
	t := time.NewTimer(timeout)
	defer t.Stop()

	resource := string(url)
	flg := make(chan bool, 1)

	go func() {

		defer func () {
			failure <- resource
		}()

		resp, err := http.Get(resource)
		if err != nil {
			log.Error("we have an error!: ", err)
			return
		}
		defer resp.Body.Close()
		log.Debug("getting,", resource)
		body, _ := ioutil.ReadAll(resp.Body)
		task := Task{url,nil, body}

		savePage(url,body)

		success <- task
		flg <- true
	}()

	//监听通道，由于设有超时，不可能泄露
	select {
	case <-t.C:
		log.Error("fetching url time out,",resource)
	case <-flg:
		log.Debug("fetching url normal exit,",resource)
		return
	}

}


func savePage(url []byte,body []byte){
	log.Info("saving page,",string(url),string(body))
}


func ThrottledCrawl(curl chan []byte, success chan Task, failure chan string) {
	maxGos := 10
	numGos := 0
	for {
		if numGos > maxGos {
			<-failure
			numGos -= 1
		}
		url := string(<-curl)
//		if _, ok := visited[url]; !ok {
		timeout := 20 * time.Second
		go fetchUrl([]byte(url), success, failure,timeout)
			numGos += 1
//		}
//		visited[url] += 1
	}
}

func Seed(curl chan []byte,seed string) {
	curl <- []byte(seed)
}

var f *bloom.Filter64
var l sync.Mutex

func init(){

//	log.Debug("[webhunter] initializing")

	// Create a bloom filter which will contain an expected 100,000 items, and which
	// allows a false positive rate of 1%.
	f = bloom.New64(1000000, 0.01)

}

func GetUrls(curl chan []byte, task Task, siteConfig SiteConfig) {
	log.Debug("parsing external links:",string(task.Url))
	if(siteConfig.LinkUrlExtractRegex==nil){
		log.Debug("use default linkUrlExtractRegex,",siteConfig.LinkUrlExtractRegex)
		siteConfig.LinkUrlExtractRegex=regexp.MustCompile("<a.*?href=[\"'](http.*?)[\"']")
	}
	matches := siteConfig.LinkUrlExtractRegex.FindAllSubmatch(task.Response, -1)
	for _, match := range matches {
		url := match[1]

		hit := false
		l.Lock();
		if(f.Test(url)){
			hit=true
		}
		l.Unlock();

		if(!hit){
			myurl:=string(url)
			if(len(siteConfig.LinkUrlMustContain)>0){
				if(!stringutil.ContainStr(myurl,siteConfig.LinkUrlMustContain)){
					log.Debug("link does not hit must-contain,ignore,",myurl,",",siteConfig.LinkUrlMustNotContain)
					break;
				}
			}

			if(len(siteConfig.LinkUrlMustNotContain)>0){
				if(stringutil.ContainStr(myurl,siteConfig.LinkUrlMustNotContain)){
					log.Debug("link hit must-not-contain,ignore,",myurl,",",siteConfig.LinkUrlMustNotContain)
					break;
				}
			}

			log.Info("enqueue:",string(url))
			curl <- match[1]
			f.Add([]byte(url))
		}else{
			log.Debug("hit bloom filter,ignore,",string(url))
		}

		//TODO 判断url是否已经请求过，并且判断url pattern，如果满足处理条件，则继续进行处理，否则放弃

	}
}
