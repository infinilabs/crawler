/** 
 * User: Medcl
 * Date: 13-7-25
 * Time: 下午9:19 
 */
package tasks

import (
	log "logging"
	. "net/url"
	"os"
	"strings"
//	bloom "github.com/zeebo/sbloom"
//	"hash/fnv"
	. "types"

)
//var saveFilter  *bloom.Filter
func init() {
//	saveFilter = bloom.NewFilter(fnv.New64(), 1000000)
//	log.Warn("init bloom filter")
}


func getSavedPath(runtimeConfig RuntimeConfig,url []byte) string{

	siteConfig:=runtimeConfig.TaskConfig

	//	if(saveFilter.Lookup(myurl)){
	//		log.Debug("hit save filter ignore,",string(myurl))
	//		return
	//	}
	//	saveFilter.Add(myurl)

	urlStr := string(url)
	log.Debug("start saving url,", urlStr)
	myurl1, _ := Parse(urlStr)
	log.Trace("url->path:", myurl1.Host, " ", myurl1.Path)

	baseDir := runtimeConfig.TaskConfig.WebDataPath+"/"+ myurl1.Host + "/"
	baseDir = strings.Replace(baseDir, `:`, `_`, -1)

	log.Trace("replaced:", baseDir)
	path := ""

	//making folders
	if strings.HasSuffix(urlStr, "/") {
		path = baseDir + myurl1.Path
		os.MkdirAll(path, 0777)
		log.Trace("making dir:", path)
		path = (path + "default.html")
		log.Trace("no page name,use default.html:", path)

	} else {
		index := strings.LastIndex(myurl1.Path, "/")
		log.Trace("index of last /:", index, ",", myurl1.Path)
		if index >= 0 {
			path = myurl1.Path[0:index]
			path = baseDir + path
			log.Trace("new path:", path)
			os.MkdirAll(path, 0777)
			log.Trace("making dir:", path)
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
			log.Trace("making dir:", path)
			path = path + "/default.html"
		}
	}




	if siteConfig.SplitByUrlParameter!=""{

		arrayStr:=strings.Split(siteConfig.SplitByUrlParameter,runtimeConfig.ArrayStringSplitter)
		breakTag:=""
		for i := 0; i < len(arrayStr); i++ {
			breakTagTemp:=myurl1.Query().Get(arrayStr[i])
			if breakTagTemp!="" {
				log.Trace("url with page parameter")
				breakTag=(breakTag+"_"+breakTagTemp)
			}
		}
		if(breakTag!=""){
			path=(path+breakTag+".html")
		}
	}

  return path
}


//func Save(siteConfig *TaskConfig,myurl []byte, body []byte, publisher *kafka.BrokerPublisher) {
func Save(path string,body []byte) {


//	pathArray:=[]byte(path)

//	if(saveFilter.Lookup(pathArray)){
//		log.Debug("hit save-path filter ignore,",string(path))
//		return
//	}
//	saveFilter.Add(pathArray)


	log.Trace("touch file,", path)
	fout, error := os.Create(path)
	if error != nil {
		log.Error(path, error)
		return
	}

	defer fout.Close()
	log.Info("saved:",path)
	fout.Write(body)

//	message:=urlStr+"|||"+path
//TODO	publisher.Publish(kafka.NewMessage([]byte(message)))

	//	log.Info("enqueue parse,", path)
//	log.Debug("end saving url,", urlStr)

}
