package model

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/errors"
	"github.com/infinitbyte/framework/core/persist"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"time"
)

const TaskCreated int = 0
const TaskFailed int = 2
const TaskSuccess int = 3
const Task404 int = 4
const TaskRedirected int = 5
const TaskTimeout int = 6
const TaskDuplicated int = 7
const TaskInterrupted int = 8
const TaskPendingFetch int = 9

//  Common pipeline context keys
const (
	CONTEXT_SNAPSHOT   pipeline.ParaKey = "SNAPSHOT"
	CONTEXT_PAGE_LINKS pipeline.ParaKey = "PAGE_LINKS"
)

func GetTaskStatusText(status int) string {
	switch status {
	case TaskCreated:
		return "created"
	case TaskFailed:
		return "failed"
	case Task404:
		return "404"
	case TaskSuccess:
		return "success"
	case TaskRedirected:
		return "redirected"
	case TaskTimeout:
		return "timeout"
	case TaskDuplicated:
		return "duplicated"
	case TaskInterrupted:
		return "interrupted"
	case TaskPendingFetch:
		return "pending_fetch"
	}
	return "unknow"
}

func NewTask(url, ref string, depth int, breadth int) *Task {
	task := Task{}
	task.Url = url
	task.Reference = ref
	task.Depth = depth
	task.Breadth = breadth
	return &task
}

