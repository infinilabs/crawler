package cluster

import (
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/modules/cluster/discovery/raft"
)

type ClusterModule struct {
}

func (this ClusterModule) Name() string {
	return "Cluster"
}

func (this ClusterModule) Start(env *Env) {

	s := raft.New()
	if err := s.Open(); err != nil {
		log.Errorf("failed to open raft: %s", err.Error())
		panic(err)
	}
}

func (this ClusterModule) Stop() error {

	return nil

}
