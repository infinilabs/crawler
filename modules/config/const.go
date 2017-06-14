package config

import (
	"github.com/medcl/gopa/core/errors"
	"github.com/medcl/gopa/core/filter"
	"github.com/medcl/gopa/core/global"
	"github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/queue"
)

const (
	FetchChannel      queue.QueueKey = "fetch"
	CheckChannel      queue.QueueKey = "check"
	DispatcherChannel queue.QueueKey = "dispatch"
	IndexChannel      queue.QueueKey = "index"

	DispatchFilter filter.FilterKey = "dispatch_filter"
	CheckFilter    filter.FilterKey = "check_filter"
	FetchFilter    filter.FilterKey = "fetch_filter"
)

const KVBucketKey string = "Kv"
const TaskBucketKey string = "Task"
const StatsBucketKey string = "Stats"
const SnapshotBucketKey string = "Snapshot"
const SnapshotMappingBucketKey string = "SnapshotMapping"

const PhraseChecker pipeline.Phrase = 1 //check url validation
const PhraseCrawler pipeline.Phrase = 2 //going to fetch
const PhraseParse pipeline.Phrase = 3   //parse content
const PhraseUpdate pipeline.Phrase = 4  //scheduled update

const REGISTER_BOLTDB global.RegisterKey = "REGISTER_BOLTDB"

const ErrorExitedPipeline errors.ErrorCode = 1000
const ErrorBrokenPipeline errors.ErrorCode = 1001

const STATS_FETCH_TOTAL_COUNT = "fetch.total"
const STATS_FETCH_SUCCESS_COUNT = "fetch.success"
const STATS_FETCH_FAIL_COUNT = "fetch.fail"
const STATS_FETCH_TIMEOUT_COUNT = "fetch.timeout"
const STATS_FETCH_TIMEOUT_IGNORE_COUNT = "fetch.timeout_ignore"

const STATS_STORAGE_FILE_SIZE = "stats.sum.file.size"
const STATS_STORAGE_FILE_COUNT = "stats.sum.file.count"
