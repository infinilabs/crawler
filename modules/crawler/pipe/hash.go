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
	"crypto/sha1"
	"fmt"
	log "github.com/cihub/seelog"
	. "github.com/gensmusic/simhash"
	"github.com/infinitbyte/gopa/core/model"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"path"
	"sync"
)

type HashJoint struct {
	Parameters
}

const simHashEnabled ParaKey = "simhash_enabled"
const simHashDictFolder ParaKey = "simhash_dict_folder"

func (joint HashJoint) Name() string {
	return "hash"
}

func (joint HashJoint) Process(context *Context) error {

	snapshot := context.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	body := string(snapshot.Payload)

	h := sha1.New()
	h.Write([]byte(body))
	bs := h.Sum(nil)

	snapshot.Hash = fmt.Sprintf("%x", bs)

	if joint.GetBool(simHashEnabled, false) {
		joint.loadDict()
		hash1 := Simhash(&body, 200)
		snapshot.SimHash = fmt.Sprintf("%x", hash1)
	}

	return nil
}

var loaded = false
var lock sync.Mutex

func (joint HashJoint) loadDict() {
	lock.Lock()
	defer lock.Unlock()
	if loaded {
		return
	}

	log.Debug("loading jieba dict files")
	mainDict := "config/dict/main.dict.txt"
	idfDict := "config/dict/idf.txt"
	stopwordsDict := "config/dict/stop_words.txt"
	if joint.Has(simHashDictFolder) {
		dictRoot := joint.MustGetString(simHashDictFolder)
		if len(dictRoot) > 0 {
			mainDict = path.Join(dictRoot, mainDict)
			idfDict = path.Join(dictRoot, idfDict)
			stopwordsDict = path.Join(dictRoot, stopwordsDict)
		}
	}
	if err := LoadDictionary(mainDict, idfDict, stopwordsDict); err != nil {
		log.Error("Failed to load dictionary:", err)
		panic(err)
	}
	loaded = true
}
