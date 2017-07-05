package user

import (
	"github.com/infinitbyte/gopa/core/http"
	"github.com/infinitbyte/gopa/modules/ui/user/search"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type UserUI struct {
	api.Handler
}

func (h UserUI) IndexPageAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	search.Search(w)
}
