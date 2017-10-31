package generator

import (
	. "github.com/infinitbyte/gopa/core/config"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
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
				queue.Push(config.CheckChannel, model.NewTaskSeed(generatorConfig.TaskUrl, generatorConfig.TaskUrl, 0, 0).MustGetBytes())
			}

			if generatorConfig.TaskID != "" {
				queue.Push(config.FetchChannel, model.EncodeFetchTask(generatorConfig.TaskID, "", ""))
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (module GeneratorModule) Stop() error {

	return nil

}
