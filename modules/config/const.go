package config

import (
	. "github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/filter"
)
const (
	FetchChannel QueueKey = "fetch"
	CheckChannel QueueKey = "check"
	DispatcherChannel QueueKey = "dispatcher"


	CheckFilter filter.FilterKey = "check_filter"
	FetchFilter filter.FilterKey = "fetch_filter"
)

const TaskBucketKey string = "Task"
const StatsBucketKey string = "Stats"
const SnapshotBucketKey string = "Snapshot"
const SnapshotMappingBucketKey string = "SnapshotMapping"