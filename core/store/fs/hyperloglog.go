package fs

import (
	log "github.com/cihub/seelog"
	. "github.com/clarkduvall/hyperloglog"
	"io/ioutil"
	"github.com/medcl/gopa/core/util"
	"github.com/medcl/gopa/core/config"
"github.com/clarkduvall/hyperloglog"
	"hash"
	"hash/fnv"
)
type HyperLogLogFilter struct{
	persistFileName string
	filter *HyperLogLogPlus
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

func (filter *HyperLogLogFilter) Init(fileName string) error{

	filter.persistFileName=fileName

	//loading or initializing hyperloglog-filter
	if util.CheckFileExists(fileName) {
		log.Debug("found hyperloglog-filter,start reload,", fileName)
		n, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Error("hyperloglog-filter:",fileName, err)
		}

		filter.filter=&HyperLogLogPlus{}
		if err := filter.filter.GobDecode(n); err != nil {
			log.Error("hyperloglog-filter:",fileName, err)
		}

		log.Info("hyperloglog-filter successfully reloaded:",fileName)
	} else {
		probItems := config.GetIntConfig(config.HyperLogLogSection,config.HyperLogLogPrecision, 16)
		log.Debug("initializing hyperloglog-filter",fileName,",virual size is,", probItems)
		var er error
		filter.filter,er = hyperloglog.NewPlus(uint8(probItems))
		if(er!=nil){
			log.Info("hyperloglog-filter successfully initialized:",fileName)
		}else{
			log.Trace("hyperloglog-filter initialize failed:",fileName)
		}
	}

	return nil
}

func (filter *HyperLogLogFilter) Persist() error{

	log.Debug("hyperloglog-filter start persist,file:",filter.persistFileName)

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

func (filter *HyperLogLogFilter) Lookup(key []byte) bool{
	var count1=filter.filter.Count()
    filter.filter.Add(hash64(key))
	var count2=filter.filter.Count()
	if(count2 == count1){
		return true
	}
	return false
}

func (filter *HyperLogLogFilter) Add(key []byte) error{
	filter.filter.Add(hash64(key))
	return nil
}