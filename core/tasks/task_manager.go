package tasks

import "github.com/medcl/gopa/core/types"
import (
	"github.com/asdine/storm"
	log "github.com/cihub/seelog"
	"time"
	"github.com/medcl/gopa/core/global"
	"path"
)


var db *storm.DB
var inited bool
func Start() error  {
	var err error
	file:= path.Join(global.Env().RuntimeConfig.PathConfig.Data,"task_db")
	db, err = storm.Open(file)
	inited=true
	return err

}

func Stop()  {
	db.Close()
}


func CreateTask(task types.PageTask)  {
	if(!inited){Start()}
	log.Trace("start create task")
	task.CreateTime=time.Now()
	err := db.Save(&task)
	if(err!=nil){
		panic(err)
	}
}

func DeleteTask(id int)  {
	if(!inited){Start()}
	log.Trace("start delete task: ",id )
	task:=types.PageTask{ID:id}
	err := db.DeleteStruct(&task)
	if(err!=nil){
		panic(err)
	}
}
func GetTask(id int) (types.PageTask,error)  {
	if(!inited){Start()}
	log.Trace("start get task: ",id)
	task:=types.PageTask{}
	err := db.One("ID", id, &task)
	return task,err
}

func GetTaskList()[]types.PageTask  {
	if(!inited){Start()}
	log.Trace("start get all tasks")
	var tasks []types.PageTask
	err := db.AllByIndex("CreateTime",&tasks)
	if(err!=nil){
		panic(err)
	}
	return tasks
}