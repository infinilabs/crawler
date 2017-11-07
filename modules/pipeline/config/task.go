package config

import (
	"github.com/infinitbyte/gopa/core/model"
)

// TaskConfig defines crawler related configs
type TaskConfig struct {
	Name string `json:"name,omitempty" config:"name"`

	MaxGoRoutine int `config:"max_go_routine"`

	//Speed Control
	ThresholdInMs int `config:"threshold_in_ms"`

	//Timeout Control
	TimeoutInMs int `config:"timeout_in_ms"`

	DefaultPipelineConfig *model.PipelineConfig `config:"default_config"`
}
