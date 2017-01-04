package cluster

import (
	. "github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/modules/cluster/discovery/raft/store"
	"github.com/medcl/gopa/modules/cluster/discovery/raft/http"
	log "github.com/cihub/seelog"
)

type ClusterModule struct {

}

func (this ClusterModule) Name() string {
	return "Cluster"
}

func (this ClusterModule) Start(env *Env) {

	s := store.New()

	//init raft
	if err := s.Open(); err != nil {
		log.Errorf("failed to open store: %s", err.Error())
		panic(err)
	}

	//init http
	h := httpd.New(s)
	if err := h.Start(); err != nil {
		log.Errorf("failed to start HTTP service: %s", err.Error())
		panic(err)
	}
}



func (this ClusterModule) Stop() error {

	return nil

}
