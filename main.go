/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	_ "expvar"
	"github.com/infinitbyte/framework"
	"github.com/infinitbyte/framework/core/module"
	"github.com/infinitbyte/framework/core/orm"
	"github.com/infinitbyte/framework/modules"
	"github.com/infinitbyte/gopa/api"
	"github.com/infinitbyte/gopa/config"
	"github.com/infinitbyte/gopa/model"
	"github.com/infinitbyte/gopa/pipeline"
	"github.com/infinitbyte/gopa/plugins"
	"github.com/infinitbyte/gopa/ui"
)

func main() {

	terminalHeader := ("  ________ ________ __________  _____   \n")
	terminalHeader += (" /  _____/ \\_____  \\\\______   \\/  _  \\  \n")
	terminalHeader += ("/   \\  ___  /   |   \\|     ___/  /_\\  \\ \n")
	terminalHeader += ("\\    \\_\\  \\/    |    \\    |  /    |    \\\n")
	terminalHeader += (" \\______  /\\_______  /____|  \\____|__  /\n")
	terminalHeader += ("        \\/         \\/                \\/ \n")

	terminalFooter := ("                         |    |                \n")
	terminalFooter += ("   _` |   _ \\   _ \\   _` |     _ \\  |  |   -_) \n")
	terminalFooter += (" \\__, | \\___/ \\___/ \\__,_|   _.__/ \\_, | \\___| \n")
	terminalFooter += (" ____/                             ___/        \n")

	app := framework.NewApp("gopa", "A Spider Written in Go.",
		config.Version, config.LastCommitLog, config.BuildDate, terminalHeader, terminalFooter)

	app.Init(nil)
	defer app.Shutdown()

	app.Start(func() {
		//modules
		module.New()

		//load core modules first
		modules.Register()

		//register API
		api.InitAPI()

		//register UI
		ui.InitUI()

		//register joints
		pipeline.InitJoints()

		//load plugins
		plugins.Register()

		//start each module, with enabled provider
		module.Start()

	}, func() {

		orm.RegisterSchema(&model.Host{})
		orm.RegisterSchema(&model.Task{})
		orm.RegisterSchema(&model.Snapshot{})
		orm.RegisterSchema(&model.HostConfig{})
		orm.RegisterSchema(&model.Project{})
		orm.RegisterSchema(&model.Index{})
	})

}
