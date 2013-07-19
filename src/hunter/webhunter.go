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
	"strings"
	"os"
.	"net/url"
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

	//Crawling within domain
	FollowSameDomain  bool
	FollowSubDomain  bool
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


func savePage(myurl []byte,body []byte){
	myurl1,_:=ParseRequestURI(string(myurl))
	log.Debug("url->path:",myurl1.Host," ",myurl1.Path)

	baseDir:="data/"+myurl1.Host+"/"
	path:=baseDir

	//making folders
	if(strings.HasSuffix(myurl1.Path,"/")){
		path=baseDir+myurl1.Path
		os.MkdirAll(path,0777)
		log.Debug("making dir:",path)
		path=(path+"default.html")
		log.Debug("no page name,use default.html:",path)

	}else{
	    index:= strings.LastIndex(myurl1.Path,"/")
		log.Info("index of last /:",index)
		if index>0{
			path= myurl1.Path[0:index]
			path=baseDir+path
			log.Debug("new path:",path)
			os.MkdirAll(path,0777)
			log.Debug("making dir:",path)
			path=(baseDir+myurl1.Path)
		}else{
			path= baseDir+path+"/"
			os.MkdirAll(path,0777)
			log.Debug("making dir:",path)
			path=path+"/default.html"
		}


	}




	log.Debug("touch file,",path)
	fout,error:=os.Create(path)
	if error!=nil{
		log.Error(path,error)
		return
	}

	defer  fout.Close()
	log.Info("file saved:",path)
	fout.Write(body)

}


func ThrottledCrawl(curl chan []byte,maxGoR int, success chan Task, failure chan string) {
	maxGos := maxGoR
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

func formatUrlForFilter(url []byte) []byte{
	src:=string(url)
	if(strings.HasSuffix(src,"/")){
		src= strings.TrimRight(src,"/");
	}
	src=strings.TrimSpace(src)
	src=strings.ToLower(src)
	return []byte(src)
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

		log.Debug("original filter url,",string(url))
		filterUrl:=formatUrlForFilter(url)
		log.Debug("format filter url,",string(filterUrl))

		hit := false
		l.Lock();
		if(f.Test(filterUrl)){
			hit=true
		}

		if(!hit){
			myurl:=string(url)
			seedUrl:=string(task.Url)

			myurl1,_:=ParseRequestURI(seedUrl)
			myurl2,_:=ParseRequestURI(myurl)
			if(siteConfig.FollowSameDomain){

				if(siteConfig.FollowSubDomain){

				   //TODO handler com.cn and .com,using a TLC-domain list

				}else if(myurl1.Host !=myurl2.Host){
					log.Debug("domain mismatch,",myurl1.Host," vs ",myurl2.Host)
					break;
				}
			}


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
			f.Add([]byte(filterUrl))
		}else{
			log.Debug("hit bloom filter,ignore,",string(url))
		}

		l.Unlock();

		//TODO 判断url是否已经请求过，并且判断url pattern，如果满足处理条件，则继续进行处理，否则放弃

	}
}
