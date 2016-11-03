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
	"errors"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/util"
	. "net/url"
	"sort"
	"strings"
	"time"
)

type UrlNormalizationJoint struct {
	timeout             time.Duration
	splitByUrlParameter []string
	FollowSubDomain     bool
}

var defaultFileName = "default.html"

func (this UrlNormalizationJoint) Name() string {
	return "url_normalization"
}

func (this UrlNormalizationJoint) Process(context *Context) (*Context, error) {
	url := context.MustGetString(CONTEXT_URL)
	var currentURI, referenceURI *URL
	var err error

	log.Trace("start parse url,", url)

	var tempUrl = url
	//adding default http protocol
	if strings.HasPrefix(url, "//") {
		tempUrl = strings.TrimLeft(url, "//")
		tempUrl = "http:" + url
	}

	currentURI, _ = Parse(tempUrl)

	log.Tracef("currentURI,schema:%s, host:%s", currentURI.Scheme, currentURI.Host)
	refUrlStr, refExists := context.GetString(CONTEXT_REFERENCE_URL)
	if refExists && refUrlStr != "" {
		log.Trace("ref url exists, ", refUrlStr)
		referenceURI, err = ParseRequestURI(refUrlStr)
		if err != nil {
			log.Trace("ref url parsed failed, ", err)
		}
	}

	//try to fix relative links
	if currentURI == nil || currentURI.Host == "" {

		log.Trace("host is nil, ", url)

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

		tempUrl = url
		if strings.HasPrefix(url, "//") {
			tempUrl = strings.TrimLeft(url, "//")
		}

		if !strings.HasPrefix(url, "http") {
			tempUrl = "http://" + url
		}
		currentURI, err = Parse(tempUrl)
		if err != nil {
			log.Error(err)
			context.Break()
			return context, err
		}
	}

	////resolve domain specific filter
	if this.FollowSubDomain && currentURI != nil && referenceURI != nil {
		log.Tracef("try to check domain rule, %s vs %s", referenceURI.Host, currentURI.Host)
		//TODO handler com.cn and .com,using a TLC-domain list

		if strings.Contains(currentURI.Host, ".") && strings.Contains(referenceURI.Host, ".") {
			ref := strings.Split(referenceURI.Host, ".")
			cur := strings.Split(currentURI.Host, ".")

			log.Tracef("%s vs %s , %s vs %s ", ref[len(ref)-1], cur[len(cur)-1], ref[len(ref)-2], cur[len(cur)-2])

			if !(ref[len(ref)-1] == cur[len(cur)-1] && ref[len(ref)-2] == cur[len(cur)-2]) {
				log.Debug("domain mismatch,", referenceURI.Host, " vs ", currentURI.Host)
				context.Break()
				return context, errors.New("domain mismatch")
			}
		} else {
			if referenceURI.Host != currentURI.Host {
				context.Break()
				return context, errors.New("domain mismatch")
			}
		}

	}

	url = tempUrl
	context.Set(CONTEXT_URL, url)
	context.Set(CONTEXT_HOST, currentURI.Host)
	context.Set(CONTEXT_URL_PATH, currentURI.RawPath)

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
				if value != nil && len(value) > 0 {
					if len(value) > 0 {
						filenamePrefix = filenamePrefix + key + "_"
						for i := 0; i < len(value); i++ {
							v := value[i]
							if v != "" && len(v) > 0 {
								filenamePrefix = filenamePrefix + v + "_"
							}
						}
					}

				}
			}
		}
	}

	//split folder and filename and also insert the prefix filename
	index := strings.LastIndex(currentURI.Path, "/")

	if index > 0 {
		//the url should has at least one folder
		//http://xx.com/1112/12
		filePath = currentURI.Path[0:index]

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

	filename = filenamePrefix + filename
	context.Set(CONTEXT_SAVE_PATH, filePath)
	context.Set(CONTEXT_SAVE_FILENAME, filename)
	log.Debugf("finished normalization,%s, %s, %s ",url,filePath,filename)

	return context, nil
}
