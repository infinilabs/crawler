package http

import (
	"net/http"
	"sync"
)

var RegisteredHandler map[string]http.Handler
var RegisteredFuncHandler map[string]func(http.ResponseWriter, *http.Request)

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
