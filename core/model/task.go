package model

import (
	"bytes"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/util"
	"strconv"
	"strings"
	"time"
)

type TaskStatus int

const TaskCreated TaskStatus = 0
const TaskFetchFailed TaskStatus = 2
const TaskFetchSuccess TaskStatus = 3
const Task404Ignore TaskStatus = 4
const TaskRedirectedIgnore TaskStatus = 5
const TaskFetchTimeout TaskStatus = 6

func GetTaskStatusText(status TaskStatus) string {
	switch status {
	case TaskCreated:
		return "created"
	case TaskFetchFailed:
		return "failed"
	case Task404Ignore:
		return "404"
	case TaskFetchSuccess:
		return "success"
	case TaskRedirectedIgnore:
		return "redirected"
	case TaskFetchTimeout:
		return "timeout"
	}
	return "unknow"
}

type Seed struct {
	// the seed url may not cleaned, may miss the domain part, need reference to provide the complete url information
	Url       string `storm:"index" json:"url,omitempty" gorm:"type:varchar(500)"`
	Reference string `json:"reference_url,omitempty"`
	Depth     int    `storm:"index" json:"depth"`
	Breadth   int    `storm:"index" json:"breadth"`
}

func (this Seed) Get(url string) Seed {
	task := Seed{}
	task.Url = url
	task.Reference = ""
	task.Depth = 0
	task.Breadth = 0
	return task
}

func (this Seed) MustGetBytes() []byte {

	bytes, err := this.GetBytes()
	if err != nil {
		panic(err)
	}
	return bytes
}

var delimiter = "|#|"

func (this Seed) GetBytes() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprint(this.Breadth))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprint(this.Depth))
	buf.WriteString(delimiter)
	buf.WriteString(this.Reference)
	buf.WriteString(delimiter)
	buf.WriteString(this.Url)

	return buf.Bytes(), nil
}

func TaskSeedFromBytes(b []byte) Seed {
	task, err := fromBytes(b)
	if err != nil {
		panic(err)
	}
	return task
}

func fromBytes(b []byte) (Seed, error) {

	str := string(b)
	array := strings.Split(str, delimiter)
	task := Seed{}
	i, _ := strconv.Atoi(array[0])
	task.Breadth = i
	i, _ = strconv.Atoi(array[1])
	task.Depth = i
	task.Reference = array[2]
	task.Url = array[3]

	return task, nil
}

func NewTaskSeed(url, ref string, depth int, breadth int) Seed {
	task := Seed{}
	task.Url = url
	task.Reference = ref
	task.Depth = depth
	task.Breadth = breadth
	return task
}

type Task struct {
	Seed
	ID          string          `gorm:"not null;unique;primary_key" json:"id" index:"id"`
	Host        string          `gorm:"index" json:"host"`
	Schema      string          `json:"schema,omitempty"`
	OriginalUrl string          `json:"original_url,omitempty"`
	Phrase      pipeline.Phrase `gorm:"index" json:"phrase"`
	Status      TaskStatus      `gorm:"index" json:"status"`
	Message     string          `json:"-"`
	Created     *time.Time      `gorm:"index" json:"created,omitempty"`
	Updated     *time.Time      `gorm:"index" json:"updated,omitempty"`
	LastFetch   *time.Time      `gorm:"index" json:"last_fetch"`
	LastCheck   *time.Time      `gorm:"index" json:"last_check"`
	NextCheck   *time.Time      `gorm:"index" json:"next_check"`

	SnapshotVersion int        `json:"snapshot_version"`
	SnapshotID      string     `json:"snapshot_id"`
	SnapshotHash    string     `json:"snapshot_hash"`
	SnapshotSimHash string     `json:"snapshot_simhash"`
	SnapshotCreated *time.Time `json:"snapshot_created"`
}

func CreateTask(task *Task) error {
	log.Trace("start create crawler task, ", task.Url)
	time := time.Now().UTC()
	task.ID = util.GetUUID()
	task.Status = TaskCreated
	task.Created = &time
	task.Updated = &time
	err := persist.Save(task)
	if err != nil {
		log.Error(task.ID, ", ", err)
	}
	return err
}

func UpdateTask(task *Task) {
	log.Trace("start update crawler task, ", task.Url)
	time := time.Now().UTC()
	task.Updated = &time
	err := persist.Update(task)
	if err != nil {
		panic(err)
	}
}

func DeleteTask(id string) error {
	log.Trace("start delete crawler task: ", id)
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
	if len(task.ID) == 0 || task.Created == nil {
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

func GetTaskList(from, size int, domain string) (int, []Task, error) {
	log.Tracef("start get crawler tasks, %v-%v, %v", from, size, domain)
	var tasks []Task
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.DESC})
	queryO := persist.Query{Sort: &sort, From: from, Size: size}
	if len(domain) > 0 {
		queryO.Conds = persist.And(persist.Eq("host", domain))
	}
	err, result := persist.Search(Task{}, &tasks, &queryO)
	if err != nil {
		log.Error(err)
		return 0, tasks, err
	}
	if result.Result != nil && tasks == nil || len(tasks) == 0 {
		convertTask(result, &tasks)
	}

	return result.Total, tasks, err
}

func GetPendingNewFetchTasks() (int, []Task, error) {
	log.Trace("start get all crawler tasks")
	var tasks []Task
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.DESC})
	queryO := persist.Query{Sort: &sort, Conds: persist.And(persist.Eq("phrase", 1))}
	err, result := persist.Search(Task{}, &tasks, &queryO)
	if err != nil {
		log.Error(err)
	}

	if result.Result != nil && tasks == nil || len(tasks) == 0 {
		convertTask(result, &tasks)
	}

	return result.Total, tasks, err
}

func GetPendingUpdateFetchTasks(offset *time.Time) (int, []Task, error) {
	t := time.Now().UTC()
	log.Tracef("start get all crawler tasks,last offset: %s,", offset.String())
	var tasks []Task
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.ASC})
	queryO := persist.Query{Sort: &sort,
		Conds: persist.And(
			persist.Lt("next_check", t),
			persist.Gt("created", offset),
			persist.Eq("status", TaskFetchSuccess)),
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
