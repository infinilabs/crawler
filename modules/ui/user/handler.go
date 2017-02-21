package user

import (
	"github.com/julienschmidt/httprouter"
	"github.com/medcl/gopa/core/http"
	"github.com/medcl/gopa/modules/ui/user/search"
	"net/http"
)

type UserUI struct {
	api.Handler
}

func (h UserUI) IndexPageAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	search.Search(w)
}
