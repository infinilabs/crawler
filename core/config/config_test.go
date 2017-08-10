package config

import (
	"fmt"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/magiconair/properties/assert"
	"os"
	"testing"
)

// Defines struct to read config from
type ExampleConfig struct {
	Counter int32 `config:"counter" validate:"min=0, max=9"`
}

// Defines default config option
var (
	defaultConfig = ExampleConfig{
		Counter: 4,
	}
)

func TestLoadDefaultCfg(t *testing.T) {

	path := "config_test.yml"
	appConfig := defaultConfig // copy default config so it's not overwritten
	config, err := yaml.NewConfigWithFile(path, ucfg.PathSep("."))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	assert.Equal(t, appConfig.Counter, int32(4))
	fmt.Println(appConfig.Counter)

	err = config.Unpack(&appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	assert.Equal(t, appConfig.Counter, int32(5))

	fmt.Println(appConfig.Counter)

}

type globalConfig struct {
	Modules []*Config `config:"modules"`
}

type crawlerConfig struct {
	Namespace  string `config:"namespace" validate:"required"`
	LikedCount int    `config:"liked"`
}

var (
	defaultCrawlerConfig = crawlerConfig{
		Namespace:  "default",
		LikedCount: 512,
	}
)

func TestLoadModules(t *testing.T) {
	cfg, _ := LoadFile("config_test.yml")

	config := globalConfig{}

	if err := cfg.Unpack(&config); err != nil {
		fmt.Println(err)
	}

	crawlerCfg := defaultCrawlerConfig

	cf1 := newConfig(t, config.Modules)
	cf1[0].Unpack(&crawlerCfg)

	assert.Equal(t, crawlerCfg.Namespace, "hello world")
	assert.Equal(t, crawlerCfg.LikedCount, 1235)

	parserConfig := struct {
		ID string `config:"parser_id" validate:"required"`
	}{}
	cf1[1].Unpack(&parserConfig)
	fmt.Println(parserConfig.ID)

}

func getModuleName(c *Config) string {
	cfgObj := struct {
		Module string `config:"module"`
	}{}

	if c == nil {
		return ""
	}
	if err := c.Unpack(&cfgObj); err != nil {
		return ""
	}

	return cfgObj.Module
}

func newConfig(t testing.TB, cfgs []*Config) []*Config {
	results := []*Config{}
	for _, cfg := range cfgs {
		//set map for modules and module config
		fmt.Println(getModuleName(cfg))
		fmt.Println(cfg.Enabled())
		config, err := NewConfigFrom(cfg)
		if err != nil {
			t.Fatal(err)
		}
		results = append(results, config)
	}

	return results
}