type Task struct {
	ID string `json:"id" index:"id"`
	// the url may not cleaned, may miss the host part, need reference to provide the complete url information
	Url         string    `storm:"index" json:"url,omitempty"`
	Reference   string    `json:"reference_url,omitempty"`
	Depth       int       `storm:"index" json:"depth"`
	Breadth     int       `storm:"index" json:"breadth"`
	Host        string    `json:"host"`
	Schema      string    `json:"schema,omitempty"`
	OriginalUrl string    `json:"original_url,omitempty"`
	Status      int       `json:"status"`
	Message     string    `json:"message,omitempty"`
	Created     time.Time `json:"created,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
	LastFetch   time.Time `json:"last_fetch,omitempty"`
	LastCheck   time.Time `json:"last_check,omitempty"`
	NextCheck   time.Time `json:"next_check,omitempty"`

	SnapshotVersion  int       `json:"snapshot_version,omitempty"`
	SnapshotID       string    `json:"snapshot_id,omitempty"`
	SnapshotHash     string    `json:"snapshot_hash,omitempty"`
	SnapshotSimHash  string    `json:"snapshot_simhash,omitempty"`
	SnapshotCreated  time.Time `json:"snapshot_created,omitempty"`
	LastScreenshotID string    `json:"last_screenshot_id,omitempty"`

	PipelineConfigID string      `json:"pipline_config_id,omitempty"`
	HostConfig       *HostConfig `json:"host_config,omitempty"`

	// transient properties
	Snapshots     []Snapshot `json:"-"`
	SnapshotCount int        `json:"-"`
}

const (
	CONTEXT_TASK_ID               pipeline.ParaKey = "GOPA_TASK_ID"
	CONTEXT_TASK_URL              pipeline.ParaKey = "GOPA_TASK_URL"
	CONTEXT_TASK_Reference        pipeline.ParaKey = "GOPA_TASK_Reference"
	CONTEXT_TASK_Depth            pipeline.ParaKey = "GOPA_TASK_Depth"
	CONTEXT_TASK_Breadth          pipeline.ParaKey = "GOPA_TASK_Breadth"
	CONTEXT_TASK_Host             pipeline.ParaKey = "GOPA_TASK_Host"
	CONTEXT_TASK_Schema           pipeline.ParaKey = "GOPA_TASK_Schema"
	CONTEXT_TASK_OriginalUrl      pipeline.ParaKey = "GOPA_TASK_OriginalUrl"
	CONTEXT_TASK_Status           pipeline.ParaKey = "GOPA_TASK_Status"
	CONTEXT_TASK_Message          pipeline.ParaKey = "GOPA_TASK_Message"
	CONTEXT_TASK_Created          pipeline.ParaKey = "GOPA_TASK_Created"
	CONTEXT_TASK_Updated          pipeline.ParaKey = "GOPA_TASK_Updated"
	CONTEXT_TASK_LastFetch        pipeline.ParaKey = "GOPA_TASK_LastFetch"
	CONTEXT_TASK_LastCheck        pipeline.ParaKey = "GOPA_TASK_LastCheck"
	CONTEXT_TASK_NextCheck        pipeline.ParaKey = "GOPA_TASK_NextCheck"
	CONTEXT_TASK_SnapshotID       pipeline.ParaKey = "GOPA_TASK_SnapshotID"
	CONTEXT_TASK_SnapshotSimHash  pipeline.ParaKey = "GOPA_TASK_SnapshotSimHash"
	CONTEXT_TASK_SnapshotHash     pipeline.ParaKey = "GOPA_TASK_SnapshotHash"
	CONTEXT_TASK_SnapshotCreated  pipeline.ParaKey = "GOPA_TASK_SnapshotCreated"
	CONTEXT_TASK_SnapshotVersion  pipeline.ParaKey = "GOPA_TASK_SnapshotVersion"
	CONTEXT_TASK_LastScreenshotID pipeline.ParaKey = "GOPA_TASK_LastScreenshotID"
	CONTEXT_TASK_PipelineConfigID pipeline.ParaKey = "GOPA_TASK_PipelineConfigID"
	CONTEXT_TASK_Cookies          pipeline.ParaKey = "GOPA_TASK_Cookies"

	CONTEXT_SNAPSHOT_ContentType pipeline.ParaKey = "GOPA_SNAPSHOT_ContentType"
)

func CreateTask(task *Task) error {
	log.Trace("start create task, ", task.Url)
	time := time.Now().UTC()
	task.ID = util.GetUUID()
	if task.OriginalUrl == "" {
		task.OriginalUrl = task.Url
	}
	task.Status = TaskCreated
	task.Created = time
	task.Updated = time
	if task.Url == "" {
		return errors.New("url can't be nil")
	}
	err := persist.Save(task)
	if err != nil {
		log.Error(task, ", ", err)
	}
	return err
}

func UpdateTask(task *Task) error {
	log.Trace("start update task, ", task.Url)
	time := time.Now().UTC()
	task.Updated = time
	if task.Url == "" {
		return errors.New("url can't be nil")
	}
	return persist.Update(task)
}

func DeleteTask(id string) error {
	log.Trace("start delete task: ", id)
	task := Task{ID: id}
	err := persist.Delete(&task)
	if err != nil {
		log.Error(id, ", ", err)
	}
	return err
}

func GetTask(id string) (Task, error) {
	log.Trace("start get seed: ", id)
	task := Task{}
	task.ID = id
	err := persist.Get(&task)
	if err != nil {
		log.Error(id, ", ", err)
	}

	if len(task.ID) == 0 || task.Updated.IsZero() {
		err = errors.New("not found," + id)
	}

	return task, err
}

func GetTaskByField(k, v string) ([]Task, error) {
	log.Trace("start get seed: ", k, ", ", v)
	task := Task{}
	tasks := []Task{}
	err, result := persist.GetBy(k, v, task, &tasks)

	if err != nil {
		log.Error(k, ", ", err)
		return tasks, err
	}
	if result.Result != nil && tasks == nil || len(tasks) == 0 {
		convertTask(result, &tasks)
	}

	return tasks, err
}

func GetTaskStatus(host string) (error, map[string]interface{}) {
	if host != "" {
		return persist.GroupBy(Task{}, "status", "host,status", "host = ?", host)

	} else {
		return persist.GroupBy(Task{}, "status", "status", "", nil)
	}
}

func GetHostStatus(status int) (error, map[string]interface{}) {
	if status >= 0 {
		return persist.GroupBy(Task{}, "host", "host,status", "status = ?", status)
	} else {
		return persist.GroupBy(Task{}, "host", "host", "", nil)
	}
}

func GetTaskList(from, size int, host string, status int) (int, []Task, error) {

	log.Tracef("start get tasks, %v-%v, %v", from, size, host)
	var tasks []Task
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.ASC})
	queryO := persist.Query{Sort: &sort, From: from, Size: size}
	if len(host) > 0 {
		queryO.Conds = persist.And(persist.Eq("host", host))
	}

	if status >= 0 {
		queryO.Conds = persist.Combine(queryO.Conds, persist.And(persist.Eq("status", status)))
	}

	err, result := persist.Search(Task{}, &tasks, &queryO)
	if err != nil {
		log.Error(err)
		return 0, tasks, err
	}
	if result.Result != nil && tasks == nil || len(tasks) == 0 {
		convertTask(result, &tasks)
	}
	log.Tracef("get %v tasks,total: %v", len(tasks), result.Total)
	return result.Total, tasks, err
}

func GetPendingNewFetchTasks(offset int64, size int) (int, []Task, error) {
	log.Tracef("start get pending fetch tasks,last offset: %v,", offset)
	var tasks []Task
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.ASC})
	queryO := persist.Query{Sort: &sort, Conds: persist.And(
		persist.Eq("status", TaskCreated),
		persist.Gt("created", offset)),
		From: 0, Size: size}
	err, result := persist.Search(Task{}, &tasks, &queryO)
	if err != nil {
		log.Error(err)
	}

	if result.Result != nil && tasks == nil || len(tasks) == 0 {
		convertTask(result, &tasks)
	}

	return result.Total, tasks, err
}

func GetFailedTasks(offset int64) (int, []Task, error) {
	log.Trace("start get all failed tasks")
	var tasks []Task
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.ASC})
	queryO := persist.Query{Sort: &sort, Conds: persist.And(
		persist.Gt("created", offset),
		persist.Eq("status", TaskFailed)),
		From: 0, Size: 100,
	}
	err, result := persist.Search(Task{}, &tasks, &queryO)
	if err != nil {
		log.Error(err)
	}

	if result.Result != nil && tasks == nil || len(tasks) == 0 {
		convertTask(result, &tasks)
	}

	return result.Total, tasks, err
}

func GetPendingUpdateFetchTasks(offset int64) (int, []Task, error) {
	t := time.Now().UTC()
	log.Tracef("start get all tasks,last offset: %v,", offset)
	var tasks []Task
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.ASC})
	queryO := persist.Query{Sort: &sort,
		Conds: persist.And(
			persist.Lt("next_check", t),
			persist.Gt("created", offset),
			persist.Eq("status", TaskSuccess)),
		From: 0, Size: 100,
	}
	err, result := persist.Search(Task{}, &tasks, &queryO)
	if err != nil {
		log.Error(err)
	}

	if result.Result != nil && tasks == nil || len(tasks) == 0 {
		convertTask(result, &tasks)
	}

	return result.Total, tasks, err
}

func convertTask(result persist.Result, tasks *[]Task) {
	if result.Result == nil {
		return
	}

	t, ok := result.Result.([]interface{})
	if ok {
		for _, i := range t {
			js := util.ToJson(i, false)
			t := Task{}
			util.FromJson(js, &t)
			*tasks = append(*tasks, t)
		}
	}
}
