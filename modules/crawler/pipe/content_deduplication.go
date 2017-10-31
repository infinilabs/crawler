package pipe

import (
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/filter"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/modules/config"
)

// ContentDeduplicationJoint used to check the hash of page body, if duplicated hash already exists, will break the pipeline
type ContentDeduplicationJoint struct {
	model.Parameters
}

// Name return: content_deduplication
func (joint ContentDeduplicationJoint) Name() string {
	return "content_deduplication"
}

// Process the content hash Deduplication
func (joint ContentDeduplicationJoint) Process(c *model.Context) error {
	task := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := c.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	if snapshot.Hash != "" {

		log.Trace("check content deduplication, ", task.Url)

		snapshot.Url = task.Url
		snapshot.TaskID = task.ID

		exist, snapshotId, url := checkByHash(snapshot, c)

		msg := fmt.Sprintf("same content hash found, %s, %s, %s, duplicated with snapshotId: %s , url: %s", task.ID, task.Url, snapshot.Hash, snapshotId, url)

		if exist {
			task.Status = model.TaskDuplicated
			task.NextCheck = nil
			c.End(msg)
			return errors.New(msg)
		}

	}

	return nil
}

func checkByHash(snapshot *model.Snapshot, c *model.Context) (bool, string, string) {

	hash := snapshot.Hash

	//Check local hash first
	if c.GetBool("check_filter", false) {
		exist, _ := filter.CheckThenAdd(config.ContentHashFilter, []byte(hash))
		if exist {
			return true, "local_filter_cache", ""
		}
	}

	//Check hash from db
	items, err := model.GetSnapshotByField("hash", hash)

	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		for _, v := range items {
			log.Tracef("%s vs  %s , %s vs %s", v.Url, snapshot.Url, snapshot.TaskID, v.TaskID)
			if v.Url != snapshot.Url && v.TaskID != snapshot.TaskID {
				return true, v.ID, v.Url
			}
		}
	}

	return false, "", ""
}
