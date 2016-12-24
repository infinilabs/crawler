package tasks

import "github.com/medcl/gopa/core/types"
import (
	log "github.com/cihub/seelog"
	"time"
	"github.com/rs/xid"
	"github.com/medcl/gopa/core/store"
)


func CreateTask(task *types.Task)  {
	log.Trace("start create crawler task")
	time:=time.Now()
	task.ID=xid.New().String()
	task.Status=types.TaskCreated
	task.CreateTime=&time
	task.UpdateTime=&time
	err := store.Save(task)
	if(err!=nil){
		panic(err)
	}
}

func LoadTaskByID(id string)(types.Task)   {
	task:=types.Task{}
	log.Trace("get id,",id)
	err := store.Get("ID",id,&task)
	if(err!=nil){
		panic(err)
	}
	return task
}

func UpdateTask(task *types.Task)  {
	log.Trace("start update crawler task")
	time:=time.Now()
	task.UpdateTime=&time
	err := store.Update(task)
	if(err!=nil){
		panic(err)
	}
}

func DeleteTask(id string)error  {
	log.Trace("start delete crawler task: ",id )
	task:=types.Task{ID:id}
	err := store.Delete(&task)
	if(err!=nil){
		panic(err)
	}
	return err
}

func GetTask(id int) (types.Task,error)  {
	log.Trace("start get seed: ",id)
	task:=types.Task{}
	err := store.Get("ID", id, &task)
	if(err!=nil){
		log.Error(id,", ",err)
	}
	return task,err
}

func GetTaskList(from,size int,skipDate string)(int,[]types.Task,error) {
	log.Trace("start get all crawler tasks")
	var tasks []types.Task
	queryO:=store.Query{Sort:"CreateTime",From:from,Size:size}
	err,result:=store.Search(&types.Task{},&tasks,&queryO)
	if(err!=nil){
		log.Debug(err)
	}
	return result.Total,tasks,err
}

func GetPendingFetchTasks()(int,[]types.Task,error) {
	log.Trace("start get all crawler tasks")
	var tasks []types.Task
	queryO:=store.Query{Sort:"CreateTime",Filter:&store.Cond{Name:"Status",Value:0}}
	err,result:=store.Search(&types.Task{},&tasks,&queryO)
	if(err!=nil){
		log.Error(err)
	}
	return result.Total,tasks,err
}
