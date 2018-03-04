package model

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/core/util"
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
	ID string `gorm:"not null;unique;primary_key" json:"id" index:"id"`
	// the url may not cleaned, may miss the host part, need reference to provide the complete url information
	Url         string `storm:"index" json:"url,omitempty" gorm:"type:varchar(500)"`
	Reference   string `json:"reference_url,omitempty"`
	Depth       int    `storm:"index" json:"depth"`
	Breadth     int    `storm:"index" json:"breadth"`
	Host        string `gorm:"index" json:"host"`
	Schema      string `json:"schema,omitempty"`
	OriginalUrl string `json:"original_url,omitempty" gorm:"type:varchar(500)"`
	Status      int    `gorm:"index" json:"status"`
	Message     string `json:"message,omitempty" gorm:"type:varchar(500)"`
	Created     int64  `gorm:"index" json:"created,omitempty"`
	Updated     int64  `gorm:"index" json:"updated,omitempty"`
	LastFetch   int64  `gorm:"index" json:"last_fetch,omitempty"`
	LastCheck   int64  `gorm:"index" json:"last_check,omitempty"`
	NextCheck   int64  `gorm:"index" json:"next_check,omitempty"`

	SnapshotVersion  int    `json:"snapshot_version,omitempty"`
	SnapshotID       string `json:"snapshot_id,omitempty"`
	SnapshotHash     string `json:"snapshot_hash,omitempty"`
	SnapshotSimHash  string `json:"snapshot_simhash,omitempty"`
	SnapshotCreated  int64  `json:"snapshot_created,omitempty"`
	LastScreenshotID string `json:"last_screenshot_id,omitempty"`

	PipelineConfigID string      `json:"pipline_config_id,omitempty"`
	HostConfig       *HostConfig `json:"host_config,omitempty"`

	// transient properties
	Snapshots     []Snapshot `json:"-"`
	SnapshotCount int        `json:"-"`
}

const (
	CONTEXT_TASK_ID               ParaKey = "GOPA_TASK_ID"
	CONTEXT_TASK_URL              ParaKey = "GOPA_TASK_URL"
	CONTEXT_TASK_Reference        ParaKey = "GOPA_TASK_Reference"
	CONTEXT_TASK_Depth            ParaKey = "GOPA_TASK_Depth"
	CONTEXT_TASK_Breadth          ParaKey = "GOPA_TASK_Breadth"
	CONTEXT_TASK_Host             ParaKey = "GOPA_TASK_Host"
	CONTEXT_TASK_Schema           ParaKey = "GOPA_TASK_Schema"
	CONTEXT_TASK_OriginalUrl      ParaKey = "GOPA_TASK_OriginalUrl"
	CONTEXT_TASK_Status           ParaKey = "GOPA_TASK_Status"
	CONTEXT_TASK_Message          ParaKey = "GOPA_TASK_Message"
	CONTEXT_TASK_Created          ParaKey = "GOPA_TASK_Created"
	CONTEXT_TASK_Updated          ParaKey = "GOPA_TASK_Updated"
	CONTEXT_TASK_LastFetch        ParaKey = "GOPA_TASK_LastFetch"
	CONTEXT_TASK_LastCheck        ParaKey = "GOPA_TASK_LastCheck"
	CONTEXT_TASK_NextCheck        ParaKey = "GOPA_TASK_NextCheck"
	CONTEXT_TASK_SnapshotID       ParaKey = "GOPA_TASK_SnapshotID"
	CONTEXT_TASK_SnapshotSimHash  ParaKey = "GOPA_TASK_SnapshotSimHash"
	CONTEXT_TASK_SnapshotHash     ParaKey = "GOPA_TASK_SnapshotHash"
	CONTEXT_TASK_SnapshotCreated  ParaKey = "GOPA_TASK_SnapshotCreated"
	CONTEXT_TASK_SnapshotVersion  ParaKey = "GOPA_TASK_SnapshotVersion"
	CONTEXT_TASK_LastScreenshotID ParaKey = "GOPA_TASK_LastScreenshotID"
	CONTEXT_TASK_PipelineConfigID ParaKey = "GOPA_TASK_PipelineConfigID"
	CONTEXT_TASK_Cookies          ParaKey = "GOPA_TASK_Cookies"

	CONTEXT_SNAPSHOT_ContentType ParaKey = "GOPA_SNAPSHOT_ContentType"
)

func CreateTask(task *Task) error {
	log.Trace("start create task, ", task.Url)
	time := time.Now().UTC().Unix()
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
	time := time.Now().UTC().Unix()
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

	if len(task.ID) == 0 || task.Updated == 0 {
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
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.DESC})
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
	log.Tracef("start get pending fetch tasks,last offset: %s,", offset)
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
	t := time.Now().UTC().Unix()
	log.Tracef("start get all tasks,last offset: %s,", offset)
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
