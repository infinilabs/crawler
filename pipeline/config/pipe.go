package config

import (
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/config"
)

// getDefaultCrawlerConfig return a default PipeRunnerConfig
func GetDefaultPipeConfig() []pipeline.PipeConfig {

	crawler := pipeline.PipelineConfig{}
	start := pipeline.JointConfig{}
	start.Enabled = true
	start.JointName = "init_task"
	crawler.StartJoint = &start
	save := pipeline.JointConfig{}
	save.Enabled = true
	save.JointName = "save_task"

	urlNormalization := pipeline.JointConfig{}
	urlNormalization.Enabled = true
	urlNormalization.JointName = "url_normalization"
	urlNormalization.Parameters = util.MapStr{
		"follow_all_domain": false,
		"follow_sub_domain": true,
	}

	fetchJoint := pipeline.JointConfig{}
	fetchJoint.Enabled = true
	fetchJoint.JointName = "fetch"

	parse := pipeline.JointConfig{}
	parse.Enabled = true
	parse.JointName = "parse"

	html2text := pipeline.JointConfig{}
	html2text.Enabled = true
	html2text.JointName = "html2text"

	hash := pipeline.JointConfig{}
	hash.Enabled = true
	hash.JointName = "hash"

	updateCheckTime := pipeline.JointConfig{}
	updateCheckTime.Enabled = true
	updateCheckTime.JointName = "update_check_time"

	contentDeduplication := pipeline.JointConfig{}
	contentDeduplication.Enabled = true
	contentDeduplication.JointName = "content_deduplication"

	langDetect := pipeline.JointConfig{}
	langDetect.Enabled = true
	langDetect.JointName = "lang_detect"

	index := pipeline.JointConfig{}
	index.Enabled = true
	index.JointName = "index"

	saveSnapshot := pipeline.JointConfig{}
	saveSnapshot.Enabled = true
	saveSnapshot.JointName = "save_snapshot_db"

	crawler.EndJoint = &save
	crawler.ProcessJoints = []*pipeline.JointConfig{
		&urlNormalization,
		&fetchJoint,
		&parse,
		&html2text,
		&hash,
		&contentDeduplication,
		&updateCheckTime,
		&langDetect,
		&saveSnapshot,
		&index,
	}

	defaultCrawlerConfig := pipeline.PipeConfig{
		Name:          "crawler",
		Enabled:       true,
		MaxGoRoutine:  10,
		TimeoutInMs:   30000,
		InputQueue:    config.FetchChannel,
		ThresholdInMs: 0,
		DefaultConfig: &crawler,
	}

	checker := pipeline.PipelineConfig{}
	checker.StartJoint = &start
	save = pipeline.JointConfig{}
	save.Enabled = true
	save.JointName = "save_task"
	save.Parameters = util.MapStr{
		"is_create": true,
	}

	urlFilter := pipeline.JointConfig{}
	urlFilter.Enabled = true
	urlFilter.JointName = "url_filter"

	urlCheckFilter := pipeline.JointConfig{}
	urlCheckFilter.Enabled = true
	urlCheckFilter.JointName = "filter_check"

	taskDeduplication := pipeline.JointConfig{}
	taskDeduplication.Enabled = true
	taskDeduplication.JointName = "task_deduplication"

	checker.EndJoint = &save
	checker.ProcessJoints = []*pipeline.JointConfig{
		&urlNormalization,
		&urlFilter,
		&urlCheckFilter,
		&taskDeduplication,
	}

	defaultCheckerConfig := pipeline.PipeConfig{
		Name:          "checker",
		Enabled:       true,
		InputQueue:    config.CheckChannel,
		MaxGoRoutine:  10,
		ThresholdInMs: 0,
		TimeoutInMs:   5000,
		DefaultConfig: &checker,
	}

	result := []pipeline.PipeConfig{defaultCrawlerConfig, defaultCheckerConfig}

	return result
}
