package api

import (
	"net/http"
	"sync"
	"github.com/julienschmidt/httprouter"
)

var RegisteredHandler map[string]http.Handler
var RegisteredFuncHandler map[string]func(http.ResponseWriter, *http.Request)
var RegisteredMethodHandler map[string]map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params)

var l sync.Mutex

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	l.Lock()
	if RegisteredFuncHandler == nil {
		RegisteredFuncHandler = map[string]func(http.ResponseWriter, *http.Request){}
	}
	RegisteredFuncHandler[pattern] = handler
	l.Unlock()
}

func Handle(pattern string, handler http.Handler) {

	l.Lock()
	if RegisteredHandler == nil {
		RegisteredHandler = map[string]http.Handler{}
	}
	RegisteredHandler[pattern] = handler
	l.Unlock()
}

func HandleMethod(method Method,pattern string, handler func (w http.ResponseWriter, req *http.Request, ps httprouter.Params)) {
	l.Lock()
	if RegisteredMethodHandler == nil {
		RegisteredMethodHandler = map[string]map[string]func (w http.ResponseWriter, req *http.Request, ps httprouter.Params){}
	}

	m:=RegisteredMethodHandler[string(method)]
	if(m==nil){
		RegisteredMethodHandler[string(method)]=map[string]func (w http.ResponseWriter, req *http.Request, ps httprouter.Params){}
	}
	RegisteredMethodHandler[string(method)][pattern]=handler
	l.Unlock()
}
