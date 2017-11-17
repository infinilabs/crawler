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

package joint

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/modules/config"
	"io/ioutil"
	"os/exec"
	"time"
)

type ChromeFetchJoint struct {
	model.Parameters
	timeout time.Duration
}

func (joint ChromeFetchJoint) Name() string {
	return "chrome_fetch"
}

func (joint ChromeFetchJoint) Process(context *model.Context) error {

	joint.timeout = time.Duration(joint.GetInt64OrDefault(timeoutInSeconds, 60)) * time.Second
	timer := time.NewTimer(joint.timeout)
	defer timer.Stop()

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	requestUrl := context.MustGetString(model.CONTEXT_TASK_URL)

	if len(requestUrl) == 0 {
		log.Error("invalid fetchUrl,", requestUrl)
		context.Exit("invalid fetch url")
		return errors.New("invalid fetchUrl")
	}

	t1 := time.Now().UTC()
	context.Set(model.CONTEXT_TASK_LastFetch, t1)

	command := joint.GetStringOrDefault("command", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome")

	log.Debug("start fetch url,", requestUrl)
	flg := make(chan signal, 1)
	go func() {

		cmd := exec.Command(command, "--headless", "-disable-gpu", "--dump-dom", requestUrl)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		content, err := ioutil.ReadAll(stdout)
		if err != nil {
			panic(err)
		}

		if err != nil {
			flg <- signal{flag: false, err: err, status: model.TaskFailed}
			return
		}

		snapshot.Payload = content
		snapshot.StatusCode = 200
		snapshot.Size = uint64(len(content))
		snapshot.ContentType = "text/html"
		log.Debug("exit fetchUrl method:", requestUrl)
		flg <- signal{flag: true, status: model.TaskSuccess}
	}()

	//check timeout
	select {
	case <-timer.C:
		host := context.MustGetString(model.CONTEXT_TASK_Host)

		log.Error("fetching url time out, ", requestUrl, ", ", joint.timeout)
		stats.Increment("host.stats", host+"."+config.STATS_FETCH_TIMEOUT_COUNT)

		context.Set(model.CONTEXT_TASK_Status, model.TaskTimeout)

		context.End(fmt.Sprintf("fetching url time out, %s, %s", requestUrl, joint.timeout))
		return errors.New("fetch url time out")
	case value := <-flg:
		host := context.MustGetString(model.CONTEXT_TASK_Host)
		if value.flag {
			log.Debug("fetching url normal exit, ", requestUrl)
			stats.Increment("host.stats", host+"."+config.STATS_FETCH_SUCCESS_COUNT)
		} else {
			log.Debug("fetching url error exit, ", requestUrl)
			if value.err != nil {
				context.End(value.err.Error())
			}
			stats.Increment("host.stats", host+"."+config.STATS_FETCH_FAIL_COUNT)
		}
		context.Set(model.CONTEXT_TASK_Status, value.status)
		return nil
	}
}
