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

package api

import (
	"encoding/json"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/jmoiron/jsonq"
	"github.com/julienschmidt/httprouter"
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

func (this Method) String() string {
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
	wroteHeader bool

	//w http.ResponseWriter
	//req *http.Request
	//
	formParsed bool
}

func (this Handler) WriteHeader(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	if len(global.Env().SystemConfig.PathConfig.Cert) > 0 {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	}
	this.wroteHeader = true
}

func (this Handler) Get(req *http.Request, key string, defaultValue string) string {
	if !this.formParsed {
		req.ParseForm()
	}
	if len(req.Form) > 0 {
		return req.Form.Get(key)
	}
	return defaultValue
}

func (w Handler) EncodeJson(v interface{}) (b []byte, err error) {

	//if(w.Get("pretty","false")=="true"){
	b, err = json.MarshalIndent(v, "", "  ")
	//}else{
	//	b, err = json.Marshal(v)
	//}

	if err != nil {
		return nil, err
	}
	return b, nil
}

func (this Handler) WriteJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	this.wroteHeader = true
}

type Result struct {
	Total  int         `json:"total"`
	Result interface{} `json:"result"`
}

func (this Handler) WriteListResultJson(w http.ResponseWriter, total int, v interface{}, statusCode int) error {
	result := Result{}
	result.Total = total
	result.Result = v
	return this.WriteJson(w, result, statusCode)
}

func (this Handler) WriteJson(w http.ResponseWriter, v interface{}, statusCode int) error {
	if !this.wroteHeader {
		this.WriteJsonHeader(w)
		w.WriteHeader(statusCode)
	}

	b, err := this.EncodeJson(v)
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

func (this Handler) GetParameter(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

func (this Handler) GetParameterOrDefault(r *http.Request, key string, defaultValue string) string {
	v := r.URL.Query().Get(key)
	if len(v) > 0 {
		return v
	}
	return defaultValue
}

func (this Handler) GetIntOrDefault(r *http.Request, key string, defaultValue int) int {

	v := this.GetParameter(r, key)
	s, ok := util.ToInt(v)
	if ok != nil {
		return defaultValue
	}
	return s

}

func (this Handler) GetJson(r *http.Request) (*jsonq.JsonQuery, error) {

	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.NewWithCode(err, errors.JSONIsEmpty, r.URL.String())
	}

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(content)))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	return jq, nil
}

func (this Handler) GetRawBody(r *http.Request) ([]byte, error) {

	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.NewWithCode(err, errors.BodyEmpty, r.URL.String())
	}
	return content, nil
}

func (this Handler) Write(w http.ResponseWriter, b []byte) (int, error) {
	if !this.wroteHeader {
		this.WriteHeader(w, http.StatusOK)
	}
	return w.Write(b)
}

func (this Handler) Error404(w http.ResponseWriter) {
	this.WriteJson(w, map[string]interface{}{"error": 404}, http.StatusNotFound)
}

func (this Handler) Error500(w http.ResponseWriter, msg string) {
	this.WriteJson(w, map[string]interface{}{"error": msg}, http.StatusInternalServerError)
}

func (this Handler) Error(w http.ResponseWriter, err error) {
	this.WriteJson(w, map[string]interface{}{"error": err.Error()}, http.StatusInternalServerError)
}

func (this Handler) Flush(w http.ResponseWriter) {
	if !this.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	flusher := w.(http.Flusher)
	flusher.Flush()
}

func BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}
