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
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/model"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/util"
	 "net/url"
	"os"
	"path"
	"strings"
)

const SaveToFileSystem JointKey = "save2fs"

type SaveToFileSystemJoint struct {
	context *Context
	baseDir string
}

func (this SaveToFileSystemJoint) Name() string {
	return string(SaveToFileSystem)
}

func (this SaveToFileSystemJoint) Process(c *Context) (*Context, error) {
	this.context = c

	if len(this.baseDir) == 0 {
		this.baseDir = global.Env().SystemConfig.Data + "/web"
	}

	url, ok := c.GetString(CONTEXT_URL)
	if !ok {
		return nil, errors.New("invalid url")
	}
	task := c.Get(CONTEXT_PAGE_ITEM).(*model.Task)
	pageItem := c.Get(CONTEXT_PAGE_ITEM).(*model.PageItem)

	domain := c.MustGetString(CONTEXT_HOST)
	dir := c.MustGetString(CONTEXT_SAVE_PATH)
	file := c.MustGetString(CONTEXT_SAVE_FILENAME)
	folder := path.Join(this.baseDir, domain, dir)

	os.MkdirAll(folder, 0777)

	fullPath := path.Join(folder, file)

	if util.FileExists(fullPath) {
		log.Warnf("file: %s already exists, ignore,url: %s", fullPath, url)
		return c, nil
	}

	log.Trace("save url,", url, ",domain,", task.Domain, ",fullpath,", fullPath)

	fout, err := os.Create(fullPath)
	if err != nil {
		log.Error(fullPath, err)
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
	myurl1, _ := url.Parse(urlStr)

	baseDir := path.Join(this.baseDir, myurl1.Host)
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

	//split folder and filename and also insert the prefix filename
	index := strings.LastIndex(myurl1.Path, "/")
	if index > 0 {
		//http://xx.com/1112/12
		filePath = myurl1.Path[0:index]
		filePath = path.Join(baseDir, filePath)

		//if the page extension is missing
		if !strings.Contains(myurl1.Path, ".") {
			filename = myurl1.Path[index:len(myurl1.Path)] + ".html"
		} else {
			filename = myurl1.Path[index:len(myurl1.Path)]
		}
	} else {
		filePath = path.Join(baseDir, filePath)
		filename = "default.html"
	}

	filename = strings.Replace(filename, "/", "", -1)

	return filePath, filenamePrefix + filename
}
