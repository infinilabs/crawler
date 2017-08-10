package config

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/global"
	"path"
)

type Command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type RaftConfig struct {
	Bind             string   `json:"bind"`
	Seeds            []string `json:"seeds"`
	DataDir          string   `json:"data"`
	EnableSingleNode bool     `json:"enable_single_node"`
}

type Node struct {
	IP   string
	Port int
}

type ClusterState struct {
	ActiveNodes   []Node
	InactiveNodes []Node
}

func (module *RaftConfig) Init() {
	module.DataDir = path.Join(global.Env().SystemConfig.GetWorkingDir(), "raft")
	module.EnableSingleNode = true

	if len(global.Env().SystemConfig.ClusterBinding) > 0 {
		module.Bind = global.Env().SystemConfig.ClusterBinding
	} else {
		module.Bind = "127.0.0.1:13001"
	}

	seeds := global.Env().SystemConfig.ClusterConfig.Seeds

	if len(seeds) > 0 {
		log.Debug("get cluster seeds: ", global.Env().SystemConfig.ClusterConfig.Seeds)
		for _, v := range seeds {
			module.Seeds = append(module.Seeds, v)
		}
	} else {
		module.Seeds = append(module.Seeds, "127.0.0.1:13001")
	}
}
