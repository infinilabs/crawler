package index

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	core "github.com/infinitbyte/gopa/core/index"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/modules/config"
)

type IndexModule struct {
}

type IndexConfig struct {
	Elasticsearch core.ElasticsearchConfig `config:"elasticsearch"`
}

func (this IndexModule) Name() string {
	return "Index"
}

var signalChannel chan bool

var (
	defaultConfig = IndexConfig{
		Elasticsearch: core.ElasticsearchConfig{
			Endpoint:    "http://localhost:9200",
			IndexPrefix: "gopa-",
		},
	}
)

func (module IndexModule) Start(cfg *Config) {

	indexConfig := defaultConfig
	cfg.Unpack(&indexConfig)

	signalChannel = make(chan bool, 1)
	client := core.ElasticsearchClient{Config: indexConfig.Elasticsearch}
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

				client.Index(doc.Index, doc.Id, doc.Source)
			}

		}
	}()
}

func (module IndexModule) Stop() error {
	signalChannel <- true
	return nil
}
