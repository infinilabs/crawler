package config

import (
	"errors"
	"fmt"
	"github.com/infinitbyte/gopa/core/util"
	"path"
)

type ClusterConfig struct {
	Name  string   `config:"name"`
	Seeds []string `config:"seeds"`
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
	workingDir string `config:"-"`

	ConfigFile string

	ClusterConfig ClusterConfig `config:"cluster"`

	NetworkConfig NetworkConfig `config:"network"`

	NodeConfig NodeConfig `config:"node"`

	PathConfig PathConfig `config:"path"`

	APIBinding     string `config:"api_bind"`
	HttpBinding    string `config:"http_bind"`
	ClusterBinding string `config:"cluster_bind"`

	AllowMultiInstance bool `config:"allow_multi_instance"`
	MaxNumOfInstance   int  `config:"max_num_of_instances"`
	TLSEnabled         bool `config:"tls_enabled"`
}

func (this *SystemConfig) GetDataDir() string {
	if this.workingDir != "" {
		return this.workingDir
	}

	if this.AllowMultiInstance == false {
		this.workingDir = path.Join(this.PathConfig.Data, this.ClusterConfig.Name, "nodes", "0")
		return this.workingDir
	} else {
		//auto select next nodes folder, eg: nodes/1 nodes/2
		i := 0
		if this.MaxNumOfInstance < 1 {
			this.MaxNumOfInstance = 5
		}
		for j := 0; j < this.MaxNumOfInstance; j++ {
			p := path.Join(this.PathConfig.Data, this.ClusterConfig.Name, "nodes", util.IntToString(i))
			if !util.FileExists(path.Join(p, ".lock")) {
				this.workingDir = p
				return this.workingDir
			}
			i++
		}
		panic(errors.New(fmt.Sprintf("reach max num of instances on this node, limit is: %v", this.MaxNumOfInstance)))
	}
}
