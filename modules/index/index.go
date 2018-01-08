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
	"github.com/infinitbyte/gopa/modules/config"
	"github.com/infinitbyte/gopa/modules/index/ui"
	common "github.com/infinitbyte/gopa/modules/index/ui/common"
	"runtime"
)

type IndexModule struct {
}

func (this IndexModule) Name() string {
	return "Index"
}

var signalChannel chan bool

var (
	defaultConfig = common.IndexConfig{
		Elasticsearch: &core.ElasticsearchConfig{
			Endpoint:    "http://localhost:9200",
			IndexPrefix: "gopa-",
		},
		UIConfig: &common.UIConfig{
			Enabled:     true,
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
		api.HandleUIMethod(api.GET, "/m/", ui.MobileIndexPageAction)
		api.HandleUIMethod(api.GET, "/m/ajax_more_item/", ui.MobileAJAXMoreItemAction)
		api.HandleUIMethod(api.GET, "/snapshot/:id", api.NeedLogin("", ui.GetSnapshotPayloadAction))
		api.HandleUIMethod(api.GET, "/suggest/", ui.SuggestAction)
	}

	go func() {
		defer func() {

			if !global.Env().IsDebug {
				if r := recover(); r != nil {

					if r == nil {
						return
					}
					var v string
					switch r.(type) {
					case error:
						v = r.(error).Error()
					case runtime.Error:
						v = r.(runtime.Error).Error()
					case string:
						v = r.(string)
					}
					log.Error("error in indexer,", v)
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
					log.Error(er)
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
