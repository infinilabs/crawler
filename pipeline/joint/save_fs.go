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

package joint

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/model"
	"net/url"
	"os"
	"path"
	"strings"
)

type SaveSnapshotToFileSystemJoint struct {
	baseDir string
}

func (joint SaveSnapshotToFileSystemJoint) Name() string {
	return "save_snapshot_fs"
}

func (joint SaveSnapshotToFileSystemJoint) Process(c *pipeline.Context) error {

	if len(joint.baseDir) == 0 {
		joint.baseDir = global.Env().SystemConfig.GetWorkingDir() + "/web"
	}

	snapshot := c.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	url := c.MustGetString(model.CONTEXT_TASK_URL)
	host := c.MustGetString(model.CONTEXT_TASK_Host)
	dir := snapshot.Path
	file := snapshot.File
	folder := path.Join(joint.baseDir, host, dir)

	fullPath := path.Join(folder, file)

	if util.FileExists(fullPath) {
		log.Warnf("file: %s already exists, ignore,url: %s", fullPath, url)
		return nil
	}

	log.Trace("save url,", url, ",host,", host, ",folder,", folder, ",file:", file, ",fullpath,", fullPath)

	err := os.MkdirAll(folder, 0777)
	if err != nil {
		log.Error(fullPath, ",", err)
		panic(err)
	}

	fout, err := os.Create(fullPath)
	if err != nil {
		log.Error(fullPath, ",", err)
		panic(err)
	}

	defer fout.Close()
	_, err = fout.Write(snapshot.Payload)
	fout.Sync()
	if err != nil {
		log.Error(fullPath, ",", err)
		panic(err)
	}

	return nil
}

func (joint SaveSnapshotToFileSystemJoint) getSavedPath(urlStr string) (string, string) {

	log.Debug("start saving url,", urlStr)
	myurl1, _ := url.Parse(urlStr)

	baseDir := path.Join(joint.baseDir, myurl1.Host)
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
