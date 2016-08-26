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

package handler

import (
	"encoding/json"
	logger "github.com/cihub/seelog"
	"github.com/jmoiron/jsonq"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/types"
	"io/ioutil"
	"net/http"
	"strings"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	HEAD   Method = "HEAD"
)

func (this Method) String() string  {
	switch this {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case HEAD:
		return "HEAD"
	}
	return "N/A"
}


type Handler struct {
	Env         *Env
	wroteHeader bool
}

func (this *Handler) WriteHeader(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	this.wroteHeader = true
}

func (w *Handler) encodeJson(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (this *Handler) WriteJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (this *Handler) WriteJson(w http.ResponseWriter, v interface{}, statusCode int) error {
	if !this.wroteHeader {
		this.WriteJsonHeader(w)
		w.WriteHeader(statusCode)
	}

	b, err := this.encodeJson(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

type ErrEmptyJson struct {
}

func (this *Handler) GetJson(r *http.Request) (*jsonq.JsonQuery, error) {

	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, types.JSONIsEmpty
	}
	logger.Trace("receive json:", string(content))

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(content)))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	return jq, nil
}

func (this *Handler) GetRawBody(r *http.Request) ([]byte, error) {

	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, types.BodyEmpty
	}
	return content, nil
}

func (this *Handler) Write(w http.ResponseWriter, b []byte) (int, error) {
	if !this.wroteHeader {
		this.WriteHeader(w, http.StatusOK)
	}
	return w.Write(b)
}

func (this *Handler) error404(w http.ResponseWriter) {
	this.WriteJson(w, map[string]interface{}{"error": 404}, http.StatusNotFound)
}

func (this *Handler) error500(w http.ResponseWriter, msg string) {
	this.WriteJson(w, map[string]interface{}{"error": msg}, http.StatusInternalServerError)
}

func (this *Handler) Flush(w http.ResponseWriter) {
	if !this.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	flusher := w.(http.Flusher)
	flusher.Flush()
}
