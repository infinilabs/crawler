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
	Url         string    `storm:"index" json:"url,omitempty" gorm:"type:varchar(500)"`
	Reference   string    `json:"reference_url,omitempty"`
	Depth       int       `storm:"index" json:"depth"`
	Breadth     int       `storm:"index" json:"breadth"`
	Host        string    `gorm:"index" json:"host"`
	Schema      string    `json:"schema,omitempty"`
	OriginalUrl string    `json:"original_url,omitempty"`
	Status      int       `gorm:"index" json:"status"`
	Message     string    `json:"message,omitempty"`
	Created     time.Time `gorm:"index" json:"created,omitempty"`
	Updated     time.Time `gorm:"index" json:"updated,omitempty"`
	LastFetch   time.Time `gorm:"index" json:"last_fetch,omitempty"`
	LastCheck   time.Time `gorm:"index" json:"last_check,omitempty"`
	NextCheck   time.Time `gorm:"index" json:"next_check,omitempty"`

	SnapshotVersion int       `json:"snapshot_version,omitempty"`
	SnapshotID      string    `json:"snapshot_id,omitempty"`
	SnapshotHash    string    `json:"snapshot_hash,omitempty"`
	SnapshotSimHash string    `json:"snapshot_simhash,omitempty"`
	SnapshotCreated time.Time `json:"snapshot_created,omitempty"`

	PipelineConfigID string `json:"pipline_config_id,omitempty"`
}

const (
	CONTEXT_TASK_ID               ParaKey = "TASK_ID"
	CONTEXT_TASK_URL              ParaKey = "TASK_URL"
	CONTEXT_TASK_Reference        ParaKey = "TASK_Reference"
	CONTEXT_TASK_Depth            ParaKey = "TASK_Depth"
	CONTEXT_TASK_Breadth          ParaKey = "TASK_Breadth"
	CONTEXT_TASK_Host             ParaKey = "TASK_Host"
	CONTEXT_TASK_Schema           ParaKey = "TASK_Schema"
	CONTEXT_TASK_OriginalUrl      ParaKey = "TASK_OriginalUrl"
	CONTEXT_TASK_Status           ParaKey = "TASK_Status"
	CONTEXT_TASK_Message          ParaKey = "TASK_Message"
	CONTEXT_TASK_Created          ParaKey = "TASK_Created"
	CONTEXT_TASK_Updated          ParaKey = "TASK_Updated"
	CONTEXT_TASK_LastFetch        ParaKey = "TASK_LastFetch"
	CONTEXT_TASK_LastCheck        ParaKey = "TASK_LastCheck"
	CONTEXT_TASK_NextCheck        ParaKey = "TASK_NextCheck"
	CONTEXT_TASK_SnapshotID       ParaKey = "TASK_SnapshotID"
	CONTEXT_TASK_SnapshotSimHash  ParaKey = "TASK_SnapshotSimHash"
	CONTEXT_TASK_SnapshotHash     ParaKey = "TASK_SnapshotHash"
	CONTEXT_TASK_SnapshotCreated  ParaKey = "TASK_SnapshotCreated"
	CONTEXT_TASK_SnapshotVersion  ParaKey = "TASK_SnapshotVersion"
	CONTEXT_TASK_PipelineConfigID ParaKey = "TASK_PipelineConfigID"
)

func CreateTask(task *Task) error {
	log.Trace("start create task, ", task.Url)
	time := time.Now().UTC()
	task.ID = util.GetUUID()
	task.Status = TaskCreated
	task.Created = time
	task.Updated = time
	err := persist.Save(task)
	if err != nil {
		log.Error(task.ID, ", ", err)
	} else {
		IncrementHostLinkCount(task.Host)
	}
	return err
}

func UpdateTask(task *Task) {
	log.Trace("start update task, ", task.Url)
	time := time.Now().UTC()
	task.Updated = time
	err := persist.Update(task)
	if err != nil {
		panic(err)
	}
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
	if len(task.ID) == 0 || task.Created.IsZero() {
		panic(errors.New("not found," + id))
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

func GetTaskList(from, size int, host string) (int, []Task, error) {
	log.Tracef("start get tasks, %v-%v, %v", from, size, host)
	var tasks []Task
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.DESC})
	queryO := persist.Query{Sort: &sort, From: from, Size: size}
	if len(host) > 0 {
		queryO.Conds = persist.And(persist.Eq("host", host))
	}
	err, result := persist.Search(Task{}, &tasks, &queryO)
	if err != nil {
		log.Error(err)
		return 0, tasks, err
	}
	if result.Result != nil && tasks == nil || len(tasks) == 0 {
		convertTask(result, &tasks)
	}
	log.Tracef("get %v tasks", result.Total)
	return result.Total, tasks, err
}

func GetPendingNewFetchTasks() (int, []Task, error) {
	log.Trace("start get all tasks")
	var tasks []Task
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.DESC})
	queryO := persist.Query{Sort: &sort, Conds: persist.And(persist.Eq("status", TaskCreated))}
	err, result := persist.Search(Task{}, &tasks, &queryO)
	if err != nil {
		log.Error(err)
	}

	if result.Result != nil && tasks == nil || len(tasks) == 0 {
		convertTask(result, &tasks)
	}

	return result.Total, tasks, err
}

func GetFailedTasks(offset time.Time) (int, []Task, error) {
	log.Trace("start get all failed tasks")
	var tasks []Task
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.ASC})
	queryO := persist.Query{Sort: &sort, Conds: persist.And(
		persist.Gt("created", offset),
		persist.Eq("status", TaskFailed)),
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

func GetPendingUpdateFetchTasks(offset time.Time) (int, []Task, error) {
	t := time.Now().UTC()
	log.Tracef("start get all tasks,last offset: %s,", offset.String())
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
