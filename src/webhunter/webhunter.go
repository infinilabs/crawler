/**
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午5:42 
 */
package webhunter

import (
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	. "net/url"
	"os"
	"regexp"
	"strings"
//	"sync"
	"time"
	. "github.com/zeebo/sbloom"
	util "util"
	. "github.com/PuerkitoBio/purell"
)

type SiteConfig struct {

	//follow page link,and walk around
	FollowLink bool

	//walking around pattern
	LinkUrlExtractRegex   *regexp.Regexp
	LinkUrlMustContain    string
	LinkUrlMustNotContain string

	//parsing url pattern,when url match this pattern,gopa will not parse urls from response of this url
	SkipPageParsePattern *regexp.Regexp

	//downloading pattern
	DownloadUrlPattern        *regexp.Regexp
	DownloadUrlMustContain    string
	DownloadUrlMustNotContain string

	//Crawling within domain
	FollowSameDomain bool
	FollowSubDomain  bool
}

type Task struct {
	Url, Request, Response []byte
}

func fetchUrl(url []byte, success chan Task, failure chan string, timeout time.Duration,config *SiteConfig) {
	t := time.NewTimer(timeout)
	defer t.Stop()

	resource := string(url)
	flg := make(chan bool, 1)

	go func() {

		defer func() {
			failure <- resource
		}()

		resp, err := http.Get(resource)
		if err != nil {
			log.Error("http.Get error!: ", resource, " , ", err)
			return
		}
		defer resp.Body.Close()
		log.Debug("getting,", resource)
		body, _ := ioutil.ReadAll(resp.Body)
		task := Task{url, nil, body}

		log.Debug("started check downloadUrlPattern,",config.DownloadUrlPattern,",",string(url))
		if config.DownloadUrlPattern.Match(url){

			if len(config.DownloadUrlMustNotContain) > 0 {
				if util.ContainStr(resource, config.DownloadUrlMustNotContain) {
					log.Debug("hit DownloadUrlMustNotContain,ignore,", resource, " , ", config.DownloadUrlMustNotContain)
					goto exitPage
				}
			}

			if len(config.DownloadUrlMustContain) > 0 {
				if !util.ContainStr(resource, config.DownloadUrlMustContain) {
					log.Debug("not hit DownloadUrlMustContain,ignore,", resource, " , ", config.DownloadUrlMustContain)
					goto exitPage
				}
			}

			savePage(url, body)
		exitPage:
		}else{
			log.Debug("does not hit DownloadUrlPattern ignoring,",resource)
		}


		success <- task
		flg <- true
	}()

	//监听通道，由于设有超时，不可能泄露
	select {
	case <-t.C:
		log.Error("fetching url time out,", resource)
	case <-flg:
		log.Debug("fetching url normal exit,", resource)
		return
	}

}

func getRootUrl(source *URL) string {
	if strings.HasSuffix(source.Path, "/") {
		return source.Host + source.Path
	} else {
		index := strings.LastIndex(source.Path, "/")
		if index > 0 {
			path := source.Path[0:index]
			return source.Host + path
		} else {
			return source.Host + "/"
		}
	}
	return ""
}

func savePage(myurl []byte, body []byte) {
	myurl1, _ := ParseRequestURI(string(myurl))
	log.Debug("url->path:", myurl1.Host, " ", myurl1.Path)

	baseDir := "data/" + myurl1.Host + "/"
	path := ""

	//making folders
	if strings.HasSuffix(myurl1.Path, "/") {
		path = baseDir + myurl1.Path
		os.MkdirAll(path, 0777)
		log.Debug("making dir:", path)
		path = (path + "default.html")
		log.Debug("no page name,use default.html:", path)

	} else {
		index := strings.LastIndex(myurl1.Path, "/")
		log.Debug("index of last /:", index)
		if index > 0 {
			path = myurl1.Path[0:index]
			path = baseDir + path
			log.Debug("new path:", path)
			os.MkdirAll(path, 0777)
			log.Debug("making dir:", path)
			path = (baseDir + myurl1.Path)
		} else {
			path = baseDir + path + "/"
			os.MkdirAll(path, 0777)
			log.Debug("making dir:", path)
			path = path + "/default.html"
		}
	}

	log.Debug("touch file,", path)
	fout, error := os.Create(path)
	if error != nil {
		log.Error(path, error)
		return
	}

	defer fout.Close()
	log.Info("saved:", path)
	fout.Write(body)

}

func ThrottledCrawl(bloomFilter *Filter,config *SiteConfig,curl chan []byte, maxGoR int, success chan Task, failure chan string) {
	maxGos := maxGoR
	numGos := 0
	for {
		if numGos > maxGos {
			log.Error("exceed maxGos,failure,")
			<-failure
			numGos -= 1
		}
		url := <-curl
		timeout := 20 * time.Second
		if !bloomFilter.Lookup(url){
			go fetchUrl(url, success, failure, timeout,config)
			numGos += 1
		}else{
			log.Debug("hit bloom filter,skipping,",string(url))
		}
	}
}

func Seed(curl chan []byte, seed string) {
	curl <- []byte(seed)
}

//var l sync.Mutex

//func init() {
//	log.Debug("[webhunter] initializing")
//}

func formatUrlForFilter(url []byte) []byte {
	src := string(url)
	log.Debug("start to normalize url:",src)
	if strings.HasSuffix(src, "/") {
		src = strings.TrimRight(src, "/")
	}
	src = strings.TrimSpace(src)
	src = strings.ToLower(src)
	return []byte(src)
}

