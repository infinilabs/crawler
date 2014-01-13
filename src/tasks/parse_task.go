/** 
 * User: Medcl
 * Date: 13-7-25
 * Time: 下午9:19 
 */
//subscribe local file channel,and parse urls
package tasks

import (
	log "logging"
	"io/ioutil"
	//	"net/http"
	. "net/url"
	"os"
	//	"regexp"
	"strings"
	time    "time"
	. "github.com/PuerkitoBio/purell"
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

//var parseFilter  *bloom.Filter
func init() {
	//	parseFilter = bloom.NewFilter(fnv.New64(), 1000000)
	//	log.Warn("init bloom filter")
}

func loadFileContent(fileName string) []byte {
	if util.CheckFileExists(fileName) {
		log.Trace("found fileName,start loading:", fileName)
		n, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Error("loadFile", err, ",", fileName)
			return nil
		}
		return n
	}
	return nil
}

func extractLinks(runtimeConfig RuntimeConfig, fileUrl string , fileName []byte, body []byte) {

	//	siteUrlStr := string(fileName)
	//	siteUrlStr = strings.TrimLeft(siteUrlStr, "data/")
	//	siteUrlStr = "http://" + siteUrlStr
	//	log.Debug("fileName to Url:", string(fileName), ",", siteUrlStr)

	siteUrlStr := fileUrl
	siteConfig:=runtimeConfig.TaskConfig

	siteUrlByte := []byte(siteUrlStr)
	log.Debug("enter links extract,", siteUrlStr)
	if siteConfig.SkipPageParsePattern.Match(siteUrlByte) {
		log.Debug("hit SkipPageParsePattern pattern,", siteUrlStr)
		return
	}

	log.Debug("parsing external links:", siteUrlStr, ",using:", siteConfig.LinkUrlExtractRegex)

	matches := siteConfig.LinkUrlExtractRegex.FindAllSubmatch(body, -1)
	log.Debug("extract links with pattern,total matchs:", len(matches), " match result,", string(fileName))
	xIndex := 0
	for _, match := range matches {
		log.Debug("dealing with match result,", xIndex)
		xIndex = xIndex + 1
		url := match[siteConfig.LinkUrlExtractRegexGroupIndex]
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

		if runtimeConfig.Storage.CheckWalkedUrl(filterUrl)||runtimeConfig.Storage.CheckFetchedUrl(filterUrl)||runtimeConfig.Storage.CheckPendingFetchUrl(filterUrl) {
			log.Debug("hit bloomFilter,continue")
			hit = true
			continue
		}

		if !hit {
			currentUrlStr := string(url)
			currentUrlStr = strings.Trim(currentUrlStr, " ")

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
					log.Debug("invalid URI for request,fix relative url,original:", currentUrlStr)
					//					log.Debug("old relatived url,", currentUrlStr)
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

			//			relative links
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
						//continue
					}
					//TODO follow all or list of domain
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
						FlagRemoveUnnecessaryHostDots | FlagRemoveDuplicateSlashes | FlagRemoveFragment)
			log.Debug("normalized url:", currentUrlStr)
			currentUrlByte := []byte(currentUrlStr)
			if !(runtimeConfig.Storage.CheckWalkedUrl(currentUrlByte)||runtimeConfig.Storage.CheckFetchedUrl(currentUrlByte)||runtimeConfig.Storage.CheckPendingFetchUrl(currentUrlByte)){
			//bloomFilter.Lookup(currentUrlByte) {

				//								if(CheckIgnore(currentUrlStr)){}

				//				log.Info("enqueue fetch: ", currentUrlStr)

				//				broker.Publish(kafka.NewMessage(currentUrlByte))


				//copied form fetchTask,TODO refactor
				//checking fetchUrlPattern
				log.Debug("started check fetchUrlPattern,", currentUrlStr)
				if siteConfig.FetchUrlPattern.Match(currentUrlByte) {
					log.Debug("match fetch url pattern,", currentUrlStr)
					if len(siteConfig.FetchUrlMustNotContain) > 0 {
						if util.ContainStr(currentUrlStr, siteConfig.FetchUrlMustNotContain) {
							log.Debug("hit FetchUrlMustNotContain,ignore,", currentUrlStr)
							continue
						}
					}

					if len(siteConfig.FetchUrlMustContain) > 0 {
						if !util.ContainStr(currentUrlStr, siteConfig.FetchUrlMustContain) {
							log.Debug("not hit FetchUrlMustContain,ignore,", currentUrlStr)
							continue
						}
					}
				} else {
					log.Debug("does not hit FetchUrlPattern ignoring,", currentUrlStr)
					continue
				}

				if(!runtimeConfig.Storage.CheckPendingFetchUrl(currentUrlByte)){
					log.Debug("log new pendingFetch url", currentUrlStr)
					runtimeConfig.Storage.LogPendingFetchUrl(runtimeConfig.PathConfig.PendingFetchLog,currentUrlStr)
					runtimeConfig.Storage.AddPendingFetchUrl(currentUrlByte)
				}else{
					log.Debug("hit new pendingFetch filter,ignore:", currentUrlStr)
				}
//				pendingUrls <- currentUrlByte

				//	TODO pendingFetchFilter			bloomFilter.Add(currentUrlByte)
			}else{
				log.Debug("hit bloom filter,ignore:", currentUrlStr)
			}
			//			bloomFilter.Add([]byte(filterUrl))
		} else {
			log.Debug("hit bloom filter,ignore,", string(url))
		}
		log.Debug("exit links extract,", siteUrlStr)

	}

	//TODO 处理ruled fetch pattern

	log.Info("all links within ", siteUrlStr, " is done")
}

