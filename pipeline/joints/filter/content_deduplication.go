package filter

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/errors"
	"github.com/infinitbyte/framework/core/filter"
	"github.com/infinitbyte/framework/core/global"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
)

// ContentDeduplicationJoint used to check the hash of page body, if duplicated hash already exists, will break the pipeline
type ContentDeduplicationJoint struct {
	pipeline.Parameters
}

// Name return: content_deduplication
func (joint ContentDeduplicationJoint) Name() string {
	return "content_deduplication"
}

// Process the content hash Deduplication
func (joint ContentDeduplicationJoint) Process(c *pipeline.Context) error {
	//task := c.MustGet(model.CONTEXT_TASK).(*model.Task)
	snapshot := c.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)
	url := c.MustGetString(model.CONTEXT_TASK_URL)
	if snapshot.Hash != "" {

		if global.Env().IsDebug {
			log.Trace("check content deduplication, ", url)
		}
		snapshot.Url = url
		taskID := c.MustGetString(model.CONTEXT_TASK_ID)
		snapshot.TaskID = taskID

		exist, depTaskID, depSnapshotId, depUrl := checkByHash(snapshot, c)

		msg := fmt.Sprintf("same content hash found, %s, %s, %s, duplicated with task: %s, snapshotID: %s, url: %s", taskID, url, snapshot.Hash, depTaskID, depSnapshotId, depUrl)

		if exist {
			c.Set(model.CONTEXT_TASK_Status, model.TaskDuplicated)
			c.End(msg)
			return errors.New(msg)
		}

	}

	return nil
}

func checkByHash(snapshot *model.Snapshot, c *pipeline.Context) (bool, string, string, string) {

	hash := snapshot.Hash

	if global.Env().IsDebug {

		log.Trace("current object hash:", hash)
	}
	//Check local hash first
	if c.GetBool("check_filter", false) {
		exist, _ := filter.CheckThenAdd(config.ContentHashFilter, []byte(hash))
		if exist {
			return true, "", "local_filter_cache", ""
		}
	}

	//Check hash from db
	items, err := model.GetSnapshotByField("hash", hash)

	if global.Env().IsDebug {
		log.Trace("get objects by hash:", items, ",", err)
	}

	if err != nil {
		panic(err)
	}

	if len(items) > 0 {
		for _, v := range items {
			if global.Env().IsDebug {
				log.Tracef("%s vs  %s , %s vs %s", v.Url, snapshot.Url, snapshot.TaskID, v.TaskID)
			}
			if v.Url != snapshot.Url && v.TaskID != snapshot.TaskID {
				return true, v.TaskID, v.ID, v.Url
			}
		}
	}

	return false, "", "", ""
}
