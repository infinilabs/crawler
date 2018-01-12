package admin

import (
	"errors"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/http"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/modules/config"
	"github.com/infinitbyte/gopa/modules/ui/admin/console"
	"github.com/infinitbyte/gopa/modules/ui/admin/dashboard"
	"github.com/infinitbyte/gopa/modules/ui/admin/explore"
	"github.com/infinitbyte/gopa/modules/ui/admin/setting"
	"github.com/infinitbyte/gopa/modules/ui/admin/tasks"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v2"
	"net/http"
)

type AdminUI struct {
	api.Handler
}

func (h AdminUI) DashboardAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	dashboard.Index(w, r)
}

func (h AdminUI) TasksPageAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var task []model.Task
	var count1, count2 int
	var host = h.GetParameterOrDefault(r, "host", "")
	var from = h.GetIntOrDefault(r, "from", 0)
	var size = h.GetIntOrDefault(r, "size", 20)
	var status = h.GetIntOrDefault(r, "status", -1)
	count1, task, _ = model.GetTaskList(from, size, host, status)

	var hosts []model.Host
	count2, hosts, _ = model.GetHostList(0, 35, "")
	tasks.Index(w, r, host, status, from, size, count1, task, count2, hosts)
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

func (h AdminUI) ConsolePageAction(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	console.Index(w, r)
}

func (h AdminUI) ExplorePageAction(w http.ResponseWriter, r *http.Request) {

	explore.Index(w, r)
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

func (h AdminUI) SettingPageAction(w http.ResponseWriter, r *http.Request) {

	o, _ := yaml.Marshal(global.Env().RuntimeConfig)
	setting.Setting(w, r, string(o))
}

func (h AdminUI) UpdateSettingAction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	body, _ := h.GetRawBody(r)
	yaml.Unmarshal(body, global.Env().RuntimeConfig) //TODO extract method, save to file

	o, _ := yaml.Marshal(global.Env().RuntimeConfig)

	setting.Setting(w, r, string(o))
}
