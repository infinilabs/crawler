package config

import (
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/util"
)

// PipeConfig defines crawler related configs
type PipeConfig struct {
	Name string `json:"name,omitempty" config:"name"`

	Enabled bool `json:"enabled,omitempty" config:"enabled"`

	MaxGoRoutine int `config:"max_go_routine"`

	//Speed Control
	ThresholdInMs int `config:"threshold_in_ms"`

	//Timeout Control
	TimeoutInMs int `config:"timeout_in_ms"`

	DefaultConfig *model.PipelineConfig `config:"default_config"`
}

// getDefaultCrawlerConfig return a default PipeConfig
func GetDefaultPipeConfig() []PipeConfig {

	crawler := model.PipelineConfig{}
	start := model.JointConfig{}
	start.Enabled = true
	start.JointName = "init_task"
	crawler.StartJoint = &start
	save := model.JointConfig{}
	save.Enabled = true
	save.JointName = "save_task"

	urlNormalization := model.JointConfig{}
	urlNormalization.Enabled = true
	urlNormalization.JointName = "url_normalization"
	urlNormalization.Parameters = util.MapStr{
		"follow_all_domain": false,
		"follow_sub_domain": true,
	}

	fetchJoint := model.JointConfig{}
	fetchJoint.Enabled = true
	fetchJoint.JointName = "fetch"

	parse := model.JointConfig{}
	parse.Enabled = true
	parse.JointName = "parse"

	html2text := model.JointConfig{}
	html2text.Enabled = true
	html2text.JointName = "html2text"

	hash := model.JointConfig{}
	hash.Enabled = true
	hash.JointName = "hash"

	updateCheckTime := model.JointConfig{}
	updateCheckTime.Enabled = true
	updateCheckTime.JointName = "update_check_time"

	contentDeduplication := model.JointConfig{}
	contentDeduplication.Enabled = true
	contentDeduplication.JointName = "content_deduplication"

	langDetect := model.JointConfig{}
	langDetect.Enabled = true
	langDetect.JointName = "lang_detect"

	index := model.JointConfig{}
	index.Enabled = true
	index.JointName = "index"

	saveSnapshot := model.JointConfig{}
	saveSnapshot.Enabled = true
	saveSnapshot.JointName = "save_snapshot_db"

	crawler.EndJoint = &save
	crawler.ProcessJoints = []*model.JointConfig{
		&urlNormalization,
		&fetchJoint,
		&parse,
		&html2text,
		&hash,
		&updateCheckTime,
		&contentDeduplication,
		&langDetect,
		&saveSnapshot,
		&index,
	}

	defaultCrawlerConfig := PipeConfig{
		Name:          "crawler",
		Enabled:       true,
		MaxGoRoutine:  10,
		TimeoutInMs:   30000,
		ThresholdInMs: 0,
		DefaultConfig: &crawler,
	}

	checker := model.PipelineConfig{}
	checker.StartJoint = &start
	save = model.JointConfig{}
	save.Enabled = true
	save.JointName = "save_task"
	save.Parameters = util.MapStr{
		"is_create": true,
	}

	urlFilter := model.JointConfig{}
	urlFilter.Enabled = true
	urlFilter.JointName = "url_filter"

	urlCheckFilter := model.JointConfig{}
	urlCheckFilter.Enabled = true
	urlCheckFilter.JointName = "filter_check"

	taskDeduplication := model.JointConfig{}
	taskDeduplication.Enabled = true
	taskDeduplication.JointName = "task_deduplication"

	checker.EndJoint = &save
	checker.ProcessJoints = []*model.JointConfig{
		&urlNormalization,
		&urlFilter,
		&urlCheckFilter,
		&taskDeduplication,
	}

	defaultCheckerConfig := PipeConfig{
		Name:          "checker",
		Enabled:       true,
		MaxGoRoutine:  10,
		ThresholdInMs: 0,
		TimeoutInMs:   5000,
		DefaultConfig: &checker,
	}

	result := []PipeConfig{defaultCrawlerConfig, defaultCheckerConfig}

	return result
}
