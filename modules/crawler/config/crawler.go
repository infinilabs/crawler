package config

type CrawlerConfig struct {
	MaxGoRoutine int `config:"max_go_routine"`
	//Fetch Speed Control
	FetchThresholdInMs int `config:"fetch_threshold_ms"`
}

var (
	defaultCrawlerConfig = CrawlerConfig{
		MaxGoRoutine:       1,
		FetchThresholdInMs: 0,
	}
)

func GetDefaultCrawlerConfig() CrawlerConfig {
	return defaultCrawlerConfig
}
