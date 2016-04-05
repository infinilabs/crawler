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

package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
    "github.com/medcl/gopa/core/code.google.com/p/go.net/websocket"
)

func ChatWith(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)

		//msg := "Received from " + ws.Request().Host + "  " + reply
		msg := "welcome to websocket do by pp"
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func Client(w http.ResponseWriter, r *http.Request) {
	html := `<!doctype html>
<html>

    <script type="text/javascript" src="http://img3.douban.com/js/packed_jquery.min6301986802.js" async="true"></script>
      <script type="text/javascript">
         var sock = null;
         var wsuri = "ws://127.0.0.1:8001";

         window.onload = function() {

            console.log("onload");


            try
            {
                sock = new WebSocket(wsuri);
            }catch (e) {
                alert(e.Message);
            }




            sock.onopen = function() {
               console.log("connected to " + wsuri);
            }

            sock.onerror = function(e) {
               console.log(" error from connect " + e);
            }



            sock.onclose = function(e) {
               console.log("connection closed (" + e.code + ")");
            }

            sock.onmessage = function(e) {
               console.log("message received: " + e.data);

               $('#log').append('<p> server say: '+e.data+'<p>');
               $('#log').get(0).scrollTop = $('#log').get(0).scrollHeight;
            }

         };

         function send() {
            var msg = document.getElementById('message').value;
            $('#log').append('<p style="color:red;">I say: '+msg+'<p>');
                $('#log').get(0).scrollTop = $('#log').get(0).scrollHeight;
                $('#msg').val('');
            sock.send(msg);
         };
      </script>
      <h1>WebSocket chat with server </h1>
          <div id="log" style="height: 300px;overflow-y: scroll;border: 1px solid #CCC;">
          </div>
          <div>
            <form>
                <p>
                    Message: <input id="message" type="text" value="Hello, world!"><button onclick="send();">Send Message</button>
                </p>
            </form>

          </div>

</html>`
	io.WriteString(w, html)
}

func StartWebsocket(){
	http.Handle("/", websocket.Handler(ChatWith))
	http.HandleFunc("/chat", Client)

	fmt.Println("listen on port 8001")
	fmt.Println("visit http://127.0.0.1:8001/chat with web browser(recommend: chrome)")

	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
