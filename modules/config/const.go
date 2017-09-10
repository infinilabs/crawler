package config

import (
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/filter"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/queue"
)

const (
	FetchChannel      queue.QueueKey = "fetch"
	CheckChannel      queue.QueueKey = "check"
	DispatcherChannel queue.QueueKey = "dispatch"
	IndexChannel      queue.QueueKey = "index"

	DispatchFilter    filter.Key = "dispatch_filter"
	CheckFilter       filter.Key = "check_filter"
	FetchFilter       filter.Key = "fetch_filter"
	ContentHashFilter filter.Key = "content_hash_filter"
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
