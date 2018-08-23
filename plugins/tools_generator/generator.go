package tools_generator

import (
	log "github.com/cihub/seelog"
	. "github.com/infinitbyte/framework/core/config"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/queue"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
	"time"
)

type GeneratorPlugin struct {
}

func (plugin GeneratorPlugin) Name() string {
	return "Generator"
}

var generatorConfig = struct {
	TaskID  string `config:"task_id"`
	TaskUrl string `config:"task_url"`
}{}

func (plugin GeneratorPlugin) Setup(cfg *Config) {

	cfg.Unpack(&generatorConfig)

}

func (plugin GeneratorPlugin) Start() error {
	go func() {
		for {
			if generatorConfig.TaskUrl != "" {
				context := pipeline.Context{IgnoreBroken: true}
				context.Set(model.CONTEXT_TASK_URL, generatorConfig.TaskUrl)
				err := queue.Push(config.CheckChannel, util.ToJSONBytes(context))
				if err != nil {
					log.Error(err)
				}
			}

			if generatorConfig.TaskID != "" {

				context := pipeline.Context{}
				context.Set(model.CONTEXT_TASK_ID, generatorConfig.TaskID)
				err := queue.Push(config.FetchChannel, util.ToJSONBytes(context))
				if err != nil {
					log.Error(err)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return nil
}
func (plugin GeneratorPlugin) Stop() error {
	return nil
}
