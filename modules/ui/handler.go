package ui

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/modules/config"
	"github.com/medcl/gopa/modules/ui/templates/boltdb"
	"github.com/medcl/gopa/modules/ui/templates/dashboard"
	"github.com/medcl/gopa/modules/ui/templates/tasks"
	"github.com/medcl/gopa/modules/ui/templates/console"
	"github.com/medcl/gopa/modules/ui/templates/explore"
)



func (h API) BoltDBStatusAction(w http.ResponseWriter, r *http.Request) {
	db := global.Lookup(config.REGISTER_BOLTDB).(*storm.DB)
	err := db.Bolt.View(func(tx *bolt.Tx) error {
		showUsage := (r.FormValue("usage") == "true")
		// Use the direct page id, if available.
		if r.FormValue("id") != "" {
			id, _ := strconv.Atoi(r.FormValue("id"))
			return boltdb.Page(w, r, tx, nil, id, showUsage)
		}

		// Otherwise extract the indexes and traverse.
		indexes, err := indexes(r)
		if err != nil {
			return err
		}

		return boltdb.Page(w, r, tx, indexes, 0, showUsage)
	})
	if err != nil {
		boltdb.Error(w, err)
	}
}

// parses and returns all indexes from a request.
func indexes(r *http.Request) ([]int, error) {
	var a = []int{0}
	if len(r.FormValue("index")) > 0 {
		for _, s := range strings.Split(r.FormValue("index"), ":") {
			i, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			a = append(a, i)
		}
	}
	return a, nil
}


func (h API) DashboardAction(w http.ResponseWriter, r *http.Request) {

	 dashboard.Index(w)
}

func (h API) TasksPageAction(w http.ResponseWriter, r *http.Request) {

	 tasks.Index(w)
}
func (h API) ConsolePageAction(w http.ResponseWriter, r *http.Request) {

	 console.Index(w)
}
func (h API) ExplorePageAction(w http.ResponseWriter, r *http.Request) {

	 explore.Index(w)
}
