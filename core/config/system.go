package config

import (
	"errors"
	"fmt"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/medcl/gopa/core/util"
	"os"
	"path"
)

type ClusterConfig struct {
	Name  string `config:"name"`
	Seeds string `config:"seeds"`
}

type NetworkConfig struct {
	Host string `config:"host"`
}

type NodeConfig struct {
	Name string `config:"name"`
}

type PathConfig struct {
	Data string `config:"data"`
	Log  string `config:"logs"`
	Cert string `config:"certs"`
}

//high priority config, init from the environment or startup, can't be changed on the fly, need to restart
type SystemConfig struct {
	ConfigFile string

	ClusterConfig ClusterConfig `config:"cluster"`

	NetworkConfig NetworkConfig `config:"network"`

	NodeConfig NodeConfig `config:"node"`

	PathConfig PathConfig `config:"path"`

	APIBinding     string `config:"api_bind"`
	HttpBinding    string `config:"http_bind"`
	ClusterBinding string `config:"cluster_bind"`

	AllowMultiInstance bool `config:"multi_instance"`
	TLSEnabled bool `config:"tls_enabled"`
}

var (
	defaultSystemConfig = SystemConfig{
		ClusterConfig: ClusterConfig{
			Name: "gopa",
		},
		NetworkConfig: NetworkConfig{
			Host: "127.0.0.1",
		},
		NodeConfig: NodeConfig{
			Name: util.RandomPickName(),
		},
		PathConfig: PathConfig{
			Data: "data",
			Log:  "log",
			Cert: "cert",
		},

		APIBinding:         ":8001",
		HttpBinding:        ":9001",
		ClusterBinding:     ":13001",
		AllowMultiInstance: false,
	}
)

func LoadSystemConfig(cfgFile string) SystemConfig {
	cfg := defaultSystemConfig
	cfg.ConfigFile = cfgFile
	config, err := yaml.NewConfigWithFile(cfgFile, ucfg.PathSep("."))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = config.Unpack(&cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg.init()
	return cfg
}

func (this *SystemConfig) init() {
	os.MkdirAll(this.GetDataDir(), 0777)
	os.MkdirAll(this.PathConfig.Log, 0777)
}

func (this *SystemConfig) GetDataDir() string {
	if this.AllowMultiInstance == false {
		return path.Join(this.PathConfig.Data, this.ClusterConfig.Name, "nodes", "0")
	}
	//TODO auto select next nodes folder, eg: nodes/1 nodes/2
	panic(errors.New("not supported yet"))
}

func GetDefaultSystemConfig()SystemConfig  {
	return defaultSystemConfig
}