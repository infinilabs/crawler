package index

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/global"
	api "github.com/infinitbyte/gopa/core/http"
	core "github.com/infinitbyte/gopa/core/index"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	"github.com/infinitbyte/gopa/modules/index/ui"
	cfg "github.com/infinitbyte/gopa/modules/index/ui/config"
	"runtime"
)

type IndexModule struct {
}

func (this IndexModule) Name() string {
	return "Index"
}

var signalChannel chan bool

var (
	defaultConfig = cfg.IndexConfig{
		Elasticsearch: &core.ElasticsearchConfig{
			Endpoint:    "http://localhost:9200",
			IndexPrefix: "gopa-",
		},
		UIConfig: &cfg.UIConfig{
			SiteName:    "Gopa",
			SiteFavicon: "/static/assets/img/favicon.ico",
			SiteLogo:    "/static/assets/img/logo.svg",
		},
	}
)

func (module IndexModule) Start(cfg *Config) {

	indexConfig := defaultConfig
	cfg.Unpack(&indexConfig)

	signalChannel = make(chan bool, 1)
	client := core.ElasticsearchClient{Config: indexConfig.Elasticsearch}

	//register UI
	if indexConfig.UIConfig.Enabled {
		ui := ui.UserUI{}
		ui.Config = indexConfig.UIConfig
		ui.SearchClient = &core.ElasticsearchClient{Config: indexConfig.Elasticsearch}
		api.HandleUIMethod(api.GET, "/", ui.IndexPageAction)
		api.HandleUIMethod(api.GET, "/snapshot/:id", ui.GetSnapshotPayloadAction)
	}

	go func() {
		defer func() {

			if !global.Env().IsDebug {
				if r := recover(); r != nil {
					if e, ok := r.(runtime.Error); ok {
						log.Error("index: ", util.GetRuntimeErrorMessage(e))
					}
					log.Error("error in indexer")
				}
			}
		}()

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

				client.Index(doc.Index, doc.ID, doc.Source)
			}

		}
	}()
}

func (module IndexModule) Stop() error {
	signalChannel <- true
	return nil
}
