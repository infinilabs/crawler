package config

import (
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/util"
	"io/ioutil"
	"path"
)

// ClusterConfig stores cluster settings
type ClusterConfig struct {
	Name  string   `config:"name"`
	Seeds []string `config:"seeds"`
}

// NetworkConfig stores network settings
type NetworkConfig struct {
	Host string `config:"host"`
}

// NodeConfig stores node settings
type NodeConfig struct {
	Name string `config:"name"`
}

// PathConfig stores path settings
type PathConfig struct {
	Data string `config:"data"`
	Log  string `config:"logs"`
	Cert string `config:"certs"`
}

// SystemConfig is a high priority config, init from the environment or startup, can't be changed on the fly, need to restart to make config apply
type SystemConfig struct {
	workingDir string `config:"-"`

	ConfigFile string

	ClusterConfig ClusterConfig `config:"cluster"`

	NetworkConfig NetworkConfig `config:"network"`

	NodeConfig NodeConfig `config:"node"`

	PathConfig PathConfig `config:"path"`

	APIBinding     string `config:"api_bind"`
	HTTPBinding    string `config:"http_bind"`
	ClusterBinding string `config:"cluster_bind"`

	AllowMultiInstance bool `config:"allow_multi_instance"`
	MaxNumOfInstance   int  `config:"max_num_of_instances"`
	TLSEnabled         bool `config:"tls_enabled"`
}

// GetWorkingDir returns root working dir of gopa instance
func (sysconfig *SystemConfig) GetWorkingDir() string {
	if sysconfig.workingDir != "" {
		return sysconfig.workingDir
	}

	if !sysconfig.AllowMultiInstance {
		sysconfig.workingDir = path.Join(sysconfig.PathConfig.Data, sysconfig.ClusterConfig.Name, "nodes", "0")
		return sysconfig.workingDir
	}

	//auto select next nodes folder, eg: nodes/1 nodes/2
	i := 0
	if sysconfig.MaxNumOfInstance < 1 {
		sysconfig.MaxNumOfInstance = 5
	}
	for j := 0; j < sysconfig.MaxNumOfInstance; j++ {
		p := path.Join(sysconfig.PathConfig.Data, sysconfig.ClusterConfig.Name, "nodes", util.IntToString(i))
		lockFile := path.Join(p, ".lock")
		if !util.FileExists(lockFile) {
			sysconfig.workingDir = p
			return sysconfig.workingDir
		}

		//check if pid is alive
		b, err := ioutil.ReadFile(lockFile)
		if err != nil {
			panic(err)
		}
		pid, err := util.ToInt(string(b))
		if err != nil {
			panic(err)
		}
		if pid <= 0 {
			panic(errors.New("invalid pid"))
		}

		procExists := util.CheckProcessExists(pid)
		if !procExists {
			util.FileDelete(lockFile)
			log.Debug("dead process with broken lock file, removed: ", lockFile)
			sysconfig.workingDir = p
			return p
		}

		i++
	}
	panic(fmt.Errorf("reach max num of instances on this node, limit is: %v", sysconfig.MaxNumOfInstance))

}
