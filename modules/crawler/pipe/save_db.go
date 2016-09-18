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

package pipe

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/types"
	"github.com/medcl/gopa/modules/storage/boltdb"
	"strings"
	"path"
)

type SaveToDBJoint struct {
	context *Context
}

func (this SaveToDBJoint) Process(c *Context) (*Context, error) {
	this.context = c

	url := c.MustGetString(CONTEXT_URL)

	pageItem := c.Get(CONTEXT_PAGE_ITEM).(*types.PageItem)
	savePath := c.MustGetString(CONTEXT_SAVE_PATH)
	saveFile := c.MustGetString(CONTEXT_SAVE_FILENAME)

	log.Debug("save url to db,", url, ",domain,", pageItem.Domain)

	this.context.Env.RuntimeConfig.Storage.AddValue(boltdb.SnapshotBucketKey,GetKey(pageItem.Domain,path.Join(savePath,saveFile)),pageItem.Body)

	return c, nil
}

const KeyDelimiter string = "||"
func GetKey( args ...string) []byte {
	return []byte(strings.Join(args,KeyDelimiter))
}
