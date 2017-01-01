/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handler

import (
	log "github.com/cihub/seelog"
	_ "github.com/jmoiron/jsonq"
	"net/http"
	"github.com/medcl/gopa/core/store"
	"github.com/medcl/gopa/modules/config"
)

func (this Handler) SnapshotAction(w http.ResponseWriter, req *http.Request) {

	if req.Method == GET.String() {
		url := this.GetParameter(req, "url")
		log.Trace("get snapsthot by url:", string(url))

		compressed := this.GetParameterOrDefault(req, "compressed", "true")
		var bytes []byte
		if compressed == "true" {
			bytes = store.GetCompressedValue(config.SnapshotBucketKey, []byte(url))
		} else {
			bytes = store.GetValue(config.SnapshotBucketKey, []byte(url))
		}

		if len(bytes) > 0 {
			this.Write(w, bytes)
			this.Write(w,[]byte("<script src=\"/ui/assets/js/snapshot_footprint.js?v=1\"></script> "))
			return
		}

	}

	this.error404(w)

}
