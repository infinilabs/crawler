package model

import (
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/errors"
	"github.com/medcl/gopa/core/store"

	"github.com/medcl/gopa/core/util"
	"time"
)

type TaskStatus int

const TaskCreated TaskStatus = 0
const TaskFetchFailed TaskStatus = 2
const TaskFetchSuccess TaskStatus = 3

type TaskPhrase int

type Task struct {
	Seed
	ID              string      `storm:"id,unique" json:"id" gorm:"not null;unique;primary_key"`
	Domain          string      `storm:"index" json:"domain,omitempty"` // elastic.co
	Scheme          string      `json:"schema,omitempty"`               // elastic.co
	OriginUrl       string      `json:"origin_url,omitempty"`           // /index.html
	UrlPath         string      `json:"path,omitempty"`                 // /index.html
	Phrase          TaskPhrase  `storm:"index" json:"phrase"`
	Status          TaskStatus  `storm:"index" json:"status"`
	Page            *PageItem   `storm:"inline" json:"page,omitempty" gorm:"-"`
	Message         interface{} `storm:"inline" json:"message,omitempty" gorm:"-"`
	CreateTime      *time.Time  `storm:"index" json:"created,omitempty" gorm:"index"`
	UpdateTime      *time.Time  `storm:"index" json:"updated,omitempty" gorm:"index"`
	LastCheckTime   *time.Time  `storm:"index" json:"checked,omitempty"`
	Snapshot        string      `json:"snapshot,omitempty"` //Last Snapshot ID
	SnapshotHash    string      `storm:"snapshot_hash" json:"snapshot_hash,omitempty"`
	SnapshotSimHash string      `storm:"snapshot_simhash" json:"snapshot_simhash,omitempty"`
}

func CreateTask(task *Task) error {
	log.Trace("start create crawler task, ", task.Url)
	time := time.Now()
	task.ID = util.GetIncrementID("task")
	task.Status = TaskCreated
	task.CreateTime = &time
	task.UpdateTime = &time
	err := store.Save(task)
	if err != nil {
		log.Debug(task.ID, ", ", err)
	}
	return err
}

func UpdateTask(task *Task) {
	log.Trace("start update crawler task, ", task.Url)
	time := time.Now()
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
		log.Debug(id, ", ", err)
	}
	return err
}

func GetTask(id string) (Task, error) {
	log.Trace("start get seed: ", id)
	task := Task{}
	err := store.GetBy("id", id, &task)
	if err != nil {
		log.Debug(id, ", ", err)
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
		log.Debug(k, ", ", err)
	}
	return task, err
}

func GetTaskList(from, size int, domain string) (int, []Task, error) {
	log.Tracef("start get crawler tasks, %v-%v, %v", from, size, domain)
	var tasks []Task
	queryO := store.Query{Sort: "create_time desc", From: from, Size: size}
	if len(domain) > 0 {
		queryO.Filter = &store.Cond{Name: "domain", Value: domain}
	}
	err, result := store.Search(&tasks, &queryO)
	if err != nil {
		log.Trace(err)
	}
	return result.Total, tasks, err
}

func GetPendingFetchTasks() (int, []Task, error) {
	log.Trace("start get all crawler tasks")
	var tasks []Task
	queryO := store.Query{Sort: "create_time desc", Filter: &store.Cond{Name: "phrase", Value: 1}}
	err, result := store.Search(&tasks, &queryO)
	if err != nil {
		log.Trace(err)
	}
	return result.Total, tasks, err
}
