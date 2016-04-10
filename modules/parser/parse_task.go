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

//subscribe local file channel,and parse urls
package parser

import (
	"bufio"
	"io/ioutil"
	. "net/url"
	"os"
	"strings"
	"time"

	. "github.com/PuerkitoBio/purell"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/util"
)

func init() {
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

func extractLinks(runtimeConfig *RuntimeConfig, fileUrl string, fileName []byte, body []byte) {

	var storage = runtimeConfig.Storage

	siteUrlStr := fileUrl
	siteConfig := runtimeConfig.TaskConfig

	siteUrlByte := []byte(siteUrlStr)
	log.Debug("enter links extract,", siteUrlStr)
	if siteConfig.SkipPageParsePattern.Match(siteUrlByte) {
		log.Debug("hit SkipPageParsePattern pattern,", siteUrlStr)
		return
	}

	//parse parent url
	seedUrlStr := siteUrlStr
	seedURI, err := ParseRequestURI(seedUrlStr)

	if err != nil {
		log.Error("ParseSeedURI failed!: ", seedUrlStr, " , ", err)
		return
	}

	var parentPath = "/"

	var parentUrlFullPath string

	if seedURI.Path != "" {
		var index = strings.LastIndex(seedURI.Path, "/")

		if index > 0 {
			parentPath = util.SubString(seedURI.Path, 0, index)

			if !strings.HasSuffix(parentPath, "/") {
				parentPath = parentPath + "/"
			}
		}
		parentUrlFullPath = "http://" + seedURI.Host + parentPath
	} else {
		parentUrlFullPath = "http://" + seedURI.Host
	}

	log.Trace("parent url fullpath:", parentUrlFullPath)

	log.Debug("parsing external links:", siteUrlStr, ",using:", siteConfig.LinkUrlExtractRegex)

	matches := siteConfig.LinkUrlExtractRegex.FindAllSubmatch(body, -1)
	log.Debug("extract links with pattern,total matchs:", len(matches), " match result,", string(fileName))
	xIndex := 0
	for _, match := range matches {
		log.Trace("dealing with match result,", xIndex)
		xIndex = xIndex + 1
		url := match[siteConfig.LinkUrlExtractRegexGroupIndex]
		filterUrl := util.FormatUrlForFilter(url)
		log.Debug("url clean result:", string(filterUrl), ",original url:", string(url))
		filteredUrl := string(filterUrl)

		//filter error link
		if filteredUrl == "" {
			log.Trace("filteredUrl is empty,continue")
			continue
		}

		result1 := strings.HasPrefix(filteredUrl, "#")
		if result1 {
			log.Trace("filteredUrl started with: # ,continue")
			continue
		}

		result2 := strings.HasPrefix(filteredUrl, "javascript:")
		if result2 {
			log.Trace("filteredUrl started with: javascript: ,continue")
			continue
		}

		hit := false

		if storage.UrlHasWalked(filterUrl) || storage.UrlHasFetched(filterUrl) || storage.PendingFetchUrlHasAdded(filterUrl) {
			log.Trace("hit Filter,continue")
			hit = true
			continue
		}

		if !hit {
			currentUrlStr := string(url)
			currentUrlStr = strings.Trim(currentUrlStr, " ")

			currentURI1, err := ParseRequestURI(currentUrlStr)
			currentURI := currentURI1
			if err != nil {

				log.Trace("invalid url,", err)

				if strings.Contains(err.Error(), "invalid URI for request") {
					log.Debug("invalid URI for request,fix relative url,original:", currentUrlStr)

					//page based relative urls
					currentUrlStr = parentUrlFullPath + currentUrlStr
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
					currentUrlStr = parentUrlFullPath + currentUrlStr
					log.Trace("new relatived url,", currentUrlStr)
				} else {
					//page based relative urls
					urlPath := util.GetRootUrl(currentURI)
					currentUrlStr = "http://" + urlPath + currentUrlStr
					log.Trace("new relatived url,", currentUrlStr)
				}
			} else {
				//resolve domain specific filter
				if siteConfig.FollowSameDomain {
					if siteConfig.FollowSubDomain {

						//TODO handler com.cn and .com,using a TLC-domain list

					}

					if seedURI.Host != currentURI.Host {
						log.Debug("domain mismatch,", seedURI.Host, " vs ", currentURI.Host)
						continue
					}
					//TODO follow all or list of domain
				}
			}

			if len(siteConfig.LinkUrlMustContain) > 0 {
				if !util.ContainStr(currentUrlStr, siteConfig.LinkUrlMustContain) {
					log.Trace("link does not hit must-contain,ignore,", currentUrlStr, " , ", siteConfig.LinkUrlMustNotContain)
					continue
				}
			}

			if len(siteConfig.LinkUrlMustNotContain) > 0 {
				if util.ContainStr(currentUrlStr, siteConfig.LinkUrlMustNotContain) {
					log.Trace("link hit must-not-contain,ignore,", currentUrlStr, " , ", siteConfig.LinkUrlMustNotContain)
					continue
				}
			}

			//normalize url
			currentUrlStr = MustNormalizeURLString(currentUrlStr, FlagLowercaseScheme|FlagLowercaseHost|FlagUppercaseEscapes|
				FlagRemoveUnnecessaryHostDots|FlagRemoveDuplicateSlashes|FlagRemoveFragment)
			log.Trace("normalized url:", currentUrlStr)
			currentUrlByte := []byte(currentUrlStr)
			if !(storage.UrlHasWalked(currentUrlByte) || storage.UrlHasFetched(currentUrlByte) || storage.PendingFetchUrlHasAdded(currentUrlByte)) {

				//copied form fetchTask,TODO refactor
				//checking fetchUrlPattern
				log.Trace("started check fetchUrlPattern,", currentUrlStr)
				if siteConfig.FetchUrlPattern.Match(currentUrlByte) {
					log.Trace("match fetch url pattern,", currentUrlStr)
					if len(siteConfig.FetchUrlMustNotContain) > 0 {
						if util.ContainStr(currentUrlStr, siteConfig.FetchUrlMustNotContain) {
							log.Trace("hit FetchUrlMustNotContain,ignore,", currentUrlStr)
							continue
						}
					}

					if len(siteConfig.FetchUrlMustContain) > 0 {
						if !util.ContainStr(currentUrlStr, siteConfig.FetchUrlMustContain) {
							log.Trace("not hit FetchUrlMustContain,ignore,", currentUrlStr)
							continue
						}
					}
				} else {
					log.Trace("does not hit FetchUrlPattern ignoring,", currentUrlStr)
					continue
				}

				if !storage.PendingFetchUrlHasAdded(currentUrlByte) {
					log.Trace("log new pendingFetch url,", currentUrlStr)
					storage.AddPendingFetchUrl(currentUrlByte)
					storage.LogPendingFetchUrl(runtimeConfig.PathConfig.PendingFetchLog, currentUrlStr)
					log.Debug("check filter result:", currentUrlStr, ":", storage.PendingFetchUrlHasAdded(currentUrlByte))

				} else {
					log.Error("hit new pendingFetch filter,ignore:", currentUrlStr)
					continue
				}

				//	TODO pendingFetchFilter			bloomFilter.Add(currentUrlByte)
			} else {
				log.Trace("hit filter,ignore:", currentUrlStr)
			}
		} else {
			log.Trace("hit filter,ignore,", string(url))
		}
		log.Trace("exit links extract,", siteUrlStr)

	}

	//TODO 处理ruled fetch pattern

	log.Debug("all links within ", siteUrlStr, " is done")
}

func ParseGo(pendingUrls chan []byte, runtimeConfig *RuntimeConfig, quit *chan bool) {
	log.Info("parsing task started.")
	path := runtimeConfig.PathConfig.SavedFileLog
	//touch local's file
	//read all of line
	//if hit the EOF,will wait 2s,and then reopen the file,and try again,may be check the time of last modified

waitFile:
	if !util.CheckFileExists(path) {
		log.Trace("waiting file create:", path)
		time.Sleep(1000 * time.Millisecond)
		goto waitFile
	}

	var storage = runtimeConfig.Storage
	var offset int64 = storage.LoadOffset(runtimeConfig.PathConfig.SavedFileLog + ".offset")
	log.Info("loaded parse offset:", offset)
	FetchFileWithOffset(runtimeConfig, path, offset)
}

func FetchFileWithOffset(runtimeConfig *RuntimeConfig, path string, skipOffset int64) {

	var offset int64 = 0

	time1, _ := util.FileMTime(path)
	log.Trace("start touch time:", time1)

	f, err := os.Open(path)
	if err != nil {
		log.Debug("error opening file,", path, " ", err)
		return
	}

	var storage = runtimeConfig.Storage

	r := bufio.NewReader(f)
	s, e := util.Readln(r)
	offset = 0

	for e == nil {
		offset = offset + 1
		//TODO use byte offset instead of lines
		if offset > skipOffset {
			ParsedSavedFileLog(runtimeConfig, s)
			storage.PersistOffset(runtimeConfig.PathConfig.SavedFileLog+".offset", offset)
		}

		s, e = util.Readln(r)
		//todo store offset
	}
	log.Trace("end offset:", offset, "vs ", skipOffset)

waitUpdate:
	time2, _ := util.FileMTime(path)

	log.Trace("2nd touch time:", time2)

	if time2 > time1 {
		log.Debug("file has been changed,restart parse")
		FetchFileWithOffset(runtimeConfig, path, offset)
	} else {
		log.Trace("waiting file update:", path)
		time.Sleep(5 * time.Second)
		goto waitUpdate
	}
}

func ParsedSavedFileLog(runtimeConfig *RuntimeConfig, fileLog string) {
	if fileLog != "" {
		var storage = runtimeConfig.Storage
		log.Debug("start parse filelog:", fileLog)
		//load file's content,and extract links

		stringArray := strings.Split(fileLog, "|||")
		fileUrl := stringArray[0]
		fileName := []byte(stringArray[1])

		if storage.FileHasParsed(fileName) {
			log.Debug("hit parse filter ignore,", string(fileName))
			return
		}

		fileContent := loadFileContent(string(fileName))
		storage.AddParsedFile(fileName)

		if fileContent != nil {

			//extract urls to fetch queue.
			extractLinks(runtimeConfig, fileUrl, fileName, fileContent)
		}
	}
}
