package model

import (
	log "github.com/cihub/seelog"
	"github.com/medcl/gopa/core/store"
	"github.com/rs/xid"
	"time"
	"github.com/medcl/gopa/core/errors"
)

func CreateTask(task *Task) error {
	log.Trace("start create crawler task, ",task.Url)
	time := time.Now()
	task.ID = xid.New().String()
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
	log.Trace("start update crawler task, ",task.Url)
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
	if(len(task.ID)==0||task.CreateTime==nil){
		panic(errors.New("not found,"+id))
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
	log.Trace("start get all crawler tasks")
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
