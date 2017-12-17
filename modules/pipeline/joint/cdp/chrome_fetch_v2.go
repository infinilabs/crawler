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

package cdp

import (
	c "context"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/persist"
	"github.com/infinitbyte/gopa/core/util"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/mafredri/cdp/rpcc"
	"strings"
	"time"
)

const timeoutInSeconds model.ParaKey = "timeout_in_seconds"
const chromeHost model.ParaKey = "chrome_host"
const saveScreenshot model.ParaKey = "save_screenshot"
const screenshotQuality model.ParaKey = "screenshot_quality"
const screenshotFormat model.ParaKey = "screenshot_format"

const bucket model.ParaKey = "bucket"

type ChromeFetchV2Joint struct {
	model.Parameters
	timeout time.Duration
}

type signal struct {
	flag   bool
	err    error
	status int
}

func (joint ChromeFetchV2Joint) Name() string {
	return "chrome_fetch_v2"
}

func (joint ChromeFetchV2Joint) Process(context *model.Context) error {

	joint.timeout = time.Duration(joint.GetInt64OrDefault(timeoutInSeconds, 10)) * time.Second

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	requestUrl := context.MustGetString(model.CONTEXT_TASK_URL)
	reference := context.MustGetString(model.CONTEXT_TASK_Reference)

	if len(requestUrl) == 0 {
		log.Error("invalid fetchUrl,", requestUrl)
		context.End("invalid fetch url")
		return errors.New("invalid fetchUrl")
	}

	t1 := time.Now().UTC()
	context.Set(model.CONTEXT_TASK_LastFetch, t1)

	log.Debug("start chrome v2 fetch url,", requestUrl)

	ctx, cancel := c.WithTimeout(c.Background(), joint.timeout)
	defer cancel()

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New(joint.GetStringOrDefault(chromeHost, "http://127.0.0.1:9223"))
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			return err
		}
	}

	// Initiate a new RPC connection to the Chrome Debugging Protocol target.
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		context.End(err)
		return err
	}
	defer conn.Close() // Leaving connections open will leak memory.

	c := cdp.NewClient(conn)

	// Open a DOMContentEventFired client to buffer this event.
	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		context.End(err)
		return err
	}
	defer domContent.Close()

	// Enable events on the Page domain, it's often preferrable to create
	// event clients before enabling events so that we don't miss any.
	if err = c.Page.Enable(ctx); err != nil {
		context.End(err)
		return err
	}

	// Enable console to evaluate scripts
	if err = c.Console.Enable(ctx); err != nil {
		context.End(err)
		return err
	}

	console, err := c.Console.MessageAdded(ctx)
	if err != nil {
		context.End(err)
		return err
	}

	go func(c *model.Context) {
		defer console.Close()
		for {
			ev, err := console.Recv()
			if err != nil {
				return
			}
			txt := ev.Message.Text
			log.Trace(requestUrl, ", console message:", txt)
			if util.PrefixStr(txt, "GOPA_") {
				if util.ContainStr(txt, ":") {
					array := strings.Split(txt, ":")
					if array[0] == string(model.CONTEXT_SNAPSHOT_ContentType) {
						c.Set(model.CONTEXT_SNAPSHOT_ContentType, array[1])
					}
				}
			}
		}
	}(context)

	// Create the Navigate arguments with the optional Referrer field set.
	navArgs := page.NewNavigateArgs(requestUrl).SetReferrer(reference)

	nav, err := c.Page.Navigate(ctx, navArgs)
	if err != nil {
		context.End(err)
		return err
	}

	if nav.ErrorText != nil {
		log.Error(nav.ErrorText)
	}

	// Wait until we have a DOMContentEventFired event.
	if _, err = domContent.Recv(); err != nil {
		context.End(err)
		return err
	}

	// Get content-type
	wait := true
	args := runtime.EvaluateArgs{
		Expression:   fmt.Sprintf("console.log('%s:'+document.contentType)", model.CONTEXT_SNAPSHOT_ContentType),
		AwaitPromise: &wait}
	c.Runtime.Evaluate(ctx, &args)

	// Fetch the document root node. We can pass nil here
	// since this method only takes optional arguments.
	doc, err := c.DOM.GetDocument(ctx, nil)
	if err != nil {
		return err
	}

	// Get the outer HTML for the page.
	result, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &doc.Root.NodeID,
	})

	if err != nil {
		context.End(err)
		return err
	}

	if strings.TrimSpace(result.OuterHTML) == "" || result.OuterHTML == "<html><head></head><body></body></html>" {
		err := errors.New("the response is empty")
		panic(err)
	}

	if joint.GetBool(saveScreenshot, false) {
		// Capture a screenshot of the current page.
		screenshotArgs := page.NewCaptureScreenshotArgs().
			SetFormat(joint.GetStringOrDefault(screenshotFormat, "jpeg")).
			SetQuality(joint.GetIntOrDefault(screenshotQuality, 10))
		screenshot, err := c.Page.CaptureScreenshot(ctx, screenshotArgs)
		if err != nil {
			context.End(err)
			return err
		}

		bucketName := joint.GetStringOrDefault(bucket, "Screenshot")
		uuid := []byte(util.GetUUID())

		//for picture, no need to compress
		err = persist.AddValue(bucketName, uuid, screenshot.Data)
		if err != nil {
			context.End(err)
			return err
		}
		snapshot.ScreenshotID = string(uuid)
		context.Set(model.CONTEXT_TASK_LastScreenshotID, snapshot.ScreenshotID)
	}

	snapshot.Payload = []byte(result.OuterHTML)
	snapshot.Size = uint64(len(result.OuterHTML))

	//snapshot.StatusCode = reply.Response.Status
	if context.Has(model.CONTEXT_SNAPSHOT_ContentType) {
		snapshot.ContentType = context.GetStringOrDefault(model.CONTEXT_SNAPSHOT_ContentType, "N/A")
	} else {
		log.Error(requestUrl, ", failed to get content-type")
		snapshot.ContentType = "N/A"
	}

	log.Debug("exit chrome v2 fetch method:", requestUrl)

	return nil
}
