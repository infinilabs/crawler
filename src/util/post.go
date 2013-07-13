package util

import (
	"io/ioutil"
	log "github.com/cihub/seelog"
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
