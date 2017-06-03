package config

import (
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/global"
	"path"
	"strings"
)

type Command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type RaftConfig struct {
	Bind    string   `json:"bind"`
	Seeds   []string `json:"seeds"`
	DataDir string   `json:"data"`
}

func (this *RaftConfig) Init() {
	this.DataDir = path.Join(global.Env().SystemConfig.GetDataDir(), "raft")

	if len(global.Env().SystemConfig.ClusterBinding) > 0 {
		this.Bind = global.Env().SystemConfig.ClusterBinding
	} else {
		this.Bind = ":13001"
	}

	join := global.Env().SystemConfig.ClusterConfig.Seeds

	log.Debug("get cluster seeds: ", global.Env().SystemConfig.ClusterConfig.Seeds)

	if len(join) > 0 {
		arr := strings.Split(join, ",")
		for _, v := range arr {
			this.Seeds = append(this.Seeds, v)
		}
	}
}
