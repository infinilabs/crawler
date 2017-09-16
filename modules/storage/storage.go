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
	"github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/modules/storage/boltdb"
	"os"
	"path"
)

var impl boltdb.BoltdbStore

func (this StorageModule) Name() string {
	return "Storage"
}

func (module StorageModule) Start(cfg *config.Config) {

	folder := path.Join(global.Env().SystemConfig.GetWorkingDir(), "blob")
	os.MkdirAll(folder, 0777)
	impl = boltdb.BoltdbStore{FileName: path.Join(folder, "/bolt.db")}
	err := impl.Open()
	if err != nil {
		panic(err)
	}
	persist.RegisterKVHandler(impl)
}

func (module StorageModule) Stop() error {
	err := impl.Close()
	return err

}

type StorageModule struct {
}
