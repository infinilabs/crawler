package user

import (
	"github.com/infinitbyte/gopa/core/http"
	"github.com/infinitbyte/gopa/modules/ui/user/search"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// UserUI is the user namespace, public web
type UserUI struct {
	api.Handler
}

// IndexPageAction index page
func (h UserUI) IndexPageAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	search.Search(w)
}
