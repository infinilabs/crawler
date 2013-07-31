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
	. "github.com/PuerkitoBio/purell"
	. "github.com/zeebo/sbloom"
	"kafka"
	util "util"
	//config	"config"
	//	"strconv"
	//	utils "util"
	. "types"
)

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

//format url,prepare for bloom filter
func formatUrlForFilter(url []byte) []byte {
	src := string(url)
	log.Debug("start to normalize url:", src)
	if strings.HasSuffix(src, "/") {
		src = strings.TrimRight(src, "/")
	}
	src = strings.TrimSpace(src)
	src = strings.ToLower(src)
	return []byte(src)
}

//func ExtractLinksFromTaskResponse(bloomFilter *Filter,curl chan []byte, task Task, siteConfig *SiteConfig) {
func ExtractLinksFromTaskResponse(bloomFilter *Filter, broker *kafka.BrokerPublisher, task Task, siteConfig *TaskConfig) {
	siteUrlStr := string(task.Url)
	log.Debug("enter links extract,", siteUrlStr)
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

		if bloomFilter.Lookup(filterUrl) {
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
			currentUrlStr = MustNormalizeURLString(currentUrlStr, FlagLowercaseScheme|FlagLowercaseHost|FlagUppercaseEscapes|
				FlagRemoveUnnecessaryHostDots|FlagRemoveDuplicateSlashes|FlagRemoveFragment)
			log.Debug("normalized url:", currentUrlStr)
			currentUrlByte := []byte(currentUrlStr)
			if !bloomFilter.Lookup(currentUrlByte) {

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
		log.Debug("exit links extract,", siteUrlStr)

	}

	log.Info("all links within ", siteUrlStr, " is done")
}
