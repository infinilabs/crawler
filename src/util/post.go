/**
 * User: Medcl
 * Date: 13-7-22
 * Time: 下午12:23
 */
package util

import (
	log "logging"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Post(url string, values url.Values) []byte {
	r, err := http.PostForm(url, values)
	if err != nil {
		log.Error("post:", err)
		return nil
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("post:", err)
		return nil
	}
	return b
}
