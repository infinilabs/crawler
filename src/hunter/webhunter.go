/** 
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午5:42 
 */
package hunter

import (
	"net/http"
	"log"
	"io/ioutil"
	"regexp"
	"util/bloom"
	"sync"
	"util/stringutil"
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

func fetchUrl(url []byte,success chan Task,failure chan string){
	resource := string(url)
	defer func () {
		failure <- resource
	}()


	resp, err := http.Get(resource)
	if err != nil {
		log.Println("we have an error!: ", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("getting %v\n", resource)
	body, _ := ioutil.ReadAll(resp.Body)
	task := Task{url,nil, body}
//	log.Printf("Response %v\n",string(body))
	success <- task

}


func ThrottledCrawl(curl chan []byte, success chan Task, failure chan string, visited map[string]int) {
	maxGos := 10
	numGos := 0
	for {
		if numGos > maxGos {
			<-failure
			numGos -= 1
		}
		url := string(<-curl)
		if _, ok := visited[url]; !ok {
			go fetchUrl([]byte(url), success, failure)
			numGos += 1
		}
		visited[url] += 1
	}
}

func Seed(curl chan []byte,seed string) {
	curl <- []byte(seed)
}

var f *bloom.Filter64
var l sync.Mutex

func init(){
	log.Print("[webhunter] initializing")

	// Create a bloom filter which will contain an expected 100,000 items, and which
	// allows a false positive rate of 1%.
	f = bloom.New64(1000000, 0.01)

}

//func ContainStr(s, substr string) bool {
//	return Index(s, substr) != -1
//}

func GetUrls(curl chan []byte, task Task, siteConfig SiteConfig) {
	log.Print("parsing external links:",string(task.Url))
	if(siteConfig.LinkUrlExtractRegex==nil){
		log.Print("use default linkUrlExtractRegex,",siteConfig.LinkUrlExtractRegex)
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
					log.Print("link does not hit must-contain,ignore,",myurl,",",siteConfig.LinkUrlMustNotContain)
					break;
				}
			}

			if(len(siteConfig.LinkUrlMustNotContain)>0){
				if(stringutil.ContainStr(myurl,siteConfig.LinkUrlMustNotContain)){
					log.Print("link hit must-not-contain,ignore,",myurl,",",siteConfig.LinkUrlMustNotContain)
					break;
				}
			}

			log.Print("enqueue:",string(url))
			curl <- match[1]
			f.Add([]byte(url))
		}else{
			log.Print("hit bloom filter,ignore,",string(url))
		}

		//TODO 判断url是否已经请求过，并且判断url pattern，如果满足处理条件，则继续进行处理，否则放弃

	}
}
