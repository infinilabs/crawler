package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

var RegisteredAPIHandler map[string]http.Handler
var RegisteredAPIFuncHandler map[string]func(http.ResponseWriter, *http.Request)
var RegisteredAPIMethodHandler map[string]map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params)

var l sync.Mutex

func HandleAPIFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	l.Lock()
	if RegisteredAPIFuncHandler == nil {
		RegisteredAPIFuncHandler = map[string]func(http.ResponseWriter, *http.Request){}
	}
	RegisteredAPIFuncHandler[pattern] = handler
	l.Unlock()
}

func HandleAPI(pattern string, handler http.Handler) {

	l.Lock()
	if RegisteredAPIHandler == nil {
		RegisteredAPIHandler = map[string]http.Handler{}
	}
	RegisteredAPIHandler[pattern] = handler
	l.Unlock()
}

func HandleAPIMethod(method Method, pattern string, handler func(w http.ResponseWriter, req *http.Request, ps httprouter.Params)) {
	l.Lock()
	if RegisteredAPIMethodHandler == nil {
		RegisteredAPIMethodHandler = map[string]map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params){}
	}

	m := RegisteredAPIMethodHandler[string(method)]
	if m == nil {
		RegisteredAPIMethodHandler[string(method)] = map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params){}
	}
	RegisteredAPIMethodHandler[string(method)][pattern] = handler
	l.Unlock()
}

var RegisteredUIHandler map[string]http.Handler
var RegisteredUIFuncHandler map[string]func(http.ResponseWriter, *http.Request)
var RegisteredUIMethodHandler map[string]map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params)

func HandleUIFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	l.Lock()
	if RegisteredUIFuncHandler == nil {
		RegisteredUIFuncHandler = map[string]func(http.ResponseWriter, *http.Request){}
	}
	RegisteredUIFuncHandler[pattern] = handler
	l.Unlock()
}

func HandleUI(pattern string, handler http.Handler) {

	l.Lock()
	if RegisteredUIHandler == nil {
		RegisteredUIHandler = map[string]http.Handler{}
	}
	RegisteredUIHandler[pattern] = handler
	l.Unlock()
}

func HandleUIMethod(method Method, pattern string, handler func(w http.ResponseWriter, req *http.Request, ps httprouter.Params)) {
	l.Lock()
	if RegisteredUIMethodHandler == nil {
		RegisteredUIMethodHandler = map[string]map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params){}
	}

	m := RegisteredUIMethodHandler[string(method)]
	if m == nil {
		RegisteredUIMethodHandler[string(method)] = map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params){}
	}
	RegisteredUIMethodHandler[string(method)][pattern] = handler
	l.Unlock()
}
