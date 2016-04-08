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
	. "github.com/medcl/gopa/core/config"
	"net/http"
)

type Handler struct {
	Config      *GopaConfig
	wroteHeader bool
	http.ResponseWriter
}

func (w *Handler) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.wroteHeader = true
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

func (this *Handler) WriteJson(w http.ResponseWriter, v interface{}) error {
	if !this.wroteHeader {
		this.WriteJsonHeader(w)
		w.WriteHeader(http.StatusOK)
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

func (w *Handler) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func (w *Handler) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	flusher := w.ResponseWriter.(http.Flusher)
	flusher.Flush()
}
