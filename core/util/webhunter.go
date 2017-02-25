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

package util

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	. "net/url"
	"strings"
	"time"

	"fmt"
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/errors"
	. "github.com/medcl/gopa/core/model"
	"golang.org/x/net/proxy"
)

func GetHost(url string) string {
	uri, err := ParseRequestURI(url)
	if err != nil {
		log.Trace(err)
		return ""
	}
	return uri.Host
}

//parse to get url root
func GetRootUrl(source *URL) string {
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

//format url, normalize url
func FormatUrlForFilter(url []byte) []byte {
	src := string(url)
	log.Trace("start to normalize url:", src)
	if strings.HasSuffix(src, "/") {
		src = strings.TrimRight(src, "/")
	}
	src = strings.TrimSpace(src)
	src = strings.ToLower(src)
	return []byte(src)
}

func GetUrlPathFolderWithoutFile(url []byte) []byte {
	src := string(url)
	log.Trace("start to get url's path folder:", src)
	if strings.HasSuffix(src, "/") {
		src = strings.TrimRight(src, "/")
	}
	src = strings.TrimSpace(src)
	src = strings.ToLower(src)
	return []byte(src)
}

func noRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("catch http redirect!")
}

func getUrl(url string) (string, error) {
	if !strings.HasPrefix(url, "http") {
		return url, errors.New("invalid url, " + url)
		//url="http://"+url
	}
	return url, nil
}

/**
proxyStr, eg: "socks5://127.0.0.1:9150"
*/
func get(page *PageItem, url string, cookie string, proxyStr string) ([]byte, error) {

	log.Debug("let's get :" + url)

	var err error
	url, err = getUrl(url)
	if err != nil {
		return nil, errors.New("invalid url: " + url)
	}

	client := &http.Client{
		CheckRedirect: noRedirect,
		//CheckRedirect: nil,
	}

	if proxyStr != "" {

		// Create a transport that uses Tor Browser's SocksPort.  If
		// talking to a system tor, this may be an AF_UNIX socket, or
		// 127.0.0.1:9050 instead.
		tbProxyURL, err := Parse(proxyStr)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to parse proxy URL: %v", err))
		}

		// Get a proxy Dialer that will create the connection on our
		// behalf via the SOCKS5 proxy.  Specify the authentication
		// and re-create the dialer/transport/client if tor's
		// IsolateSOCKSAuth is needed.
		tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to obtain proxy dialer: %v", err))
		}

		// Make a http.Transport that uses the proxy dialer, and a
		// http.Client that uses the transport.
		tbTransport := &http.Transport{Dial: tbDialer.Dial}
		//http.DefaultClient.Transport = &http.Transport{Dial: tbDialer.Dial}

		client.Transport = tbTransport

	}

	reqest, _ := http.NewRequest("GET", url, nil)

	//req.SetBasicAuth("user", "password")

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

	log.Trace("response: ", err, ", ", resp)

	if resp != nil && (resp.StatusCode == 301 || resp.StatusCode == 302) {
		log.Debug("got redirect: ", url, " => ", resp.Header.Get("Location"))
		location := resp.Header.Get("Location")
		if len(location) > 0 && location != url {
			return nil, errors.NewWithPayload(errors.URLRedirected, fmt.Sprint("got redirect: ", url, " => ", location), location)
		}
	}

	if err != nil {
		log.Error(url, ", ", err)
		return nil, err
	}

	log.Trace("status code,", resp.StatusCode, ",size,", resp.ContentLength)

	log.Trace("host: ", resp.Request.Host, " url: ", resp.Request.URL.String())

	//update host, redirects may change the host
	page.Host = resp.Request.Host
	page.Url = resp.Request.URL.String()

	page.StatusCode = resp.StatusCode
	if resp.Header != nil {
		page.Headers = resp.Header
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Error(url, ", ", err)
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	if reader != nil {
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Error(url, ", ", err)
			return nil, err
		}
		page.Body = body
		page.Size = len(body)
		return body, nil

	}
	return nil, nil
}

func post(url string, cookie string, postStr string) []byte {

	log.Debug("let's post: " + url)

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
		log.Error(url, ", ", err)
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

func HttpGetWithCookie(page *PageItem, resource string, cookie string, proxy string) (msg []byte, err error) {
	out, err := get(page, resource, cookie, proxy)
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

func sendHTTPRequest(method, url string, body io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	newReq, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer newReq.Body.Close()
	response, err := ioutil.ReadAll(newReq.Body)
	if err != nil {
		return nil, err
	}

	if newReq.StatusCode > http.StatusCreated && newReq.StatusCode < http.StatusNotFound {
		return nil, errors.New(string(response))
	}

	return response, nil
}