func ParseGo(pendingUrls chan []byte, runtimeConfig RuntimeConfig, quit *chan bool, offsets *RoutingOffset) {
	log.Info("parsing task started.")
	path := runtimeConfig.PathConfig.SavedFileLog
	//touch local's file
	//read all of line
	//if hit the EOF,will wait 2s,and then reopen the file,and try again,may be check the time of last modified

waitFile:
	if (!util.CheckFileExists(path)) {
		log.Trace("waiting file create",path)
		time.Sleep(10*time.Millisecond)
		goto waitFile
	}

	var offset int64= runtimeConfig.Storage.LoadOffset(runtimeConfig.PathConfig.SavedFileLog + ".offset")
	FetchFileWithOffset(runtimeConfig,path, offset)
}

func FetchFileWithOffset(runtimeConfig RuntimeConfig,path string, skipOffset int64) {

	var offset int64= 0

	time1, _ := util.FileMTime(path)
	log.Debug("start touch time:", time1)

	f, err := os.Open(path)
	if err != nil {
		log.Debug("error opening file,", path, " ", err)
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
			ParsedSavedFileLog(runtimeConfig,s)
		}

		runtimeConfig.Storage.PersistOffset(runtimeConfig.PathConfig.SavedFileLog + ".offset",offset)

		s, e = util.Readln(r)
		//todo store offset
	}
	log.Trace("end offset:", offset, "vs ", skipOffset)

waitUpdate:
	time2, _ := util.FileMTime(path)

	log.Trace("2nd touch time:", time2)

	if (time2 > time1) {
		log.Trace("file has been changed,restart parse")
		FetchFileWithOffset(runtimeConfig,path, offset)
	}else {
		log.Trace("waiting file update",path)
		time.Sleep(10*time.Millisecond)
		goto waitUpdate
	}
}

