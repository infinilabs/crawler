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

package api

import (
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/framework/core/api"
	"github.com/infinitbyte/framework/core/api/router"
	"github.com/infinitbyte/framework/core/kv"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
	_ "github.com/jmoiron/jsonq"
	"net/http"
	"strconv"
)

func (this API) SnapshotAction(w http.ResponseWriter, req *http.Request) {

	if req.Method == GET.String() {
		url := this.GetParameter(req, "id")
		log.Trace("get snapsthot by url:", string(url))

		compressed := this.GetParameterOrDefault(req, "compressed", "true")
		var bytes []byte
		var err error
		if compressed == "true" {
			bytes, err = kv.GetCompressedValue(config.SnapshotBucketKey, []byte(url))
		} else {
			bytes, err = kv.GetValue(config.SnapshotBucketKey, []byte(url))
		}

		if err != nil {
			this.Error(w, err)
			return
		}

		if len(bytes) > 0 {
			this.Write(w, bytes)
			this.Write(w, []byte("<script src=\"/ui/assets/js/snapshot_footprint.js?v=1\"></script> "))
			return
		}

	}

	this.Error404(w)

}

func (this API) SnapshotListAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	fr := this.GetParameter(req, "from")
	si := this.GetParameter(req, "size")
	taskId := this.GetParameter(req, "task_id")

	from, err := strconv.Atoi(fr)
	if err != nil {
		from = 0
	}
	size, err := strconv.Atoi(si)
	if err != nil {
		size = 10
	}

	total, snapshots, err := model.GetSnapshotList(from, size, taskId)
	if err != nil {
		this.Error(w, err)
	} else {
		this.WriteJSONListResult(w, total, snapshots, http.StatusOK)
	}

}

func (this API) SnapshotGetAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	snapshot, err := model.GetSnapshot(id)
	if err != nil {
		this.Error(w, err)
	} else {
		this.WriteJSON(w, snapshot, http.StatusOK)

	}

}

func (this API) SnapshotGetPayloadAction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var err error
	compressed := this.GetParameterOrDefault(req, "compressed", "true")
	var bytes []byte
	if compressed == "true" {
		bytes, err = kv.GetCompressedValue(config.SnapshotBucketKey, []byte(id))
	} else {
		bytes, err = kv.GetValue(config.SnapshotBucketKey, []byte(id))
	}

	if err != nil {
		this.Error(w, err)
		return
	}

	if len(bytes) > 0 {
		this.Write(w, bytes)
	}

	this.Error404(w)

}
