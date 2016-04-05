/**
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午5:42
 */
package tasks

import (
	"bytes"
	"compress/gzip"
	log "github.com/cihub/seelog"
	"io"
	"io/ioutil"
	net "net"
	"net/http"
	. "net/url"
	"strings"
	"time"
)

//parse to get url root
func getRootUrl(source *URL) string {
	if strings.HasSuffix(source.Path, "/") {
		return source.Host + source.Path
	} else {
		index := strings.LastIndex(source.Path, "/")
		if index > 0 {
			path := source.Path[0:index]
			return source.Host + path
		} else {
			return source.Host + "/"
		}
	}
	return ""
}

//format url,prepare for bloom filter
func formatUrlForFilter(url []byte) []byte {
	src := string(url)
	log.Trace("start to normalize url:", src)
	if strings.HasSuffix(src, "/") {
		src = strings.TrimRight(src, "/")
	}
	src = strings.TrimSpace(src)
	src = strings.ToLower(src)
	return []byte(src)
}

func getUrlPathFolderWithoutFile(url []byte) []byte {
	src := string(url)
	log.Trace("start to get url's path folder:", src)
	if strings.HasSuffix(src, "/") {
		src = strings.TrimRight(src, "/")
	}
	src = strings.TrimSpace(src)
	src = strings.ToLower(src)
	return []byte(src)
}

func get(url string, cookie string) ([]byte, error) {

	log.Debug("let's get :" + url)

	client := &http.Client{
		CheckRedirect: nil,
	}
	reqest, _ := http.NewRequest("GET", url, nil)

	reqest.Header.Set("User-Agent", " Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	reqest.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Set("Cache-Control", "max-age=0")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("Referer", url)

	if len(cookie) > 0 {
		log.Debug("dealing with cookie:" + cookie)
		array := strings.Split(cookie, ";")
		for item := range array {
			array2 := strings.Split(array[item], "=")
			if len(array2) == 2 {
				cookieObj := http.Cookie{}
				cookieObj.Name = array2[0]
				cookieObj.Value = array2[1]
				reqest.AddCookie(&cookieObj)
			} else {
				log.Info("error,index out of range:" + array[item])
			}
		}
	}

	resp, err := client.Do(reqest)

	if err != nil {
		log.Error(url, err)
		return nil, err
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Error(url, err)
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	if reader != nil {
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Error(url, err)
			return nil, err
		}
		return body, nil

	}
	return nil, nil
}

func post(url string, cookie string, postStr string) []byte {

	log.Debug("let's post :" + url)

	client := &http.Client{
		CheckRedirect: nil,
	}

	postBytesReader := bytes.NewReader([]byte(postStr))
	reqest, _ := http.NewRequest("POST", url, postBytesReader)

	reqest.Header.Set("User-Agent", " Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	reqest.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	//	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Set("Cache-Control", "max-age=0")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("Referer", url)

	if len(cookie) > 0 {
		log.Debug("dealing with cookie:" + cookie)
		array := strings.Split(cookie, ";")
		for item := range array {
			array2 := strings.Split(array[item], "=")
			if len(array2) == 2 {
				cookieObj := http.Cookie{}
				cookieObj.Name = array2[0]
				cookieObj.Value = array2[1]
				reqest.AddCookie(&cookieObj)
			} else {
				log.Info("error,index out of range:" + array[item])
			}
		}
	}

	resp, err := client.Do(reqest)

	if err != nil {
		log.Error(url, err)
		return nil
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Error(url, err)
			return nil
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	if reader != nil {
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Error(url, err)
			return nil
		}
		return body
	}
	return nil
}

func HttpGetWithCookie(resource string, cookie string) (msg []byte, err error) {

	out, err := get(resource, cookie)
	return out, err
}

func HttpGet(resource string) (msg []byte, err error) {

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err := net.DialTimeout(netw, addr, 5*time.Second) //连接超时时间
				if err != nil {
					log.Error(resource, err)
					return nil, err
				}

				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	req, err := http.NewRequest("GET", resource, nil)

	if err != nil {
		log.Error(resource, err)
		return nil, err
	}

	//support gzip
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; gopa/0.1; +http://infinitbyte.com/gopa)")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		log.Error(resource, err)
		return nil, err
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Error(resource, err)
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}
	if reader != nil {
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Error(resource, err)
			return nil, err
		}
		return body, nil
	}
	return nil, http.ErrNotSupported
}
