package stats

import (
	stats "github.com/dmuth/golang-stats"
	"encoding/json"
)

func Increment(key string )  {
	stats.IncrStat(key)
}

func Decrement(key string)  {
	stats.DecrStat(key)
}

func Stats(key string)  {
	stats.Stat(key)
}

func StatsAll()string  {
	obj,_ :=json.Marshal(stats.StatAll())
	return string(obj)
}