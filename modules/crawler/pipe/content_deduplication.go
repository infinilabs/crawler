package pipe

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/filter"
	"github.com/infinitbyte/gopa/core/model"
	api "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/modules/config"
)

// ContentDeduplicationJoint used to check the hash of page body, if duplicated hash already exists, will break the pipeline
type ContentDeduplicationJoint struct {
	api.Parameters
}

// Name return: content_deduplication
func (joint ContentDeduplicationJoint) Name() string {
	return "content_deduplication"
}

// Process the content hash Deduplication
func (joint ContentDeduplicationJoint) Process(c *api.Context) error {
	task := c.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := c.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	if snapshot.Hash != "" {

		log.Trace("check content deduplication, ", task.Url)
		msg := fmt.Sprintf("same content hash found, %s, %s, %s", task.ID, task.Url, snapshot.Hash)

		exist := checkByHash(snapshot.Hash)

		if exist {
			task.Status = model.TaskDuplicated
			c.End(msg)
			return errors.New(msg)
		}

	}

	return nil
}

func checkByHash(hash string) bool {

	//Check local hash first
	exist, _ := filter.CheckThenAdd(config.ContentHashFilter, []byte(hash))
	if exist {
		return true
	}

	//Check hash from db
	items, err := model.GetSnapshotByField("hash", hash)

	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		return true
	}

	return false
}
