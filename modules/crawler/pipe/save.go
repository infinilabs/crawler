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
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/types"
	"github.com/syndtr/goleveldb/leveldb/errors"
	. "net/url"
	"os"
	"strings"
	"path"
	"github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/util"
)

type SaveToFileSystemJoint struct {
	context *Context
	baseDir string
}

func (this SaveToFileSystemJoint) Process(c *Context) (*Context, error) {
	this.context = c

	if(len(this.baseDir)==0){
		this.baseDir =this.context.Env.RuntimeConfig.PathConfig.WebData
	}

	url, ok := c.GetString(CONTEXT_URL)
	if !ok {
		return nil, errors.New("invalid url")
	}
	pageItem := c.Get(CONTEXT_PAGE_ITEM).(*types.PageItem)

	log.Debug("save url,", url, ",domain,", pageItem.Domain)

	dir, file := this.getSavedPath(url)
	fullPath:=path.Join(dir,file)
	log.Trace("saving file,", fullPath)
	os.MkdirAll(dir, 0777)
	fout, err := os.Create(fullPath)
	if err != nil {
		log.Error(dir, err)
		return nil, err
	}

	defer fout.Close()
	_, err = fout.Write(pageItem.Body)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (this SaveToFileSystemJoint) getSavedPath(urlStr string) (string, string) {

	log.Debug("start saving url,", urlStr)
	myurl1, _ := Parse(urlStr)

	baseDir := path.Join(this.baseDir,myurl1.Host)
	baseDir = strings.Replace(baseDir, `:`, `_`, -1)

	log.Trace("replaced:", baseDir)
	filePath := ""
	filename := ""

	filenamePrefix := ""

	//the url is a folder, making folders
	if strings.HasSuffix(urlStr, "/") {
		filename = "default.html"
		log.Trace("no page name,use default.html:", urlStr)
	}

	// if the url have parameters
	if len(myurl1.Query()) > 0 {

		////TODO 不处理非网页内容，去除js 图片 css 压缩包等
		//
		//if siteConfig.SplitByUrlParameter != "" {
		//
		//	arrayStr := strings.Split(siteConfig.SplitByUrlParameter, this.context.Env.RuntimeConfig.ArrayStringSplitter)
		//	for i := 0; i < len(arrayStr); i++ {
		//		breakTagTemp := myurl1.Query().Get(arrayStr[i])
		//		if breakTagTemp != "" {
		//			filenamePrefix = filenamePrefix + arrayStr[i] + "_" + breakTagTemp + "_"
		//		}
		//	}
		//} else {
		//	queryMap := myurl1.Query()
		//	//			queryMap = sort.Sort(queryMap) //TODO sort the parameters by parameter key
		//	for key, value := range queryMap {
		//		if value != nil && len(value) > 0 {
		//			if len(value) > 0 {
		//				filenamePrefix = filenamePrefix + key + "_"
		//				for i := 0; i < len(value); i++ {
		//					v := value[i]
		//					if v != "" && len(v) > 0 {
		//						filenamePrefix = filenamePrefix + v + "_"
		//					}
		//				}
		//			}
		//
		//		}
		//	}
		//}
	}

	//split folder and filename and also insert the prefix filename
	index := strings.LastIndex(myurl1.Path, "/")
	if index > 0 {
		//http://xx.com/1112/12
		filePath = myurl1.Path[0:index]
		filePath = path.Join(baseDir ,filePath)

		//if the page extension is missing
		if !strings.Contains(myurl1.Path, ".") {
			filename = myurl1.Path[index:len(myurl1.Path)] + ".html"
		} else {
			filename = myurl1.Path[index:len(myurl1.Path)]
		}
	} else {
		filePath = path.Join(baseDir , filePath)
		filename = "default.html"
	}

	filename = strings.Replace(filename, "/", "", -1)

	return filePath , filenamePrefix + filename
}

func checkIfUrlWillBeSave(taskConfig *config.TaskConfig,url []byte,)bool  {

	requestUrl:=string(url)

	log.Debug("started check savingUrlPattern,", taskConfig.SavingUrlPattern, ",", string(url))
	if taskConfig.SavingUrlPattern.Match(url) {


		log.Debug("match saving url pattern,", requestUrl)
		if len(taskConfig.SavingUrlMustNotContain) > 0 {
			if util.ContainStr(requestUrl, taskConfig.SavingUrlMustNotContain) {
				log.Debug("hit SavingUrlMustNotContain,ignore,", requestUrl, " , ", taskConfig.SavingUrlMustNotContain)
				return false
			}
		}

		if len(taskConfig.SavingUrlMustContain) > 0 {
			if !util.ContainStr(requestUrl, taskConfig.SavingUrlMustContain) {
				log.Debug("not hit SavingUrlMustContain,ignore,", requestUrl, " , ", taskConfig.SavingUrlMustContain)
				return false
			}
		}

		return true

	} else {
		log.Debug("does not hit SavingUrlPattern ignoring,", requestUrl)
	}
	return false
}
