package fs

import (
	log "github.com/cihub/seelog"
	. "github.com/zeebo/sbloom"
	"io/ioutil"
	"github.com/medcl/gopa/src/util"
	"github.com/medcl/gopa/src/config"
	"hash/fnv"
)

type BloomFilter struct{
	persistFileName string
	filter *Filter
}


func (filter *BloomFilter) Init(fileName string) error{

	filter.persistFileName=fileName

	//loading or initializing bloom filter
	if util.CheckFileExists(fileName) {
		log.Debug("found bloomFilter,start reload,", fileName)
		n, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Error("bloomFilter:",fileName, err)
		}
		filter.filter=new (Filter)
		if err := filter.filter.GobDecode(n); err != nil {
			log.Error("bloomFilter:",fileName, err)
		}
		log.Info("bloomFilter successfully reloaded:",fileName)
	} else {
		probItems := config.GetIntConfig("BloomFilter", "ItemSize", 100000)
		log.Debug("initializing bloom-filter",fileName,",virual size is,", probItems)
		filter.filter = NewFilter(fnv.New64(), probItems)
		log.Info("bloomFilter successfully initialized:",fileName)
	}

	return nil
}

func (filter *BloomFilter) Persist() error{

	log.Debug("bloomFilter start persist,file:",filter.persistFileName)

	//save bloom-filter
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
	log.Info("bloomFilter safety persisted.")

	return nil
}

func (filter *BloomFilter) Lookup(key []byte) bool{
	return filter.filter.Lookup(key)
}

func (filter *BloomFilter) Add(key []byte) error{
	filter.filter.Add(key)
	return nil
}