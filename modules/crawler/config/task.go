package config

import "github.com/infinitbyte/gopa/core/pipeline"

// TaskConfig defines crawler related configs
type TaskConfig struct {
	MaxGoRoutine int `config:"max_go_routine"`
	//Fetch Speed Control
	FetchThresholdInMs int `config:"fetch_threshold_ms"`

	DefaultPipelineConfig *pipeline.PipelineConfig `config:"default_pipeline_config"`
}