func ExtractLinksFromTaskResponse(bloomFilter *Filter,curl chan []byte, task Task, siteConfig *SiteConfig) {
	siteUrlStr := string(task.Url)
	if siteConfig.SkipPageParsePattern.Match(task.Url) {
		log.Debug("hit SkipPageParsePattern pattern,", siteUrlStr)
		return
	}

	log.Debug("parsing external links:", siteUrlStr, ",using:", siteConfig.LinkUrlExtractRegex)
	if siteConfig.LinkUrlExtractRegex == nil {
		siteConfig.LinkUrlExtractRegex = regexp.MustCompile("src=\"(?<url1>.*?)\"|href=\"(?<url2>.*?)\"")
		log.Debug("use default linkUrlExtractRegex,", siteConfig.LinkUrlExtractRegex)
	}

	matches := siteConfig.LinkUrlExtractRegex.FindAllSubmatch(task.Response, -1)
	log.Debug("extract links with pattern:", len(matches), " match result")
	xIndex := 0
	for _, match := range matches {
		log.Debug("dealing with match result,", xIndex)
		xIndex = xIndex + 1
		url := match[2]
		filterUrl := formatUrlForFilter(url)
		log.Debug("url clean result:", string(filterUrl), ",original url:", string(url))
		filteredUrl := string(filterUrl)

		//filter error link
		if filteredUrl == "" {
			log.Debug("filteredUrl is empty,continue")
			continue
		}

		result1 := strings.HasPrefix(filteredUrl, "#")
		if result1 {
			log.Debug("filteredUrl started with: # ,continue")
			continue
		}

		result2 := strings.HasPrefix(filteredUrl, "javascript:")
		if result2 {
			log.Debug("filteredUrl started with: javascript: ,continue")
			continue
		}

		hit := false


		//		l.Lock();
		//		defer l.Unlock();

		if bloomFilter.Lookup(filterUrl) {
			log.Debug("hit bloomFilter,continue")
			hit = true
			continue
		}

		if !hit {
			currentUrlStr := string(url)

			seedUrlStr := siteUrlStr
			seedURI, err := ParseRequestURI(seedUrlStr)

			if err != nil {
				log.Error("ParseSeedURI failed!: ", seedUrlStr, " , ", err)
				continue
			}

			currentURI1, err := ParseRequestURI(currentUrlStr)
			currentURI := currentURI1
			if err != nil {
				if strings.Contains(err.Error(), "invalid URI for request") {
					log.Warn("invalid URI for request,fix relative url,", currentUrlStr)
					log.Debug("old relatived url,", currentUrlStr)
					//page based relative urls

					currentUrlStr = "http://" + seedURI.Host + "/" + currentUrlStr
					currentURI1, err = ParseRequestURI(currentUrlStr)
					currentURI = currentURI1
					if err != nil {
						log.Error("ParseCurrentURI internal failed!: ", currentUrlStr, " , ", err)
						continue
					}

					log.Debug("new relatived url,", currentUrlStr)

				} else {
					log.Error("ParseCurrentURI failed!: ", currentUrlStr, " , ", err)
					continue
				}
			}

			//relative links
			if currentURI == nil || currentURI.Host == "" {
				if strings.HasPrefix(currentURI.Path, "/") {
					//root based relative urls
					log.Debug("old relatived url,", currentUrlStr)
					currentUrlStr = "http://" + seedURI.Host + currentUrlStr
					log.Debug("new relatived url,", currentUrlStr)
				} else {
					log.Debug("old relatived url,", currentUrlStr)
					//page based relative urls
					urlPath := getRootUrl(currentURI)
					currentUrlStr = "http://" + urlPath + currentUrlStr
					log.Debug("new relatived url,", currentUrlStr)
				}
			} else {
				log.Debug("host:", currentURI.Host, " ", currentURI.Host == "")

				//resolve domain specific filter
				if siteConfig.FollowSameDomain {

					if siteConfig.FollowSubDomain {

						//TODO handler com.cn and .com,using a TLC-domain list

					} else if seedURI.Host != currentURI.Host {
						log.Debug("domain mismatch,", seedURI.Host, " vs ", currentURI.Host)
						continue
					}
				}
			}

			if len(siteConfig.LinkUrlMustContain) > 0 {
				if !util.ContainStr(currentUrlStr, siteConfig.LinkUrlMustContain) {
					log.Debug("link does not hit must-contain,ignore,", currentUrlStr, " , ", siteConfig.LinkUrlMustNotContain)
					continue
				}
			}

			if len(siteConfig.LinkUrlMustNotContain) > 0 {
				if util.ContainStr(currentUrlStr, siteConfig.LinkUrlMustNotContain) {
					log.Debug("link hit must-not-contain,ignore,", currentUrlStr, " , ", siteConfig.LinkUrlMustNotContain)
					continue
				}
			}

			//normalize url
			currentUrlStr = MustNormalizeURLString(currentUrlStr, FlagLowercaseScheme | FlagLowercaseHost | FlagUppercaseEscapes |
						FlagsUsuallySafeGreedy | FlagRemoveDuplicateSlashes | FlagRemoveFragment)
			log.Debug("normalized url:", currentUrlStr)
			currentUrlByte := []byte(currentUrlStr)
			if (!bloomFilter.Lookup(currentUrlByte)) {

//				if(CheckIgnore(currentUrlStr)){}

				log.Info("enqueue:", currentUrlStr)
				curl <- currentUrlByte
				bloomFilter.Add(currentUrlByte)
			}
			bloomFilter.Add([]byte(filterUrl))
		} else {
			log.Debug("hit bloom filter,ignore,", string(url))
		}
	}

	log.Info("all links within ", siteUrlStr," is done")
}

