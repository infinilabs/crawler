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
	"kafka"
	. "types"
	bloom "github.com/zeebo/sbloom"
	"hash/fnv"
)
var saveFilter  *bloom.Filter
func init() {
	saveFilter = bloom.NewFilter(fnv.New64(), 1000000)
}

func Save(siteConfig *TaskConfig,myurl []byte, body []byte, publisher *kafka.BrokerPublisher) {
	if(saveFilter.Lookup(myurl)){
		log.Debug("hit save filter ignore,",string(myurl))
		return
	}
	saveFilter.Add(myurl)

	urlStr := string(myurl)
	log.Info("start saving url,", urlStr)
	myurl1, _ := Parse(urlStr)
	log.Debug("url->path:", myurl1.Host, " ", myurl1.Path)

	baseDir :=  myurl1.Host + "/"
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
		log.Debug("index of last /:", index, ",", myurl1.Path)
		if index >= 0 {
			path = myurl1.Path[0:index]
			path = baseDir + path
			log.Debug("new path:", path)
			os.MkdirAll(path, 0777)
			log.Debug("making dir:", path)
			path = (baseDir + myurl1.Path)
			log.Trace("fileUrl:",urlStr)
//			myurl1.Query().Encode();
//			log.Error("fileArgs:",myurl1.Query().Get("pn"))
//			log.Error("fileArgs:",myurl1.Query().Get("p"))

//			log.Error("fileArgs:",myurl1.Path)
			log.Trace("fileArgs:",myurl1.RawQuery)
			//check to see if we have paging info，TODO configable
//			if strings.Contains(myurl1.RawQuery, "p"){
//		      	getParameters:=strings.Split(myurl1.RawQuery, "&")
//				strings.
//			}



//			log.Error("fileArgs:",myurl1.Query().Encode())
			log.Trace("fileName:",path)
		} else {
			path = baseDir + path + "/"
			os.MkdirAll(path, 0777)
			log.Debug("making dir:", path)
			path = path + "/default.html"
		}
	}




	if siteConfig.SplitByUrlParameter!=""{

		breakTag:=myurl1.Query().Get(siteConfig.SplitByUrlParameter)
		if breakTag!="" {
			log.Debug("url with page parameter")
			path=path+"_"+breakTag+".html"
		}
	}


	log.Debug("touch file,", path)
	fout, error := os.Create(path)
	if error != nil {
		log.Error(path, error)
		return
	}

	defer fout.Close()
	log.Info("saved:", urlStr, ",", path)
	fout.Write(body)

	publisher.Publish(kafka.NewMessage([]byte(path)))

	//	log.Info("enqueue parse,", path)
	log.Debug("end saving url,", urlStr)

}
