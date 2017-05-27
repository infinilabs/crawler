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
	"github.com/medcl/gopa/core/errors"
	"github.com/medcl/gopa/core/model"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/queue"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/core/util"
	"github.com/medcl/gopa/modules/config"
	"time"
)

const Fetch JointKey = "fetch"
const Proxy ParaKey = "proxy"
const Cookie ParaKey = "cookie"

type FetchJoint struct {
	Parameters
	context *Context
	timeout time.Duration
}

func (this FetchJoint) Name() string {
	return string(Fetch)
}

type signal struct {
	flag bool
	err  error
}

func (this FetchJoint) Process(context *Context) (*Context, error) {

	this.timeout = 10 * time.Second
	this.context = context
	t := time.NewTimer(this.timeout)
	defer t.Stop()
	requestUrl := context.MustGetString(CONTEXT_URL)

	if len(requestUrl) == 0 {
		log.Error("invalid fetchUrl")
		return context, errors.New("invalid fetchUrl")
	}

	log.Debug("start fetch url,", requestUrl)
	flg := make(chan signal, 1)
	go func() {
		pageItem := model.PageItem{}
		context.Set(CONTEXT_PAGE_LAST_FETCH, time.Now().UTC())

		cookie, _ := this.GetString(Cookie)
		proxy, _ := this.GetString(Proxy) //"socks5://127.0.0.1:9150"  //TODO 这个是全局配置,后续的url应该也使用同样的配置,应该在domain setting里面
		//先全局,再domain,再task,再pipeline,层层覆盖
		log.Trace("proxy:", proxy)

		//start to fetch remote content
		body, err := util.HttpGetWithCookie(&pageItem, requestUrl, cookie, proxy)

		if err == nil {
			if body != nil {
				if pageItem.StatusCode == 404 {
					log.Info("skip while 404, ", requestUrl, " , ", pageItem.StatusCode)
					context.Break("fetch 404")
					flg <- signal{flag: false, err: errors.New("404 NOT FOUND")}
					return
				}
			}

			//update url, in case catch redirects
			context.Set(CONTEXT_URL, pageItem.Url) //should update here, if url got redirected, the url is change
			context.Set(CONTEXT_HOST, pageItem.Host)
			context.Set(CONTEXT_PAGE_BODY_BYTES, body)
			context.Set(CONTEXT_PAGE_ITEM, &pageItem)
			log.Debug("exit fetchUrl method:", requestUrl)
			flg <- signal{flag: true}

		} else {

			e, ok := err.(*errors.Error)
			if ok {
				if e.Code == errors.URLRedirected {
					log.Trace(util.ToJson(context, true))
					depth := context.MustGetInt(CONTEXT_DEPTH)
					breadth := context.MustGetInt(CONTEXT_BREADTH)
					task := model.NewTaskSeed(e.Payload.(string), requestUrl, depth, breadth)
					log.Trace(err)
					queue.Push(config.CheckChannel, task.MustGetBytes())
				}
			}
			flg <- signal{flag: false, err: err}
		}
	}()

	domain := context.MustGetString(CONTEXT_HOST)

	//监听通道，由于设有超时，不可能泄露
	select {
	case <-t.C:
		log.Error("fetching url time out, ", requestUrl, ", ", this.timeout)
		stats.Increment("domain.stats", domain+"."+stats.STATS_FETCH_TIMEOUT_COUNT)
		context.Break(fmt.Sprintf("fetching url time out, %s, %s", requestUrl, this.timeout))
		return nil, errors.New("fetch url time out")
	case value := <-flg:
		if value.flag {
			log.Debug("fetching url normal exit, ", requestUrl)
			stats.Increment("domain.stats", domain+"."+stats.STATS_FETCH_SUCCESS_COUNT)
		} else {
			log.Debug("fetching url error exit, ", requestUrl)
			if value.err != nil {
				context.Break(value.err.Error())
			}
			stats.Increment("domain.stats", domain+"."+stats.STATS_FETCH_FAIL_COUNT)
		}
		return context, nil
	}

	return context, nil
}
