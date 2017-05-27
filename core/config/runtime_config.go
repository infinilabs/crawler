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

type ChannelConfig struct {
	CheckerChanBuffSize int `checker_chan_buff_size`
	FetchChanBuffSize   int `fetch_chan_buff_size`
}

func (this *ChannelConfig) Init() *ChannelConfig {
	this.CheckerChanBuffSize = 0
	this.FetchChanBuffSize = 0
	return this
}

type StatsdConfig struct {
	Host string `host`
}

type IndexingConfig struct {
	Host  string `host`
	Index string `index`
}

func (this *IndexingConfig) Init() *IndexingConfig {
	this.Host = "http://127.0.0.1:9200"
	this.Index = "gopa"
	return this
}

type SaveConfig struct {
	DefaultExtension string
}

func (this *SaveConfig) Init() *SaveConfig {
	this.DefaultExtension = ".html"
	return this
}

type CrawlerConfig struct {
	Enabled              bool `enabled`
	LoadPendingFetchJobs bool `load_pending_fetch_from_file` //fetch url parse and extracted from saved page,load data from:"pending_fetch.urls"

}

func (this *CrawlerConfig) Init() *CrawlerConfig {
	this.Enabled = true
	return this
}

type ParserConfig struct {
	Enabled                          bool `enabled`
	ParseUrlsFromSavedFileLog        bool `parse_file_log`
	ReParseUrlsFromPreviousSavedPage bool `reparse_file_log` //extract urls from previous saved page

}

type RuledFetchConfig struct {
	UrlTemplate        string
	From               int
	To                 int
	Step               int
	LinkExtractPattern string
	LinkTemplate       string
}

type Test struct {
	Name     string
	LastName string `last_name`
}

type RuntimeConfig struct {
	Test *Test

	IndexingConfig *IndexingConfig `indexing`

	ChannelConfig *ChannelConfig `channel`

	//crawler config
	CrawlerConfig *CrawlerConfig `crawler`
	ParserConfig  *ParserConfig  `parser`

	//task
	TaskConfig *TaskConfig `task.default`

	RuledFetchConfig *RuledFetchConfig

	//splitter of joined array string
	//ArrayStringSplitter string

	//StoreWebPageTogether bool

	MaxGoRoutine int `max_go_routine`

	//switch config
	//ParseUrlsFromSavedFileLog      bool
	LoadTemplatedFetchJob bool
	//ParseUrlsFromPreviousSavedPage bool //extract urls from previous saved page
	LoadRuledFetchJob bool //extract urls from previous saved page
	HttpEnabled       bool //extract urls from previous saved page

	WalkBloomFilterFileName         string
	FetchBloomFilterFileName        string
	ParseBloomFilterFileName        string
	PendingFetchBloomFilterFileName string
}
