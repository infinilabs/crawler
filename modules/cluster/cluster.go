package cluster

import (
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/modules/cluster/discovery/raft"
)

type ClusterModule struct {
}

func (module ClusterModule) Name() string {
	return "Cluster"
}

func (module ClusterModule) Start(cfg *Config) {

	s := raft.New()
	if err := s.Open(); err != nil {
		log.Errorf("failed to open raft: %s", err.Error())
		panic(err)
	}
}

func (module ClusterModule) Stop() error {

	return nil

}
