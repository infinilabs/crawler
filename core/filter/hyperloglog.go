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

package filter

import (
	"hash"
	"hash/fnv"
	"io/ioutil"

	log "github.com/cihub/seelog"
	"github.com/clarkduvall/hyperloglog"
	. "github.com/clarkduvall/hyperloglog"
	"github.com/medcl/gopa/core/util"
	"sync"
)

type HyperLogLogFilter struct {
	hyperLogLogPrecision uint8
	persistFileName string
	filter          *HyperLogLogPlus
	l sync.Mutex
}

func hash32(s []byte) hash.Hash32 {
	h := fnv.New32a()
	h.Write(s)
	return h
}

func hash64(s []byte) hash.Hash64 {
	h := fnv.New64a()
	h.Write(s)
	return h
}

func (filter *HyperLogLogFilter) Init(fileName string) error {
	filter.l.Lock()
	defer filter.l.Unlock()

	filter.persistFileName = fileName
	filter.hyperLogLogPrecision = 18

	//loading or initializing hyperloglog-filter
	if util.FileExists(fileName) {
		log.Debug("found hyperloglog-filter,start reload,", fileName)
		n, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Error("hyperloglog-filter:", fileName, err)
		}

		filter.filter = &HyperLogLogPlus{}
		if err := filter.filter.GobDecode(n); err != nil {
			log.Error("hyperloglog-filter:", fileName, err)
		}

		log.Info("hyperloglog-filter successfully reloaded:", fileName)
	} else {
		log.Debug("initializing hyperloglog-filter", fileName, ",precision,", filter.hyperLogLogPrecision)
		var er error
		filter.filter, er = hyperloglog.NewPlus(filter.hyperLogLogPrecision)
		if er != nil {
			log.Info("hyperloglog-filter successfully initialized:", fileName)
		} else {
			log.Trace("hyperloglog-filter initialize failed:", fileName)
		}
	}

	return nil
}

func (filter *HyperLogLogFilter) Persist() error {
	filter.l.Lock()
	defer filter.l.Unlock()

	log.Debug("hyperloglog-filter start persist,file:", filter.persistFileName)

	//save
	m, err := filter.filter.GobEncode()
	if err != nil {
		log.Error(err)
		return nil
	}
	err = ioutil.WriteFile(filter.persistFileName, m, 0600)
	if err != nil {
		panic(err)
		return nil
	}
	log.Info("hyperloglog-filter safety persisted.")

	return nil
}

func (filter *HyperLogLogFilter) Exists(key []byte) bool {
	filter.l.Lock()
	defer filter.l.Unlock()

	var count1 = filter.filter.Count()
	filter.filter.Add(hash64(key))
	var count2 = filter.filter.Count()
	if count2 == count1 {
		return false
	}
	if(count2==count1+1){
		return true
	}
	log.Errorf("error hyperloglog behavior, %d vs %d",count1,count2)
	return true
}

func (filter *HyperLogLogFilter) Add(key []byte) error {
	filter.l.Lock()
	defer filter.l.Unlock()

	filter.filter.Add(hash64(key))
	return nil
}
