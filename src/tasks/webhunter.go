/**
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午5:42 
 */
package tasks

import (
	log "github.com/cihub/seelog"
//	"io/ioutil"
//	"net/http"
	. "net/url"
//	"os"
	"regexp"
	"strings"
//	"time"
	. "github.com/zeebo/sbloom"
	util "util"
	. "github.com/PuerkitoBio/purell"
	"kafka"
//config	"config"
//	"strconv"
//	utils "util"
	. "types"
)


//
////fetch url's content
//func fetchUrl(url []byte, failure chan string, timeout time.Duration,config *TaskConfig,bloomFilter *Filter) {
//	t := time.NewTimer(timeout)
//	defer t.Stop()
//	log.Debug("enter fetchUrl method")
//
//	resource := string(url)
//	bloomFilter.Add(url)
//
//	//checking fetchUrlPattern
//	log.Debug("started check fetchUrlPattern,",config.FetchUrlPattern,",",string(url))
//	if config.FetchUrlPattern.Match(url){
//		log.Debug("match fetch url pattern,",resource)
//		if len(config.FetchUrlMustNotContain) > 0 {
//			if util.ContainStr(resource, config.FetchUrlMustNotContain) {
//				log.Debug("hit FetchUrlMustNotContain,ignore,", resource, " , ", config.FetchUrlMustNotContain)
//				return
//			}
//		}
//
//		if len(config.FetchUrlMustContain) > 0 {
//			if !util.ContainStr(resource, config.FetchUrlMustContain) {
//				log.Debug("not hit FetchUrlMustContain,ignore,", resource, " , ", config.FetchUrlMustContain)
//				return
//			}
//		}
//	}else{
//		log.Debug("does not hit FetchUrlPattern ignoring,",resource)
//	}
//
//
//
//	log.Debug("start fetch url,",resource)
//	flg := make(chan bool, 1)
//
//	go func() {
//
//		defer func() {
//			failure <- resource
//		}()
//
//		resp, err := http.Get(resource)
//		if err != nil {
//			log.Error("http.Get error!: ", resource, " , ", err)
//			return
//		}
//		defer resp.Body.Close()
//		log.Debug("getting,", resource)
//		body, _ := ioutil.ReadAll(resp.Body)
//		task := Task{url, nil, body}
//
//		log.Debug("started check savingUrlPattern,",config.SavingUrlPattern,",",string(url))
//		if config.SavingUrlPattern.Match(url){
//			log.Debug("match saving url pattern,",resource)
//			if len(config.SavingUrlMustNotContain) > 0 {
//				if util.ContainStr(resource, config.SavingUrlMustNotContain) {
//					log.Debug("hit SavingUrlMustNotContain,ignore,", resource, " , ", config.SavingUrlMustNotContain)
//					goto exitPage
//				}
//			}
//
//			if len(config.SavingUrlMustContain) > 0 {
//				if !util.ContainStr(resource, config.SavingUrlMustContain) {
//					log.Debug("not hit SavingUrlMustContain,ignore,", resource, " , ", config.SavingUrlMustContain)
//					goto exitPage
//				}
//			}
//
//			task.Save(url, body)
//
//		exitPage:
//		}else{
//			log.Debug("does not hit SavingUrlPattern ignoring,",resource)
//		}
//
//		log.Debug("task enqueue,",resource)
////		success <- task
//		flg <- true
//	}()
//
//	//监听通道，由于设有超时，不可能泄露
//	select {
//	case <-t.C:
//		log.Error("fetching url time out,", resource)
//	case <-flg:
//		log.Debug("fetching url normal exit,", resource)
//		return
//	}
//
//}

//parse to get url root
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
//
////saving page
//func savePage(myurl []byte, body []byte) {
//	urlStr := string(myurl)
//	log.Debug("start saving url,", urlStr)
//	myurl1, _ := ParseRequestURI(urlStr)
//	log.Debug("url->path:", myurl1.Host, " ", myurl1.Path)
//
//	baseDir := "data/" + myurl1.Host + "/"
//	baseDir = strings.Replace(baseDir, `:`, `_`, -1)
//
//	log.Debug("replaced:", baseDir)
//	path := ""
//
//	//making folders
//	if strings.HasSuffix(urlStr, "/") {
//		path = baseDir + myurl1.Path
//		os.MkdirAll(path, 0777)
//		log.Debug("making dir:", path)
//		path = (path + "default.html")
//		log.Debug("no page name,use default.html:", path)
//
//	} else {
//		index := strings.LastIndex(myurl1.Path, "/")
//		log.Debug("index of last /:", index,",",myurl1.Path)
//		if index >= 0 {
//			path = myurl1.Path[0:index]
//			path = baseDir + path
//			log.Debug("new path:", path)
//			os.MkdirAll(path, 0777)
//			log.Debug("making dir:", path)
//			path = (baseDir + myurl1.Path)
//		} else {
//			path = baseDir + path + "/"
//			os.MkdirAll(path, 0777)
//			log.Debug("making dir:", path)
//			path = path + "/default.html"
//		}
//	}
//
//	log.Debug("touch file,", path)
//	fout, error := os.Create(path)
//	if error != nil {
//		log.Error(path, error)
//		return
//	}
//
//	defer fout.Close()
//	log.Info("saved:", urlStr,",", path)
//	fout.Write(body)
//	log.Debug("end saving url,", urlStr)
//
//}



