package config

import (
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/global"
)

const (
	FetchChannel      string = "fetch"
	UpdateChannel     string = "update"
	CheckChannel      string = "check"
	DispatcherChannel string = "dispatch"
	IndexChannel      string = "index"

	DispatchFilter    string = "dispatch_filter"
	CheckFilter       string = "check_filter"
	FetchFilter       string = "fetch_filter"
	ContentHashFilter string = "content_hash_filter"
)

const KVBucketKey string = "Kv"
const TaskBucketKey string = "Task"
const StatsBucketKey string = "Stats"
const SnapshotBucketKey string = "Snapshot"
const ScreenshotBucketKey string = "Screenshot"
const SnapshotMappingBucketKey string = "SnapshotMapping"

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
