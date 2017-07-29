package model

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/store"

	"bytes"
	"fmt"
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
	Url       string `storm:"index" json:"url,omitempty" gorm:"type:not null;varchar(500)"` // the seed url may not cleaned, may miss the domain part, need reference to provide the complete url information
	Reference string `json:"reference_url,omitempty"`
	Depth     int    `storm:"index" json:"depth,omitempty"`
	Breadth   int    `storm:"index" json:"breadth,omitempty"`
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
	ID            string          `gorm:"not null;unique;primary_key" json:"id"`
	Host          string          `gorm:"index" json:"-"`
	Schema        string          `json:"schema,omitempty"`
	OriginalUrl   string          `json:"original_url,omitempty"`
	Phrase        pipeline.Phrase `gorm:"index" json:"phrase"`
	Status        TaskStatus      `gorm:"index" json:"status"`
	Message       string          `json:"-"`
	CreateTime    *time.Time      `gorm:"index" json:"created,omitempty"`
	UpdateTime    *time.Time      `gorm:"index" json:"updated,omitempty"`
	LastFetchTime *time.Time      `gorm:"index" json:"last_fetch"`
	LastCheckTime *time.Time      `gorm:"index" json:"last_check"`
	NextCheckTime *time.Time      `gorm:"index" json:"next_check"`

	SnapshotVersion    int        `json:"snapshot_version"`
	SnapshotID         string     `json:"snapshot_id"`      //Last Snapshot's ID
	SnapshotHash       string     `json:"snapshot_hash"`    //Last Snapshot's Hash
	SnapshotSimHash    string     `json:"snapshot_simhash"` //Last Snapshot's Simhash
	SnapshotCreateTime *time.Time `json:"snapshot_created"` //Last Snapshot's Simhash
}

func CreateTask(task *Task) error {
	log.Trace("start create crawler task, ", task.Url)
	time := time.Now().UTC()
	task.ID = util.GetIncrementID("task")
	task.Status = TaskCreated
	task.CreateTime = &time
	task.UpdateTime = &time
	err := store.Save(task)
	if err != nil {
		log.Error(task.ID, ", ", err)
	}
	return err
}

func UpdateTask(task *Task) {
	log.Trace("start update crawler task, ", task.Url)
	time := time.Now().UTC()
	task.UpdateTime = &time
	err := store.Update(task)
	if err != nil {
		panic(err)
	}
}

func DeleteTask(id string) error {
	log.Trace("start delete crawler task: ", id)
	task := Task{ID: id}
	err := store.Delete(&task)
	if err != nil {
		log.Error(id, ", ", err)
	}
	return err
}

func GetTask(id string) (Task, error) {
	log.Trace("start get seed: ", id)
	task := Task{}
	err := store.GetBy("id", id, &task)
	if err != nil {
		log.Error(id, ", ", err)
	}
	if len(task.ID) == 0 || task.CreateTime == nil {
		panic(errors.New("not found," + id))
	}

	return task, err
}

func GetTaskByField(k, v string) (Task, error) {
	log.Trace("start get seed: ", k, ", ", v)
	task := Task{}
	err := store.GetBy(k, v, &task)
	if err != nil {
		log.Error(k, ", ", err)
	}
	return task, err
}

func GetTaskList(from, size int, domain string) (int, []Task, error) {
	log.Tracef("start get crawler tasks, %v-%v, %v", from, size, domain)
	var tasks []Task
	queryO := store.Query{Sort: "create_time desc", From: from, Size: size}
	if len(domain) > 0 {
		queryO.Conds = store.And(store.Eq("host", domain))
	}
	err, result := store.Search(&tasks, &queryO)
	if err != nil {
		log.Error(err)
	}
	return result.Total, tasks, err
}

func GetPendingNewFetchTasks() (int, []Task, error) {
	log.Trace("start get all crawler tasks")
	var tasks []Task
	queryO := store.Query{Sort: "create_time desc", Conds: store.And(store.Eq("phrase", 1))}
	err, result := store.Search(&tasks, &queryO)
	if err != nil {
		log.Error(err)
	}
	return result.Total, tasks, err
}

func GetPendingUpdateFetchTasks(offset *time.Time) (int, []Task, error) {
	t := time.Now().UTC()
	log.Tracef("start get all crawler tasks,last offset: %s,", offset.String())
	var tasks []Task
	queryO := store.Query{Sort: "create_time asc",
		Conds: store.And(
			store.Lt("next_check_time", t),
			store.Gt("create_time", offset),
			store.Eq("status", TaskFetchSuccess))}
	err, result := store.Search(&tasks, &queryO)
	if err != nil {
		log.Error(err)
	}
	return result.Total, tasks, err
}
