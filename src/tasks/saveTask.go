/** 
 * User: Medcl
 * Date: 13-7-25
 * Time: 下午9:19 
 */
package tasks

import (
	log "github.com/cihub/seelog"
	. "net/url"
	"os"
	"strings"
	//	. "github.com/zeebo/sbloom"
//	. "github.com/PuerkitoBio/purell"
	"kafka"
)

func Save(myurl []byte, body []byte,publisher *kafka.BrokerPublisher) {
	urlStr := string(myurl)
	log.Debug("start saving url,", urlStr)
	myurl1, _ := ParseRequestURI(urlStr)
	log.Debug("url->path:", myurl1.Host, " ", myurl1.Path)

	baseDir := "data/" + myurl1.Host + "/"
	baseDir = strings.Replace(baseDir, `:`, `_`, -1)

	log.Debug("replaced:", baseDir)
	path := ""

	//making folders
	if strings.HasSuffix(urlStr, "/") {
		path = baseDir + myurl1.Path
		os.MkdirAll(path, 0777)
		log.Debug("making dir:", path)
		path = (path + "default.html")
		log.Debug("no page name,use default.html:", path)

	} else {
		index := strings.LastIndex(myurl1.Path, "/")
		log.Debug("index of last /:", index,",",myurl1.Path)
		if index >= 0 {
			path = myurl1.Path[0:index]
			path = baseDir + path
			log.Debug("new path:", path)
			os.MkdirAll(path, 0777)
			log.Debug("making dir:", path)
			path = (baseDir + myurl1.Path)
		} else {
			path = baseDir + path + "/"
			os.MkdirAll(path, 0777)
			log.Debug("making dir:", path)
			path = path + "/default.html"
		}
	}

	log.Debug("touch file,", path)
	fout, error := os.Create(path)
	if error != nil {
		log.Error(path, error)
		return
	}

	defer fout.Close()
	log.Info("saved:", urlStr,",", path)
	fout.Write(body)

	publisher.Publish(kafka.NewMessage([]byte(path)))

//	log.Info("enqueue parse,", path)
	log.Debug("end saving url,", urlStr)

}


