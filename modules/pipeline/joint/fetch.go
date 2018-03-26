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
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/queue"
	"github.com/infinitbyte/gopa/core/stats"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/infinitbyte/gopa/modules/config"
	"net/http"
	"strings"
	"time"
)

const proxy model.ParaKey = "proxy"
const cookie model.ParaKey = "cookie"
const timeoutInSeconds model.ParaKey = "timeout_in_seconds"

type FetchJoint struct {
	model.Parameters
	timeout time.Duration
}

func (joint FetchJoint) Name() string {
	return "fetch"
}

type signal struct {
	flag   bool
	err    error
	status int
}

func (joint FetchJoint) Process(context *model.Context) error {

	joint.timeout = time.Duration(joint.GetInt64OrDefault(timeoutInSeconds, 60)) * time.Second
	timer := time.NewTimer(joint.timeout)
	defer timer.Stop()

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	requestUrl := context.MustGetString(model.CONTEXT_TASK_URL)

	if len(requestUrl) == 0 {
		log.Error("invalid fetchUrl,", requestUrl)
		context.End("invalid fetch url")
		return errors.New("invalid fetchUrl")
	}

	t1 := time.Now().UTC()
	context.Set(model.CONTEXT_TASK_LastFetch, t1)

	log.Debug("start fetch url,", requestUrl)
	flg := make(chan signal, 1)
	go func() {

		cookie, _ := joint.GetString(cookie)
		proxy, _ := joint.GetString(proxy)

		//start to fetch remote content
		result, err := util.HttpGetWithCookie(requestUrl, cookie, proxy)

		if err == nil && result != nil {

			//update url, in case catch redirects
			context.Set(model.CONTEXT_TASK_URL, result.Url)
			context.Set(model.CONTEXT_TASK_Host, result.Host)

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
								// only use the first 512 bytes to detect the content type.
								n := 512
								buffer := make([]byte, n)
								if len(snapshot.Payload) < n {
									n = len(snapshot.Payload)
								}
								// always returns a valid content-type and "application/octet-stream" if no others seemed to match.
								contentType := http.DetectContentType(buffer[:n])
								snapshot.ContentType = contentType
							}
							//normalize content-type
							snapshot.ContentType = strings.ToLower(snapshot.ContentType)
							snapshot.ContentType = strings.TrimSpace(snapshot.ContentType)
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

				depth := context.GetOrDefault(model.CONTEXT_TASK_Depth, 0)
				breadth := context.GetOrDefault(model.CONTEXT_TASK_Breadth, 0)

				context := model.Context{IgnoreBroken: true}
				context.Set(model.CONTEXT_TASK_URL, payload.(string))
				context.Set(model.CONTEXT_TASK_Reference, requestUrl)
				context.Set(model.CONTEXT_TASK_Depth, depth)
				context.Set(model.CONTEXT_TASK_Breadth, breadth)
				err = queue.Push(config.CheckChannel, util.ToJSONBytes(context))
				if err != nil {
					log.Error(err)
				}
				flg <- signal{flag: false, err: err, status: model.TaskRedirected}
				return
			}

			flg <- signal{flag: false, err: err, status: model.TaskFailed}
		}
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
