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
	. "github.com/medcl/gopa/core/pipeline"
	"encoding/hex"
	"crypto/md5"
)


type PublishJoint struct {
}

func (this PublishJoint) Process(c *Context) (*Context, error) {

	m := md5.Sum([]byte(c.MustGetString(CONTEXT_URL.String())))
	id:=hex.EncodeToString(m[:])

	data:=map[string]interface{}{}
	meta,b:= c.GetMap(CONTEXT_PAGE_METADATA.String())
	if(b){
		data["metadata"]=meta
	}
	_,err:= c.Env.ESClient.IndexDoc(id,data)
	if(err!=nil){
		return c,err
	}

	return c, nil
}