////crawl with limited gorouting
//func ThrottledCrawl(bloomFilter *Filter,taskConfig *TaskConfig,kafkaConfig *config.KafkaConfig,curl chan []byte,
//maxGoR int, success chan Task, failure chan string,quit []*chan bool,offsets []*RoutingOffset) {
//	log.Debug("enter kafka consume")
////	for i := 0; i < maxGoR-1; i++ {
//	   i:=0
//		//TODO
//		go func(partition int) {
//			log.Debug("init go routing,",i," of ",maxGoR)
//
//			offset:= *offsets[partition]
//
//			broker := kafka.NewBrokerConsumer(kafkaConfig.Hostname,taskConfig.Name , partition, offset.Offset, kafkaConfig.MaxSize)
//
////			printmessage:=true
//
//			consumerCallback := func(msg *kafka.Message) {
////				if printmessage {
//
//				url :=  msg.Payload()
//					log.Debug("kafka message offset: " + strconv.FormatUint(msg.Offset(), 10) )
//					timeout := 10 * time.Second
//					if !bloomFilter.Lookup(url){
//						fetchUrl(url, success, failure, timeout,taskConfig,bloomFilter)
////						bloomFilter.Add(url)
//					}else{
//						log.Debug("hit bloom filter,skipping,",string(url))
//					}
//					offsetV:=msg.Offset()
//					offset.Offset=offsetV
//
//					path:="offset_"+strconv.FormatInt(int64(i),10)+".tmp"
//					path_new:="offset_"+strconv.FormatInt(int64(i),10)
//					fout, error := os.Create(path)
//					if error != nil {
//						log.Error(path, error)
//						return
//					}
//
//					defer fout.Close()
//					log.Debug("saved offset: ", offsetV)
//					fout.Write([]byte(strconv.FormatUint(msg.Offset(), 10)))
//					utils.CopyFile(path,path_new)
//				}
////			}
//
////			consumerForever:=true
////			if consumerForever {
//				msgChan := make(chan *kafka.Message)
//				go broker.ConsumeOnChannel(msgChan, 10, *quit[i])
//				for msg := range msgChan {
//					if msg != nil {
//						log.Debug("consume messaging.",string(msg.Payload()))
//						consumerCallback(msg)
//					} else {
//						break
//					}
//				}
////			} else {
////				broker.Consume(consumerCallback)
////			}
//		}(i)
////	}
//	log.Debug("exit kafka consume")
//}

//var l sync.Mutex

//func init() {
//	log.Debug("[webhunter] initializing")
//}

//format url,prepare for bloom filter
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

//func ExtractLinksFromTaskResponse(bloomFilter *Filter,curl chan []byte, task Task, siteConfig *SiteConfig) {
func ExtractLinksFromTaskResponse(bloomFilter *Filter,broker *kafka.BrokerPublisher, task Task, siteConfig *TaskConfig) {
	siteUrlStr := string(task.Url)
	log.Debug("enter links extract,",siteUrlStr)
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
			currentUrlStr=strings.Trim(currentUrlStr," ")

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
					log.Warn("invalid URI for request,fix relative url,original:", currentUrlStr)
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
						FlagRemoveUnnecessaryHostDots |FlagRemoveDuplicateSlashes | FlagRemoveFragment)
			log.Debug("normalized url:", currentUrlStr)
			currentUrlByte := []byte(currentUrlStr)
			if (!bloomFilter.Lookup(currentUrlByte)) {

//				if(CheckIgnore(currentUrlStr)){}

				log.Debug("enqueue:", currentUrlStr)

				//TODO 如果使用分布式队列，则不使用go的channel，抽象出接口
//				curl <- currentUrlByte

				broker.Publish(kafka.NewMessage(currentUrlByte))

//				bloomFilter.Add(currentUrlByte)
			}
//			bloomFilter.Add([]byte(filterUrl))
		} else {
			log.Debug("hit bloom filter,ignore,", string(url))
		}
		log.Debug("exit links extract,",siteUrlStr)

	}

	log.Info("all links within ", siteUrlStr," is done")
}

