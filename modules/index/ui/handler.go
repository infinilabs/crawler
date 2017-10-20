package ui

import (
	"github.com/julienschmidt/httprouter"

	"github.com/infinitbyte/gopa/core/http"
	"net/http"
	"github.com/infinitbyte/gopa/core/util"
	"strings"
	core "github.com/infinitbyte/gopa/core/index"
	"github.com/infinitbyte/gopa/modules/index"
)

// UserUI is the user namespace, public web
type UserUI struct {
	api.Handler
	Config *index.SearchUIConfig
	SearchClient *core.ElasticsearchClient
}

// IndexPageAction index page
func (h *UserUI) IndexPageAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	query:=h.GetParameter(req,"q")
	query=util.XSSHandle(query)
	if(strings.TrimSpace(query)==""){
		index.Index(w,h.Config)
	}else{
		h.SearchClient.Search()
		index.Search(w,req,query,h.Config)
	}
}
