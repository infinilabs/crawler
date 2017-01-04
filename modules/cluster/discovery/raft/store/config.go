package store

import (
	"github.com/medcl/gopa/core/global"
	"path"
	"strings"
	log "github.com/cihub/seelog"
)

type RaftConfig struct {
	Bind    string `json:bind`
	Seeds   []string `json:seeds`
	DataDir string `json:data`
}

func (this *RaftConfig)Init()  {
	this.DataDir=path.Join(global.Env().SystemConfig.Data,"raft")

	if(len(global.Env().SystemConfig.ClusterBinding)>0){
		this.Bind=global.Env().SystemConfig.ClusterBinding
	}else{
		this.Bind=":13001"
	}

	join:=global.Env().SystemConfig.ClusterSeeds

	log.Debug("get cluster seeds: ",global.Env().SystemConfig.ClusterSeeds)

	if(len(join)>0){
		arr:=strings.Split(join,",")
		for _,v:=range arr{
			this.Seeds =append(this.Seeds,v)
		}
	}
}