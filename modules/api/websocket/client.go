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

package websocket

//
//import (
//	"github.com/gorilla/websocket"
//	"log"
//	"net/http"
//	"time"
//)
//
//type client struct {
//	ws   *websocket.Conn
//	send chan []byte
//}
//
//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  maxMessageSize,
//	WriteBufferSize: maxMessageSize,
//}
//
//func Broadcast(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "GET" {
//		http.Error(w, "Method not allowed", 405)
//		return
//	}
//
//	ws, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	c := &client{
//		send: make(chan []byte, maxMessageSize),
//		ws:   ws,
//	}
//
//	h.register <- c
//
//	c.broadcast("hello world")
//}
//
//func (c *client) broadcast(msg string) {
//	defer func() {
//		h.unregister <- c
//		c.ws.Close()
//	}()
//	h.broadcast <- msg
//}
//
//func (c *client) write(mt int, message []byte) error {
//	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
//	return c.ws.WriteMessage(mt, message)
//}
