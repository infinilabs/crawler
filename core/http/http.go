package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

// RegisteredAPIHandler is a hub for registered api
var RegisteredAPIHandler map[string]http.Handler

// RegisteredAPIFuncHandler is a hub for registered api
var RegisteredAPIFuncHandler map[string]func(http.ResponseWriter, *http.Request)

// RegisteredAPIMethodHandler is a hub for registered api
var RegisteredAPIMethodHandler map[string]map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params)

var l sync.Mutex

// HandleAPIFunc register api handler to specify pattern
func HandleAPIFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	l.Lock()
	if RegisteredAPIFuncHandler == nil {
		RegisteredAPIFuncHandler = map[string]func(http.ResponseWriter, *http.Request){}
	}
	RegisteredAPIFuncHandler[pattern] = handler
	l.Unlock()
}

// HandleAPI register api handler
func HandleAPI(pattern string, handler http.Handler) {

	l.Lock()
	if RegisteredAPIHandler == nil {
		RegisteredAPIHandler = map[string]http.Handler{}
	}
	RegisteredAPIHandler[pattern] = handler
	l.Unlock()
}

// HandleAPIMethod register api handler
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

// RegisteredUIHandler is a hub for registered ui handler
var RegisteredUIHandler map[string]http.Handler

// RegisteredUIFuncHandler is a hub for registered ui handler
var RegisteredUIFuncHandler map[string]func(http.ResponseWriter, *http.Request)

// RegisteredUIMethodHandler is a hub for registered ui handler
var RegisteredUIMethodHandler map[string]map[string]func(w http.ResponseWriter, req *http.Request, ps httprouter.Params)

// HandleUIFunc register ui request handler
func HandleUIFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	l.Lock()
	if RegisteredUIFuncHandler == nil {
		RegisteredUIFuncHandler = map[string]func(http.ResponseWriter, *http.Request){}
	}
	RegisteredUIFuncHandler[pattern] = handler
	l.Unlock()
}

// HandleUI register ui request handler
func HandleUI(pattern string, handler http.Handler) {

	l.Lock()
	if RegisteredUIHandler == nil {
		RegisteredUIHandler = map[string]http.Handler{}
	}
	RegisteredUIHandler[pattern] = handler
	l.Unlock()
}

// HandleUIMethod register ui request handler
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
