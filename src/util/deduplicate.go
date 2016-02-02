package util

import (
	"github.com/clarkduvall/hyperloglog"
	"github.com/medcl/gopa/src/config"
)

func checkWithValue(val string)  {
	var precision int
	precision=config.GetIntConfig(config.HyperLogLogSection,config.HyperLogLogPrecision,50000)

	hyperloglog.NewPlus(uint8(precision))
}

type DeduplicatePlugin interface{
	Init(fileName string) error
	Persist() error
	Lookup(key []byte) bool
	Add(key []byte) error
}