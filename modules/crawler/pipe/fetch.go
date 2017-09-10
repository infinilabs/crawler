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

package pipe

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/model"
	. "github.com/infinitbyte/gopa/core/pipeline"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	"net/http"
	"time"
)

const proxy ParaKey = "proxy"
const cookie ParaKey = "cookie"
const timeoutInSeconds ParaKey = "timeout_in_seconds"

type FetchJoint struct {
	Parameters
	timeout time.Duration
}

func (joint FetchJoint) Name() string {
	return "fetch"
}

type signal struct {
	flag   bool
	err    error
	status model.TaskStatus
}

func (joint FetchJoint) Process(context *Context) error {

	joint.timeout = time.Duration(joint.MustGetInt64(timeoutInSeconds)) * time.Second
	timer := time.NewTimer(joint.timeout)
	defer timer.Stop()

	task := context.MustGet(CONTEXT_CRAWLER_TASK).(*model.Task)
	snapshot := context.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	requestUrl := task.Url

	if len(requestUrl) == 0 {
		log.Error("invalid fetchUrl,", requestUrl)
		context.Exit("invalid fetch url")
		return errors.New("invalid fetchUrl")
	}

	t1 := time.Now().UTC()
	task.LastFetch = &t1

	log.Debug("start fetch url,", requestUrl)
	flg := make(chan signal, 1)
	go func() {

		cookie, _ := joint.GetString(cookie)
		proxy, _ := joint.GetString(proxy)

		//start to fetch remote content
		result, err := util.HttpGetWithCookie(requestUrl, cookie, proxy)

		if err == nil && result != nil {

			task.Url = result.Url //update url, in case catch redirects
			task.Host = result.Host

			snapshot.Payload = result.Body
			snapshot.StatusCode = result.StatusCode
			snapshot.Size = result.Size
			snapshot.Headers = result.Headers

			if result.Body != nil {

				if snapshot.StatusCode == 404 {
					log.Debug("skip while 404, ", requestUrl, " , ", snapshot.StatusCode)
					context.End("fetch 404")
					flg <- signal{flag: false, err: errors.New("404 NOT FOUND"), status: model.Task404}
					return
				}

				//detect content-type
				if snapshot.Headers != nil && len(snapshot.Headers) > 0 {
					v, ok := snapshot.Headers["content-type"]
					if ok {
						if len(v) > 0 {
							s := v[0]
							if s != "" {
								snapshot.ContentType = s
							} else {
								n := 512 // Only the first 512 bytes are used to sniff the content type.
								buffer := make([]byte, n)
								if len(snapshot.Payload) < n {
									n = len(snapshot.Payload)
								}
								// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
								contentType := http.DetectContentType(buffer[:n])
								snapshot.ContentType = contentType
							}
						}

					}

				}

			}
			log.Debug("exit fetchUrl method:", requestUrl)
			flg <- signal{flag: true, status: model.TaskSuccess}

		} else {

			code, payload := errors.CodeWithPayload(err)

			if code == errors.URLRedirected {
				if global.Env().IsDebug {
					log.Debug(util.ToJson(context, true))
				}
				task := model.NewTaskSeed(payload.(string), requestUrl, task.Depth, task.Breadth)
				log.Trace(err)
				queue.Push(config.CheckChannel, task.MustGetBytes())
				flg <- signal{flag: false, err: err, status: model.TaskRedirected}
				return
			}

			flg <- signal{flag: false, err: err, status: model.TaskFailed}
		}
	}()

	//check timeout
	select {
	case <-timer.C:
		log.Error("fetching url time out, ", requestUrl, ", ", joint.timeout)
		stats.Increment("domain.stats", task.Host+"."+config.STATS_FETCH_TIMEOUT_COUNT)
		task.Status = model.TaskTimeout
		context.End(fmt.Sprintf("fetching url time out, %s, %s", requestUrl, joint.timeout))
		return errors.New("fetch url time out")
	case value := <-flg:
		if value.flag {
			log.Debug("fetching url normal exit, ", requestUrl)
			stats.Increment("domain.stats", task.Host+"."+config.STATS_FETCH_SUCCESS_COUNT)
		} else {
			log.Debug("fetching url error exit, ", requestUrl)
			if value.err != nil {
				context.End(value.err.Error())
			}
			stats.Increment("domain.stats", task.Host+"."+config.STATS_FETCH_FAIL_COUNT)
		}
		task.Status = value.status
		return nil
	}
}
