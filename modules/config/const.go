package config

import (
	. "github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/filter"
	"github.com/medcl/gopa/core/model"
)
const (
	FetchChannel QueueKey = "fetch"
	CheckChannel QueueKey = "check"
	DispatcherChannel QueueKey = "dispatcher"


	DispatchFilter filter.FilterKey = "dispatch_filter"
	CheckFilter filter.FilterKey = "check_filter"
	FetchFilter filter.FilterKey = "fetch_filter"
)

const TaskBucketKey string = "Task"
const StatsBucketKey string = "Stats"
const SnapshotBucketKey string = "Snapshot"
const SnapshotMappingBucketKey string = "SnapshotMapping"

const PhraseChecker  model.TaskPhrase =1 //check url validation
const PhraseCrawler  model.TaskPhrase =2 //going to fetch
const PhraseUpdate  model.TaskPhrase =3 //scheduled update
