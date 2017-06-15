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

package pipe

import (
	"github.com/PuerkitoBio/purell"
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/errors"
	"github.com/medcl/gopa/core/model"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/util"
	u "net/url"
	"path"
	"sort"
	"strings"
	"time"
)

const UrlNormalization JointKey = "url_normalization"

type UrlNormalizationJoint struct {
	timeout              time.Duration
	splitByUrlParameter  []string
	FollowAllDomain      bool
	FollowDomainSettings bool
	FollowSubDomain      bool
	maxFileNameLength    int
}

var defaultFileName = "default.html"

func (this UrlNormalizationJoint) Name() string {
	return string(UrlNormalization)
}

func (this UrlNormalizationJoint) Process(context *Context) error {

	task := context.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := context.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	url := task.Url
	var currentURI, referenceURI *u.URL
	var err error

	if len(url) <= 0 {
		context.ErrorExit("url can't be null")
	}

	log.Trace("start parse url,", url)

	var tempUrl = url

	//adding default http protocol
	if strings.HasPrefix(url, "//") {
		tempUrl = strings.TrimLeft(url, "//")
		tempUrl = "http:" + url
	}

	currentURI, err = u.Parse(tempUrl)
	if err != nil {
		log.Debug("url parsed failed, ", err, ",", tempUrl)
		context.ErrorExit(err.Error())
	}

	log.Tracef("currentURI,schema:%s, host:%s", currentURI.Scheme, currentURI.Host)
	refUrlStr := task.Reference
	var refExists bool
	if refUrlStr != "" {
		log.Trace("ref url exists, ", refUrlStr)
		referenceURI, err = u.ParseRequestURI(refUrlStr)
		if err == nil {
			refExists = true
		} else {
			log.Warn("ref url parsed failed, ", err)
		}
	}

	//try to fix relative links
	if currentURI == nil || currentURI.Host == "" {

		log.Trace("host is nil, ", url)

		//try to fix link with reference
		if refExists && referenceURI != nil {

			log.Trace("ref is not nil, try to fix relative link: ", url)

			var parentPath = "/"

			if strings.HasPrefix(url, "/") {
				url = "http://" + referenceURI.Host + url
				log.Trace("new relatived url,", url)
			} else {
				var parentUrlFullPath string

				if referenceURI.Path != "" {
					var index = strings.LastIndex(referenceURI.Path, "/")

					if index > 0 {
						parentPath = util.SubString(referenceURI.Path, 0, index)

						if !strings.HasSuffix(parentPath, "/") {
							parentPath = parentPath + "/"
						}
					}
					parentUrlFullPath = "http://" + referenceURI.Host + parentPath
				} else {
					parentUrlFullPath = "http://" + referenceURI.Host
				}

				log.Trace("parent url fullpath:", parentUrlFullPath)

				if strings.HasPrefix(referenceURI.Path, "/") {
					//root based relative urls
					url = parentUrlFullPath + url
					log.Trace("new relatived url,", url)
				} else {
					//page based relative urls
					urlPath := util.GetRootUrl(referenceURI)
					url = "http://" + urlPath + url
					log.Trace("new relatived url,", url)
				}
			}

			log.Trace("fixed link: ", url)
		}

		//try to fix link with host
		if task.Host != "" {
			if strings.HasPrefix(url, "/") {
				url = "http://" + task.Host + url
			} else {
				url = "http://" + task.Host + "/" + url
			}
			log.Trace("new relatived url with host,", url)
		}

		tempUrl = url
		if strings.HasPrefix(url, "//") {
			tempUrl = strings.TrimLeft(url, "//")
		}

		if !strings.HasPrefix(url, "http") {
			tempUrl = "http://" + url
		}
		currentURI, err = u.Parse(tempUrl)
		if err != nil {
			log.Error(err)
			context.Break(err.Error())
			return err
		}
	}

	url = tempUrl

	if strings.Contains(url, "..") {

		url = purell.NormalizeURL(currentURI, purell.FlagsUsuallySafeGreedy|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveFragment)
		//update currentURI
		currentURI, _ = u.Parse(url)
		log.Trace("purell parsed url,", url)
	}

	////resolve domain specific filter
	if !this.FollowAllDomain && this.FollowSubDomain && currentURI != nil && referenceURI != nil {
		log.Tracef("try to check domain rule, %s vs %s", referenceURI.Host, currentURI.Host)
		//TODO handler com.cn and .com,using a TLC-domain list

		if strings.Contains(currentURI.Host, ".") && strings.Contains(referenceURI.Host, ".") {
			ref := strings.Split(referenceURI.Host, ".")
			cur := strings.Split(currentURI.Host, ".")

			log.Tracef("%s vs %s , %s vs %s ", ref[len(ref)-1], cur[len(cur)-1], ref[len(ref)-2], cur[len(cur)-2])

			if !(ref[len(ref)-1] == cur[len(cur)-1] && ref[len(ref)-2] == cur[len(cur)-2]) {
				log.Debug("domain mismatch,", referenceURI.Host, " vs ", currentURI.Host)
				context.Break("domain missmatch," + referenceURI.Host + " vs " + currentURI.Host)
				return nil //known exception, not error
			}
		} else {
			if referenceURI.Host != currentURI.Host {
				context.Break("domain missmatch," + referenceURI.Host + " vs " + currentURI.Host)
				return nil //known exception, not error
			}
		}

	}

	task.Url = url
	task.Host = currentURI.Host
	task.Schema = currentURI.Scheme

	filePath := ""
	filename := ""

	filenamePrefix := ""

	//the url is a folder, making folders
	if strings.HasSuffix(url, "/") {
		filename = defaultFileName
		log.Trace("no page name found,use default.html:", url)
	}

	// if the url have parameters
	if len(currentURI.Query()) > 0 {

		//TODO 不处理非网页内容，去除js 图片 css 压缩包等

		if len(this.splitByUrlParameter) > 0 {

			for i := 0; i < len(this.splitByUrlParameter); i++ {
				breakTagTemp := currentURI.Query().Get(this.splitByUrlParameter[i])
				if breakTagTemp != "" {
					filenamePrefix = filenamePrefix + this.splitByUrlParameter[i] + "_" + breakTagTemp + "_"
				}
			}
		} else {
			queryMap := currentURI.Query()

			//sort the parameters by parameter key
			keys := make([]string, 0, len(queryMap))
			for key := range queryMap {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, key := range keys {
				value := queryMap[key]
				len := len(value)
				if value != nil && len > 0 {
					if len > 0 {
						filenamePrefix = filenamePrefix + key + "_"
						for i := 0; i < len; i++ {
							v := value[i]
							if v != "" && len > 0 {
								filenamePrefix = (filenamePrefix + v + "_")
							}
						}
					}
				}
			}
			filenamePrefix = strings.TrimRight(filenamePrefix, "_")
		}
	}

	//split folder and filename and also insert the prefix filename
	index := strings.LastIndex(currentURI.Path, "/")

	if index > 0 {
		//the url should has at least one folder
		//http://xx.com/1112/12
		filePath = currentURI.Path[0:index]
		log.Trace("filepath: ", filePath)
		//if the page extension is missing
		if !strings.Contains(currentURI.Path, ".") {
			if strings.HasSuffix(currentURI.Path, "/") {
				filename = currentURI.Path[index:len(currentURI.Path)] + defaultFileName
			} else {
				filename = currentURI.Path[index:len(currentURI.Path)] + ".html"
			}
		} else {
			filename = currentURI.Path[index:len(currentURI.Path)]
		}
	} else {
		//file in the root folder
		log.Tracef("no / in path, %s", currentURI.Path)
		filePath = ""
		filename = currentURI.Path
	}

	if len(filenamePrefix) > 0 {
		log.Tracef("get file prefix: %s", filenamePrefix)
		if strings.Contains(filename, "/") {
			log.Tracef("filename contains / : %s", filename)
			index := strings.LastIndex(filename, "/") + 1
			start := filename[0:index]
			end := filename[index:]
			log.Tracef("filename start: %s, end: %s", start, end)

			if strings.Contains(end, ".") {
				index := strings.LastIndex(end, ".")
				start1 := end[0:index]
				end1 := end[index:]
				filename = start + start1 + "_" + filenamePrefix + end1
			} else {
				filename = start + "_" + filenamePrefix + end
			}

		} else {
			filename = filenamePrefix + "_" + filename
		}
	}

	if len(filename) == 0 {
		filename = defaultFileName
	}

	if !strings.Contains(filename, ".") {
		filePath = path.Join(filePath, filename)
		filename = defaultFileName
	}

	//set default filename limit
	if this.maxFileNameLength <= 0 {
		this.maxFileNameLength = 200
	}

	//verify filename
	if len(filename) > this.maxFileNameLength {
		panic(errors.New("file name too long"))
	}

	snapshot.Path = filePath
	snapshot.File = filename
	log.Debugf("finished normalization,%s, %s, %s, %s ",task.ID, url, filePath, filename)

	return nil
}
