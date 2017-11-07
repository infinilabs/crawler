package config

import (
	"github.com/infinitbyte/gopa/core/model"
)

// PipeConfig defines crawler related configs
type PipeConfig struct {
	Name string `json:"name,omitempty" config:"name"`

	MaxGoRoutine int `config:"max_go_routine"`

	//Speed Control
	ThresholdInMs int `config:"threshold_in_ms"`

	//Timeout Control
	TimeoutInMs int `config:"timeout_in_ms"`

	DefaultConfig *model.PipelineConfig `config:"default_config"`
}
