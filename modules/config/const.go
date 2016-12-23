package config

import . "github.com/medcl/gopa/core/queue"
const (
	FetchChannel QueueKey = "fetch"
	CheckChannel QueueKey = "check"
)

const TaskBucketKey string = "Task"
const StatsBucketKey string = "Stats"
const SnapshotBucketKey string = "Snapshot"
const SnapshotMappingBucketKey string = "SnapshotMapping"