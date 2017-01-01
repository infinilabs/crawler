package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/medcl/gopa/modules/api/http/templates/boltdb"
	"github.com/medcl/gopa/core/global"
	"github.com/asdine/storm"
	"github.com/medcl/gopa/modules/config"
)

func (h Handler) BoltDBStatusAction(w http.ResponseWriter, r *http.Request) {
	db:=global.Lookup(config.REGISTER_BOLTDB).(*storm.DB)

	//go func() {
	//	// Grab the initial stats.
	//	prev := db.Stats()
	//
	//	for {
	//		// Wait for 10s.
	//		time.Sleep(10 * time.Second)
	//
	//		// Grab the current stats and diff them.
	//		stats := db.Stats()
	//		diff := stats.Sub(&prev)
	//
	//		// Encode stats to JSON and print to STDERR.
	//		json.NewEncoder(os.Stderr).Encode(diff)
	//
	//		// Save stats for the next loop.
	//		prev = stats
	//	}
	//}()

	err := db.Bolt.View(func(tx *bolt.Tx) error {
		showUsage := (r.FormValue("usage") == "true")
		// Use the direct page id, if available.
		if r.FormValue("id") != "" {
			id, _ := strconv.Atoi(r.FormValue("id"))
			return templates.Page(w, r, tx, nil, id, showUsage)
		}

		// Otherwise extract the indexes and traverse.
		indexes, err := indexes(r)
		if err != nil {
			return err
		}

		return templates.Page(w, r, tx, indexes, 0, showUsage)
	})
	if err != nil {
		templates.Error(w, err)
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
