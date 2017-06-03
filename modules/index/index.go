package index

import (
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/modules/config"
	."github.com/medcl/gopa/core/config"
)

type IndexModule struct {
}

func (this IndexModule) Name() string {
	return "Index"
}

var signalChannel chan bool

func (this IndexModule) Start(cfg *Config) {
	signalChannel = make(chan bool, 1)
	go func() {
		for {
			select {
			case <-signalChannel:
				log.Trace("indexer exited")
				return
			default:
				log.Trace("waiting index signal")
				er, v := queue.Pop(config.IndexChannel)
				log.Trace("got index signal, ", string(v))
				if er != nil {
					log.Trace(er)
					continue
				}
				//indexing to es or blevesearch

			}

		}
	}()
}

func (this IndexModule) Stop() error {
	signalChannel <- true
	return nil
}
