package index

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/config"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/modules/config"
	. "github.com/medcl/gopa/modules/index/elasticsearch"
)

type IndexModule struct {
}

func (this IndexModule) Name() string {
	return "Index"
}

var signalChannel chan bool

var (
	defaultESConfig = ElasticsearchConfig{
		Endpoint: "http://localhost:9200",
		Index:    "gopa",
	}
)

func (this IndexModule) Start(cfg *Config) {

	elasticsearchConfig := defaultESConfig
	cfg.Unpack(&elasticsearchConfig)

	signalChannel = make(chan bool, 1)
	client := ElasticsearchClient{Endpoint: elasticsearchConfig.Endpoint, Index: elasticsearchConfig.Index}
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
				doc := model.IndexDocument{}
				err := json.Unmarshal(v, &doc)
				if err != nil {
					panic(err)
				}

				client.IndexDoc(doc.Id, doc.Source)
			}

		}
	}()
}

func (this IndexModule) Stop() error {
	signalChannel <- true
	return nil
}
