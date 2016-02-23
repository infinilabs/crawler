/**
 * Created with IntelliJ IDEA.
 * User: medcl
 * Date: 13-11-8
 * Time: 下午6:32
 * To change this template use File | Settings | File Templates.
 */
package http

import (
	"net/http"
	"github.com/pantsing/gograce/ghttp"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
)

func index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("gopa!"))
	w.Write([]byte("\nversion: "+config.Version))
	w.Write([]byte("\ncluster: "+config.ClusterConfig.Name))
}

var config RuntimeConfig
func Start(runtimeConfig *RuntimeConfig) {
	config=*runtimeConfig
	http.HandleFunc("/", index)
	log.Info("http server listen at: http://localhost:8001/")
	ghttp.ListenAndServe(":8001", nil)
}
