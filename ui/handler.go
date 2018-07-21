package ui

import (
	"errors"
	"github.com/infinitbyte/framework/core/api"
	"github.com/infinitbyte/framework/core/api/router"
	"github.com/infinitbyte/framework/core/persist"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
	"github.com/infinitbyte/gopa/ui/tasks"
	"net/http"
)

type AdminUI struct {
	api.Handler
}

func (h AdminUI) TasksPageAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var task []model.Task
	var count1, count2 int
	var host = h.GetParameterOrDefault(r, "host", "")
	var from = h.GetIntOrDefault(r, "from", 0)
	var size = h.GetIntOrDefault(r, "size", 20)
	var status = h.GetIntOrDefault(r, "status", -1)
	count1, task, _ = model.GetTaskList(from, size, host, status)

	//err, hvs := model.GetHostStatus(status)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err, kvs := model.GetTaskStatus(host)
	//if err != nil {
	//	panic(err)
	//}

	tasks.Index(w, r, host, status, from, size, count1, task, count2, nil, nil)
}

func (h AdminUI) TaskViewPageAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if id == "" {
		panic(errors.New("id is nill"))
	}

	task, err := model.GetTask(id)
	if err != nil {
		panic(err)
	}

	total, snapshots, err := model.GetSnapshotList(0, 10, id)
	task.Snapshots = snapshots
	task.SnapshotCount = total

	tasks.View(w, r, task)
}

func (h AdminUI) GetScreenshotAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	bytes, err := persist.GetValue(config.ScreenshotBucketKey, []byte(id))
	if err != nil {
		h.Error(w, err)
		return
	}
	w.Write(bytes)
}
