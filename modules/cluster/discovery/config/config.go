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

func (this *RaftConfig) Init() {
	this.DataDir = path.Join(global.Env().SystemConfig.GetDataDir(), "raft")
	this.EnableSingleNode = true

	if len(global.Env().SystemConfig.ClusterBinding) > 0 {
		this.Bind = global.Env().SystemConfig.ClusterBinding
	} else {
		this.Bind = "127.0.0.1:13001"
	}

	seeds := global.Env().SystemConfig.ClusterConfig.Seeds

	if len(seeds) > 0 {
		log.Debug("get cluster seeds: ", global.Env().SystemConfig.ClusterConfig.Seeds)
		for _, v := range seeds {
			this.Seeds = append(this.Seeds, v)
		}
	} else {
		this.Seeds = append(this.Seeds, "127.0.0.1:13001")
	}
}
