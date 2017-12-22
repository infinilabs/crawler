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

// Method is object of http method
type Method string

const (
	// GET is http get method
	GET Method = "GET"
	// POST is http post method
	POST Method = "POST"
	// PUT is http put method
	PUT Method = "PUT"
	// DELETE is http delete method
	DELETE Method = "DELETE"
	// HEAD is http head method
	HEAD Method = "HEAD"
)

// String return http method as string
func (method Method) String() string {
	switch method {
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

// Handler is the object of http handler
type Handler struct {
	wroteHeader bool

	//w http.ResponseWriter
	//req *http.Request
	//
	formParsed bool
}

// WriteHeader write status code to http header
func (handler Handler) WriteHeader(w http.ResponseWriter, code int) {

	if len(global.Env().SystemConfig.PathConfig.Cert) > 0 {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	}

	w.WriteHeader(code)
	handler.wroteHeader = true
}

// Get http parameter or return default value
func (handler Handler) Get(req *http.Request, key string, defaultValue string) string {
	if !handler.formParsed {
		req.ParseForm()
	}
	if len(req.Form) > 0 {
		return req.Form.Get(key)
	}
	return defaultValue
}

// EncodeJSON encode the object to json string
func (handler Handler) EncodeJSON(v interface{}) (b []byte, err error) {

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

// WriteJSONHeader will write standard json header
func (handler Handler) WriteJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	handler.wroteHeader = true
}

// Result is a general json result
type Result struct {
	Total  int         `json:"total"`
	Result interface{} `json:"result"`
}

// WriteJSONListResult output result list to json format
func (handler Handler) WriteJSONListResult(w http.ResponseWriter, total int, v interface{}, statusCode int) error {
	result := Result{}
	result.Total = total
	result.Result = v
	return handler.WriteJSON(w, result, statusCode)
}

// WriteJSON output signal result with json format
func (handler Handler) WriteJSON(w http.ResponseWriter, v interface{}, statusCode int) error {
	if !handler.wroteHeader {
		handler.WriteJSONHeader(w)
		w.WriteHeader(statusCode)
	}

	b, err := handler.EncodeJSON(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// GetParameter return query parameter with argument name
func (handler Handler) GetParameter(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetParameterOrDefault return query parameter or return default value
func (handler Handler) GetParameterOrDefault(r *http.Request, key string, defaultValue string) string {
	v := r.URL.Query().Get(key)
	if len(v) > 0 {
		return v
	}
	return defaultValue
}

// GetIntOrDefault return parameter or default, data type is int
func (handler Handler) GetIntOrDefault(r *http.Request, key string, defaultValue int) int {

	v := handler.GetParameter(r, key)
	s, ok := util.ToInt(v)
	if ok != nil {
		return defaultValue
	}
	return s

}

// GetJSON return json input
func (handler Handler) GetJSON(r *http.Request) (*jsonq.JsonQuery, error) {

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

// GetRawBody return raw http request body
func (handler Handler) GetRawBody(r *http.Request) ([]byte, error) {

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

// Write response to client
func (handler Handler) Write(w http.ResponseWriter, b []byte) (int, error) {
	return w.Write(b)
}

// Error404 output 404 response
func (handler Handler) Error404(w http.ResponseWriter) {
	handler.WriteJSON(w, map[string]interface{}{"error": 404}, http.StatusNotFound)
}

// Error500 output 500 response
func (handler Handler) Error500(w http.ResponseWriter, msg string) {
	handler.WriteJSON(w, map[string]interface{}{"error": msg}, http.StatusInternalServerError)
}

// Error output custom error
func (handler Handler) Error(w http.ResponseWriter, err error) {
	handler.WriteJSON(w, map[string]interface{}{"error": err.Error()}, http.StatusInternalServerError)
}

// Flush flush response message
func (handler Handler) Flush(w http.ResponseWriter) {
	flusher := w.(http.Flusher)
	flusher.Flush()
}

// BasicAuth register api with basic auth
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
