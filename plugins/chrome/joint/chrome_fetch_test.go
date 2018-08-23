package joint

import (
	"context"
	"errors"
	"fmt"
	"github.com/infinitbyte/framework/core/util"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/mafredri/cdp/rpcc"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"
)

func ChromeFetchV2(t *testing.T) {
	//func TestChromeFetchV2(t *testing.T) {
	err := run(30 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

func run(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New("http://127.0.0.1:9223")
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
		return err
	}
	defer conn.Close() // Leaving connections open will leak memory.

	i := 611
	c := cdp.NewClient(conn)

	// Open a DOMContentEventFired client to buffer this event.
	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		return err
	}
	defer domContent.Close()

	// Enable events on the Page domain, it's often preferrable to create
	// event clients before enabling events so that we don't miss any.
	if err = c.Page.Enable(ctx); err != nil {
		return err
	}

	args := network.EnableArgs{}
	if err = c.Network.Enable(ctx, &args); err != nil {
		return err
	}

	// Enable console to evaluate scripts
	if err = c.Console.Enable(ctx); err != nil {
		return err
	}

	console, err := c.Console.MessageAdded(ctx)
	if err != nil {
		return err
	}

	go func() {
		defer console.Close()
		for {
			ev, err := console.Recv()
			if err != nil {
				return
			}
			fmt.Println("reply:", ev.Message.Text)
			if util.PrefixStr(ev.Message.Text, "gopa-") {
				array := strings.Split(ev.Message.Text, ":")
				fmt.Println(array[0], array[1])
				//c.Set(array[0],array[1])
			}
		}
	}()

	// Create the Navigate arguments with the optional Referrer field set.
	navArgs := page.NewNavigateArgs("http://localhost:8081/")
	//SetReferrer(fmt.Sprintf("https://duckduckgo.com?r=%v",i))
	nav, err := c.Page.Navigate(ctx, navArgs)
	if err != nil {
		return err
	}

	// Wait until we have a DOMContentEventFired event.
	if _, err = domContent.Recv(); err != nil {
		return err
	}

	fmt.Printf("Page loaded with frame ID: %s\n", nav.FrameID)

	// Fetch the document root node. We can pass nil here
	// since this method only takes optional arguments.
	doc, err := c.DOM.GetDocument(ctx, nil)
	if err != nil {
		return err
	}

	//Get the outer HTML for the page.
	result, err := c.DOM.GetOuterHTML(ctx, &dom.GetOuterHTMLArgs{
		NodeID: &doc.Root.NodeID,
	})
	if err != nil {
		return err
	}

	fmt.Println(result.OuterHTML)

	if strings.TrimSpace(result.OuterHTML) == "" || result.OuterHTML == "<html><head></head><body></body></html>" {
		panic(errors.New("empty body"))
	}

	//args1:=network.GetResponseBodyArgs{RequestID:""}
	//r1,_:=c.Network.GetResponseBody(ctx,args1)

	client, err := c.Network.ResponseReceived(ctx)
	go func() {
		select {
		case <-client.Ready():
			reply, err := client.Recv()
			if err != nil {
				return
			}
			fmt.Println(reply.Type)
			fmt.Println(reply.Response.Status)
			fmt.Println(reply.Response.StatusText)
			//fmt.Println(*reply.Response.HeadersText)
			v, _ := reply.Response.Headers.Map()
			fmt.Println(v)
			fmt.Println(reply.FrameID)
			fmt.Println(reply.RequestID)
			//return nil
		case <-time.After(3000 * time.Millisecond):
			//return errors.New("timeout")
		}
	}()
	//fmt.Printf("HTML: %s\n", result.OuterHTML)

	// Capture a screenshot of the current page.
	screenshotName := fmt.Sprintf("%v_screenshot.jpg", i)
	screenshotArgs := page.NewCaptureScreenshotArgs().
		SetFormat("jpeg").
		SetQuality(1)
	screenshot, err := c.Page.CaptureScreenshot(ctx, screenshotArgs)

	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(screenshotName, screenshot.Data, 0644); err != nil {
		return err
	}

	wait := true
	// get content-type
	args1 := runtime.EvaluateArgs{
		Expression:   "console.log('gopa-meta-content-type:'+document.contentType)",
		AwaitPromise: &wait}
	c.Runtime.Evaluate(ctx, &args1)

	fmt.Printf("Saved screenshot: %s\n", screenshotName)

	return nil
}
