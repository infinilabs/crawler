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
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	cdp "github.com/knq/chromedp"
	cdptypes "github.com/knq/chromedp/cdp"
	"github.com/knq/chromedp/cdp/network"
	"log"
	"net/http"
	"time"
)

var (
	flagPort = flag.Int("port", 8544, "port")
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := cdp.New(ctxt, cdp.WithErrorf(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	var buf string
	err = c.Run(ctxt, text(`http://git.infinitbyte.com/`, `body`, &buf))
	fmt.Println(buf)

	//// run task list
	//err = c.Run(ctxt, click())
	//if err != nil {
	//	log.Fatal(err)
	//}

	time.Sleep(1 * time.Hour)

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	//err = ioutil.WriteFile("contact-form.png", buf, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}

}

func main1() {
	var err error

	flag.Parse()

	// setup http server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		buf, err := json.MarshalIndent(req.Cookies(), "", "  ")
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(res, indexHTML, string(buf))
	})
	go http.ListenAndServe(fmt.Sprintf(":%d", *flagPort), mux)

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := cdp.New(ctxt, cdp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var res string
	err = c.Run(ctxt, setcookies(fmt.Sprintf("http://localhost:%d", *flagPort), &res))
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("passed cookies: %s", res)
	time.Sleep(1 * time.Minute)
}

func setcookies(host string, res *string) cdp.Tasks {
	return cdp.Tasks{
		cdp.ActionFunc(func(ctxt context.Context, h cdptypes.Handler) error {
			expr := cdptypes.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			success, err := network.SetCookie("cookiename", "cookievalue").
				WithExpires(&expr).
				WithDomain("localhost").
				WithHTTPOnly(true).
				Do(ctxt, h)
			if err != nil {
				return err
			}
			if !success {
				return errors.New("could not set cookie")
			}
			return nil
		}),
		cdp.Navigate(host),
		cdp.Text(`#result`, res, cdp.ByID, cdp.NodeVisible),
		cdp.ActionFunc(func(ctxt context.Context, h cdptypes.Handler) error {
			cookies, err := network.GetAllCookies().Do(ctxt, h)
			if err != nil {
				return err
			}
			fmt.Println("get all cookies:")
			for i, cookie := range cookies {
				log.Printf("cookie %d: %+v", i, cookie)
			}

			return nil
		}),
	}
}

const (
	indexHTML = `<!doctype html>
<html>
<body>
  <div id="result">%s</div>
</body>
</html>`
)

func click() cdp.Tasks {
	//var buf []byte
	return cdp.Tasks{
		cdp.Navigate(`https://golang.org/pkg/time/`),
		cdp.WaitVisible(`#footer`),
		cdp.Click(`#pkg-overview`, cdp.NodeVisible),
		//cdp.Screenshot(sel, res, cdp.NodeVisible, cdp.ByID),
		cdp.Sleep(150 * time.Second),
	}
}

func text(urlstr, sel string, res *string) cdp.Tasks {
	return cdp.Tasks{
		cdp.Navigate(urlstr),
		cdp.Sleep(3 * time.Second),
		cdp.WaitVisible(sel, cdp.ByID),
		cdp.Text(sel, res),
		//cdp.WaitNotVisible(`div.v-middle > div.la-ball-clip-rotate`, cdp.ByQuery),
		//cdp.Screenshot(sel, res, cdp.NodeVisible, cdp.ByID),
	}
}
