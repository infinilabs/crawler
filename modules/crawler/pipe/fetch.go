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
	"errors"
	log "github.com/cihub/seelog"
	. "github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/core/stats"
	"github.com/medcl/gopa/core/types"
	"github.com/medcl/gopa/core/util"
	"time"
)

type FetchJoint struct {
	context *Context
	timeout time.Duration
	cookie  string
}

func (this FetchJoint) Name() string {
	return "fetch"
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
	flg := make(chan bool, 1)

	go func() {
		pageItem := types.PageItem{}
		context.Set(CONTEXT_PAGE_LAST_FETCH, time.Now().UTC())

		//start to fetch remote content
		body, err := util.HttpGetWithCookie(&pageItem, requestUrl, this.cookie)

		if err == nil {
			if body != nil {
				if pageItem.StatusCode == 404 {
					log.Info("skip while 404, ", requestUrl, " , ", pageItem.StatusCode)
					context.Break()
					flg <- false
					return
				}
			}

			//update url, in case catch redirects
			//context.Set(CONTEXT_URL, pageItem.Url)//don't update, so what?
			context.Set(CONTEXT_PAGE_BODY_BYTES, body)
			context.Set(CONTEXT_PAGE_ITEM, &pageItem)
			log.Debug("exit fetchUrl method:", requestUrl)
			flg <- true

		} else {
			flg <- false
		}
	}()

	domain := context.MustGetString(CONTEXT_HOST)

	//监听通道，由于设有超时，不可能泄露
	select {
	case <-t.C:
		log.Error("fetching url time out, ", requestUrl)
		stats.Increment(domain, stats.STATS_FETCH_TIMEOUT_COUNT)
		context.Break()
		return nil, errors.New("fetch url time out")
	case value := <-flg:
		if value {
			log.Debug("fetching url normal exit, ", requestUrl)
			stats.Increment(domain, stats.STATS_FETCH_SUCCESS_COUNT)
		} else {
			log.Debug("fetching url error exit, ", requestUrl)
			context.Break()
			stats.Increment(domain, stats.STATS_FETCH_FAIL_COUNT)
		}
		return context, nil
	}

	return context, nil
}
