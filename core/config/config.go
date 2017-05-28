/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	log "github.com/cihub/seelog"
	cfg "github.com/robfig/config"
)

var loadingConfig *cfg.Config
var runtimeConfig RuntimeConfig

func GetStringConfig(configSection string, configKey string, defaultValue string) string {
	if loadingConfig == nil {
		log.Trace("loadingConfig is nil,just return")
		return defaultValue
	}

	//loading or initializing bloom filter
	value, error := loadingConfig.String(configSection, configKey)
	if error != nil {
		value = defaultValue
	}
	log.Trace("get config value,", configSection, ".", configKey, ":", value)
	return value
}

func GetFloatConfig(configSection string, configKey string, defaultValue float64) float64 {
	if loadingConfig == nil {
		log.Trace("loadingConfig is nil,just return")
		return defaultValue
	}

	//loading or initializing bloom filter
	value, error := loadingConfig.Float(configSection, configKey)
	if error != nil {
		value = defaultValue
	}
	log.Trace("get config value,", configSection, ".", configKey, ":", value)
	return value
}

func GetIntConfig(configSection string, configKey string, defaultValue int) int {
	if loadingConfig == nil {
		log.Trace("loadingConfig is nil,just return")
		return defaultValue
	}

	//loading or initializing bloom filter
	value, error := loadingConfig.Int(configSection, configKey)
	if error != nil {
		value = defaultValue
	}
	log.Trace("get config value,", configSection, ".", configKey, ":", value)
	return value
}

func GetBoolConfig(configSection string, configKey string, defaultValue bool) bool {
	if loadingConfig == nil {
		log.Trace("loadingConfig is nil,just return")
		return defaultValue
	}

	//loading or initializing bloom filter
	value, error := loadingConfig.Bool(configSection, configKey)
	if error != nil {
		value = defaultValue
	}
	log.Trace("get config value,", configSection, ".", configKey, ":", value)
	return value
}
