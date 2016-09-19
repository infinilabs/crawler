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
	. "github.com/medcl/gopa/core/types"
)

type PathConfig struct {
	Data     string  `data`
	Log      string  `log`
	TaskData string
	QueueData string
	WebData  string

	SavedFileLog    string //path of saved files
	PendingFetchLog string //path of pending fetch
	FetchFailedLog  string //path of failed fetch
}

func (this *PathConfig) Init() (*PathConfig)  {
	this.Data="data"
	this.Log="log"
	return this
}

type ClusterConfig struct {
	Name string
}

func (this *ClusterConfig)Init() (*ClusterConfig)  {
	this.Name="gopa"
	return this
}


type LoggingConfig struct {
	Level     string `level`

	//config string of seelog
	ConfigStr string
}

func (this *LoggingConfig)Init() (*LoggingConfig)  {
	this.Level="info"
	return this
}


type IndexingConfig struct {
	Host string `host`
	Index string `index`
}

type IndexingConfig struct {
	Host string `host`
	Index string `index`
}

func (this *IndexingConfig)Init() (*IndexingConfig)  {
	this.Host="http://127.0.0.1:9200"
	this.Index="gopa"
	return this
}

type SaveConfig struct {
	DefaultExtension string
}

func (this *SaveConfig)Init() (*SaveConfig)  {
	this.DefaultExtension=".html"
	return this
}

type CrawlerConfig struct {
	Enabled bool `enabled`
	LoadPendingFetchJobs           bool  `load_pending_fetch_from_file`//fetch url parse and extracted from saved page,load data from:"pending_fetch.urls"

}
type ParserConfig struct {
	Enabled bool `enabled`
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

type RuntimeConfig struct {
	//cluster
	ClusterConfig *ClusterConfig `cluster`

	//logging related
	LoggingConfig *LoggingConfig `logging`
	IndexingConfig *IndexingConfig `indexing`


	//path related
	PathConfig *PathConfig `path`


	//crawler config
	CrawlerConfig *CrawlerConfig `crawler`
	ParserConfig *ParserConfig `parser`


	//task
	TaskConfig *TaskConfig `task.default`

	RuledFetchConfig *RuledFetchConfig

	//splitter of joined array string
	//ArrayStringSplitter string


	//StoreWebPageTogether bool

	MaxGoRoutine int `max_go_routine`

	//switch config
	//ParseUrlsFromSavedFileLog      bool
	LoadTemplatedFetchJob          bool
	//ParseUrlsFromPreviousSavedPage bool //extract urls from previous saved page
	LoadRuledFetchJob              bool //extract urls from previous saved page
	HttpEnabled                    bool //extract urls from previous saved page

	//runtime variables
	Storage Store

	WalkBloomFilterFileName         string
	FetchBloomFilterFileName        string
	ParseBloomFilterFileName        string
	PendingFetchBloomFilterFileName string
}
