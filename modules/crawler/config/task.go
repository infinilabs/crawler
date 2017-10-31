package config

import (
	"github.com/infinitbyte/gopa/core/model"
)

// TaskConfig defines crawler related configs
type TaskConfig struct {
	MaxGoRoutine int `config:"max_go_routine"`
	//Fetch Speed Control
	FetchThresholdInMs int `config:"fetch_threshold_ms"`

	DefaultPipelineConfig *model.PipelineConfig `config:"default_pipeline_config"`
}