func ParsedSavedFileLog(runtimeConfig RuntimeConfig,fileLog string) {
	if (fileLog != "") {
		log.Debug("start parse filelog:", fileLog)
		//load file's content,and extract links

		stringArray:=strings.Split(fileLog,"|||");
		fileUrl:=stringArray[0]
		fileName:=[]byte(stringArray[1])

		if(runtimeConfig.Storage.CheckParsedFile(fileName)){
			log.Debug("hit parse filter ignore,",string(fileName))
			return
		}

		fileContent := loadFileContent(string(fileName))
		runtimeConfig.Storage.AddParsedFile(fileName)

		if fileContent != nil {
//			log.Debug("partition:", partition, ",parse fileName:", string(fileName))

			//extract urls to fetch queue.
			extractLinks(runtimeConfig,fileUrl, fileName, fileContent)
//			offsetV := msg.Offset()
//			offset.Offset = offsetV
//
//			path := taskConfig.BaseStoragePath+     "task/parse_offset_" + strconv.FormatInt(int64(partition), 10) + ".tmp"
//			path_new := taskConfig.BaseStoragePath+     "task/parse_offset_" + strconv.FormatInt(int64(partition), 10)
//			fout, error := os.Create(path)
//			if error != nil {
//				log.Error(path, error)
//				return
//			}
//
//			defer fout.Close()
//			log.Debug("partition:", partition, ",saved offset:", offsetV)
//			fout.Write([]byte(strconv.FormatUint(msg.Offset(), 10)))
//			utils.CopyFile(path, path_new)
		}
	}
}

func ParseLinks(pendingUrls chan []byte, runtimeConfig *RuntimeConfig, quit *chan bool, offsets *RoutingOffset, MaxGoRoutine int) {
	//func ParseLinks(pendingUrls chan []byte, bloomFilter *Filter, taskConfig *TaskConfig, kafkaConfig *config.KafkaConfig, quit *chan bool, offsets *RoutingOffset, MaxGoRoutine int) {
	//
	//	partition := 0
	//	log.Debug("partition:", partition, "start parse local file")
	//	offset := *offsets
	//
	//	broker := kafka.NewBrokerConsumer(kafkaConfig.Hostname, taskConfig.Name+"_parse", partition, offset.Offset, kafkaConfig.MaxSize)
	//
	////	randomPartition := 0
	////	if MaxGoRoutine > 1 {
	////		randomPartition = rand.Intn(MaxGoRoutine - 1)
	////	}
	////	//		log.Debug("random partition:",random)
	////	publisher := kafka.NewBrokerPublisher(kafkaConfig.Hostname, taskConfig.Name+"_fetch", randomPartition)
	//
	//	consumerCallback := func(msg *kafka.Message) {
	//
	//		message := msg.Payload()
	//		stringArray:=strings.Split(string(message),"|||");
	//		fileUrl:=stringArray[0]
	//		fileName:=[]byte(stringArray[1])
	//
	//		if(parseFilter.Lookup(fileName)){
	//			log.Debug("hit parse filter ignore,",string(fileName))
	//			return
	//		}
	//		parseFilter.Add(fileName)
	//
	//		fileContent := loadFileContent(string(fileName))
	//
	//		if fileContent != nil {
	//			log.Debug("partition:", partition, ",parse fileName:", string(fileName))
	//
	//			extractLinks(pendingUrls, bloomFilter,fileUrl, fileName, fileContent, taskConfig)
	//			offsetV := msg.Offset()
	//			offset.Offset = offsetV
	//
	//			path := taskConfig.BaseStoragePath+     "task/parse_offset_" + strconv.FormatInt(int64(partition), 10) + ".tmp"
	//			path_new := taskConfig.BaseStoragePath+     "task/parse_offset_" + strconv.FormatInt(int64(partition), 10)
	//			fout, error := os.Create(path)
	//			if error != nil {
	//				log.Error(path, error)
	//				return
	//			}
	//
	//			defer fout.Close()
	//			log.Debug("partition:", partition, ",saved offset:", offsetV)
	//			fout.Write([]byte(strconv.FormatUint(msg.Offset(), 10)))
	//			utils.CopyFile(path, path_new)
	//		}
	//
	//	}
	//	msgChan := make(chan *kafka.Message)
	//	go broker.ConsumeOnChannel(msgChan, 10, *quit)
	//	for msg := range msgChan {
	//		if msg != nil {
	//			log.Debug("partition:", partition, ",consume messaging,parsing.", string(msg.Payload()))
	//			consumerCallback(msg)
	//		} else {
	//			break
	//		}
	//	}
	//	log.Debug("partition:", partition, ",exit parse local file")
}
