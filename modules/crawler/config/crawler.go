package config

type CrawlerConfig struct {
	MaxGoRoutine         int  `config:"max_go_routine"`
	LoadPendingFetchJobs bool `config:"load_pending_fetch_from_file"` //fetch url parse and extracted from saved page,load data from:"pending_fetch.urls"
	//Fetch Speed Control
	FetchDelayThresholdInMs int `config:"fetch_delay_threshold_ms"`
}

var (
	defaultCrawlerConfig = CrawlerConfig{
		MaxGoRoutine:            1,
		FetchDelayThresholdInMs: 0,
	}
)

func GetDefaultCrawlerConfig() CrawlerConfig {
	return defaultCrawlerConfig
}
