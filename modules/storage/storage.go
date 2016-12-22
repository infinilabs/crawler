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

package storage

import (
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/store"
	"github.com/medcl/gopa/modules/storage/boltdb"
)

var impl boltdb.BoltdbStore

func (this StorageModule) Name() string {
	return "Storage"
}

func (this StorageModule) Start(env *Env) {

	impl = boltdb.BoltdbStore{FileName: env.RuntimeConfig.PathConfig.Data + "/boltdb"}
	err := impl.Open()
	if err != nil {
		panic(err)
	}
	store.Register(impl)
}

func (this StorageModule) Stop() error {
	err := impl.Close()
	return err

}

type StorageModule struct {
}
