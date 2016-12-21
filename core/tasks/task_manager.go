package tasks

import "github.com/medcl/gopa/core/types"
import (
	"github.com/asdine/storm"
	log "github.com/cihub/seelog"
	"time"
	"github.com/medcl/gopa/core/global"
	"path"
	"github.com/rs/xid"
	"github.com/asdine/storm/q"
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

func CreateSeed(task types.TaskSeed)  {
	if(!inited){Start()}
	log.Trace("start create seed")
	time:=time.Now()
	task.CreateTime=&time
	err := db.Save(&task)
	if(err!=nil){
		panic(err)
	}
	global.Env().Channels.PushUrlToCheck(task)
	log.Trace("end create seed")
}

func DeleteSeed(id int)  {
	if(!inited){Start()}
	log.Trace("start delete seed: ",id )
	task:=types.TaskSeed{ID:id}
	err := db.DeleteStruct(&task)
	if(err!=nil){
		panic(err)
	}
	log.Trace("end delete seed")
}

func GetSeed(id int) (types.TaskSeed,error)  {
	if(!inited){Start()}
	log.Trace("start get seed: ",id)
	task:=types.TaskSeed{}
	err := db.One("ID", id, &task)
	log.Trace("end get seed: ",id)
	return task,err
}

func GetSeedList()[]types.TaskSeed {
	if(!inited){Start()}
	log.Trace("start get all seeds")
	var tasks []types.TaskSeed
	err := db.AllByIndex("CreateTime",&tasks, storm.Reverse())
	if(err!=nil){
		panic(err)
	}
	log.Trace("end get all seeds")
	return tasks
}

func CreateTask(task *types.CrawlerTask)  {
	if(!inited){Start()}
	log.Trace("start create crawler task")
	time:=time.Now()
	task.ID=xid.New().String()
	task.CreateTime=&time
	err := db.Save(task)
	if(err!=nil){
		panic(err)
	}
	log.Trace("end create crawler task")
}

func UpdateTask(task *types.CrawlerTask)  {
	if(!inited){Start()}
	log.Trace("start update crawler task")
	time:=time.Now()
	task.UpdateTime=&time
	err := db.Update(task)
	if(err!=nil){
		panic(err)
	}
	log.Trace("end update crawler task")
}

func DeleteTask(id string)error  {
	if(!inited){Start()}
	log.Trace("start delete crawler task: ",id )
	task:=types.CrawlerTask{ID:id}
	err := db.DeleteStruct(&task)
	log.Trace("end delete crawler task: ",id )
	return err
}

func GetTask(id int) (types.CrawlerTask,error)  {
	if(!inited){Start()}
	log.Trace("start get seed: ",id)
	task:=types.CrawlerTask{}
	err := db.One("ID", id, &task)
	log.Trace("end get seed: ",id)
	return task,err
}

func GetTaskList(from,size int,skipDate string)(int,[]types.CrawlerTask,error) {
	if(!inited){Start()}
	log.Trace("start get all crawler tasks")
	var tasks []types.CrawlerTask
	total,err:=db.Count(&types.CrawlerTask{})
	if(err!=nil){
		log.Error(err)
	}
	if(skipDate!=""){
		log.Error(skipDate)
		layout := "2016-12-20T22:38:53.456485578+08:00"
		t, _ := time.Parse(layout, skipDate)
		query := db.Select(q.Gt("CreateTime", t)).Limit(size).Skip(from).Reverse().OrderBy("CreateTime")
		err = query.Find(&tasks)
	}else{
		err = db.AllByIndex("CreateTime",&tasks,storm.Skip(from),storm.Limit(size), storm.Reverse())
	}
	log.Trace("end get all crawler tasks")
	return total,tasks,err
}


func Create(o interface{})  {
	if(!inited){Start()}
	err := db.Save(o)
	if(err!=nil){
		panic(err)
	}
}

func Update(o interface{})  {
	if(!inited){Start()}
	err := db.Update(o)
	if(err!=nil){
		panic(err)
	}
}

func Delete(o interface{})error  {
	if(!inited){Start()}
	err := db.DeleteStruct(o)
	return err
}