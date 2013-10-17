/**
 * Created with IntelliJ IDEA.
 * User: medcl
 * Date: 13-10-17
 * Time: 下午5:21
 */
package config

import (
	"store"
	. "github.com/zeebo/sbloom"
	. "types"
)

type RuntimeConfig struct{
	Storage store.Store
	BloomFilter *Filter
	TaskConfig *TaskConfig
}

