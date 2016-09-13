// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type WebsocketConnection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	handlers map[string]WebsocketHandlerFunc
}

// readPump pumps messages from the websocket connection to the hub.
func (c *WebsocketConnection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Errorf("error: %v", err)
			}
			break
		}
		h.broadcast <- message
	}
}

// write writes a message with the given message type and payload.
func (c *WebsocketConnection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *WebsocketConnection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.parseMessage(message)

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *WebsocketConnection) WriteMessage(msg []byte)error {
	return c.write(websocket.TextMessage, msg)
}

func (c *WebsocketConnection) parseMessage(msg []byte) {

	message := string(msg)
	array := strings.Split(message, " ")
	if len(array) > 0 {
		cmd := strings.ToLower(strings.TrimSpace(array[0]))
		if c.handlers!=nil {
			handler := c.handlers[cmd]
			if handler != nil {
				handler(c,array)
				return
			}
		}
	}

	if err := c.write(websocket.TextMessage, []byte("\nCommand not found\n")); err != nil {
		return
	}

}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	c := &WebsocketConnection{send: make(chan []byte, 256), ws: ws, handlers: h.handlers}
	h.register <- c
	go c.writePump()
	c.readPump()
}
