package config

import (
	"github.com/medcl/gopa/core/filter"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/model"
	"github.com/medcl/gopa/core/queue"
)

const (
	FetchChannel      queue.QueueKey = "fetch"
	CheckChannel      queue.QueueKey = "check"
	DispatcherChannel queue.QueueKey = "dispatcher"

	DispatchFilter filter.FilterKey = "dispatch_filter"
	CheckFilter    filter.FilterKey = "check_filter"
	FetchFilter    filter.FilterKey = "fetch_filter"
)

const TaskBucketKey string = "Task"
const StatsBucketKey string = "Stats"
const SnapshotBucketKey string = "Snapshot"
const SnapshotMappingBucketKey string = "SnapshotMapping"

const PhraseChecker model.TaskPhrase = 1 //check url validation
const PhraseCrawler model.TaskPhrase = 2 //going to fetch
const PhraseParse model.TaskPhrase = 3   //parse content
const PhraseUpdate model.TaskPhrase = 4  //scheduled update

const REGISTER_BOLTDB global.RegisterKey = "REGISTER_BOLTDB"
