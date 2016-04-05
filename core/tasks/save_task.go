/**
 * User: Medcl
 * Date: 13-7-25
 * Time: 下午9:19
 */
package tasks

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	. "net/url"
	"os"
	"strings"
)

func init() {
}

func getSavedPath(runtimeConfig *RuntimeConfig, url []byte) string {

	siteConfig := runtimeConfig.TaskConfig

	urlStr := string(url)
	log.Debug("start saving url,", urlStr)
	myurl1, _ := Parse(urlStr)
	log.Trace("url->path:", myurl1.Host, " ", myurl1.Path)

	baseDir := runtimeConfig.TaskConfig.WebDataPath + "/" + myurl1.Host + "/"
	baseDir = strings.Replace(baseDir, `:`, `_`, -1)

	log.Trace("replaced:", baseDir)
	path := ""
	filename := ""

	//the url is a folder, making folders
	if strings.HasSuffix(urlStr, "/") {
		path = baseDir + myurl1.Path
		os.MkdirAll(path, 0777)
		log.Trace("making dir:", path)
		filename = "default.html"
		log.Trace("no page name,use default.html:", path)

	}

	filenamePrefix := ""

	if len(myurl1.Query()) > 0 {

		//TODO 不处理非网页内容，去除js 图片 css 压缩包等

		if siteConfig.SplitByUrlParameter != "" {

			arrayStr := strings.Split(siteConfig.SplitByUrlParameter, runtimeConfig.ArrayStringSplitter)
			for i := 0; i < len(arrayStr); i++ {
				breakTagTemp := myurl1.Query().Get(arrayStr[i])
				if breakTagTemp != "" {
					filenamePrefix = filenamePrefix + arrayStr[i] + "_" + breakTagTemp + "_"
				}
			}
		} else {
			queryMap := myurl1.Query()
			//			queryMap = sort.Sort(queryMap) //TODO sort the parameters by parameter key
			for key, value := range queryMap {
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
	index := strings.LastIndex(myurl1.Path, "/")
	if index > 0 {
		//http://xx.com/1112/12
		path = myurl1.Path[0:index]
		path = baseDir + path
		os.MkdirAll(path, 0777)

		//if the page extension is missing
		if !strings.Contains(myurl1.Path, ".") {
			filename = myurl1.Path[index:len(myurl1.Path)] + ".html"
		} else {
			filename = myurl1.Path[index:len(myurl1.Path)]
		}
	} else {
		path = baseDir + path + "/"
		os.MkdirAll(path, 0777)
		filename = "default.html"
	}

	filename = strings.Replace(filename, "/", "", -1)

	path = path + "/" + filenamePrefix + filename

	log.Trace(urlStr, " will save to file: ", path)

	return path
}

func Save(runtimeConfig *RuntimeConfig, path string, body []byte) (int, error) {

	log.Trace("saving file,", path)
	fout, error := os.Create(path)
	if error != nil {
		log.Error(path, error)
		return 5, error
	}

	defer fout.Close()
	rt, err := fout.Write(body)
	return rt, err
}
