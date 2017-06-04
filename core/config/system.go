package config

import (
	"errors"
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
	TLSEnabled         bool `config:"tls_enabled"`
}

func (this *SystemConfig) GetDataDir() string {
	if this.AllowMultiInstance == false {
		return path.Join(this.PathConfig.Data, this.ClusterConfig.Name, "nodes", "0")
	}
	//TODO auto select next nodes folder, eg: nodes/1 nodes/2
	panic(errors.New("not supported yet"))
}
