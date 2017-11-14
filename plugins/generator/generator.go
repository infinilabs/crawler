package generator

import (
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	"time"
)

type GeneratorModule struct {
}

func (module GeneratorModule) Name() string {
	return "Generator"
}

func (module GeneratorModule) Start(cfg *Config) {

	generatorConfig := struct {
		TaskID  string `config:"task_id"`
		TaskUrl string `config:"task_url"`
	}{}

	cfg.Unpack(&generatorConfig)

	go func() {
		for {
			if generatorConfig.TaskUrl != "" {
				context := model.Context{IgnoreBroken: true}
				context.Set(model.CONTEXT_TASK_URL, generatorConfig.TaskUrl)
				queue.Push(config.CheckChannel, util.ToJSONBytes(context))
			}

			if generatorConfig.TaskID != "" {

				context := model.Context{}
				context.Set(model.CONTEXT_TASK_ID, generatorConfig.TaskID)
				queue.Push(config.FetchChannel, util.ToJSONBytes(context))
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (module GeneratorModule) Stop() error {

	return nil

}
