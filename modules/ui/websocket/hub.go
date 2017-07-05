// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/infinitbyte/gopa/core/env"
	"github.com/infinitbyte/gopa/core/global"
	"github.com/infinitbyte/gopa/core/logger"
	"github.com/infinitbyte/gopa/core/stats"
	"strings"
	"time"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	connections map[*WebsocketConnection]bool

	// Inbound messages from the connections.
	broadcast chan string

	// Register requests from the connections.
	register chan *WebsocketConnection

	// Unregister requests from connections.
	unregister chan *WebsocketConnection

	// Command handlers
	handlers map[string]WebsocketHandlerFunc
}

type WebsocketHandlerFunc func(c *WebsocketConnection, array []string)

var h = Hub{
	broadcast:   make(chan string, 5),
	register:    make(chan *WebsocketConnection),
	unregister:  make(chan *WebsocketConnection),
	connections: make(map[*WebsocketConnection]bool),
	handlers:    make(map[string]WebsocketHandlerFunc),
}

var runningHub = false

// Register command handlers
func (h *Hub) registerHandlers(env *env.Env) {
	handler := Command{}
	HandleWebSocketCommand("HELP", handler.Help)
	HandleWebSocketCommand("SEED", handler.AddSeed)
	HandleWebSocketCommand("LOG", handler.UpdateLogLevel)
	HandleWebSocketCommand("DIS", handler.Dispatch)
	HandleWebSocketCommand("GET_TASK", handler.GetTask)
}

func InitWebSocket(env *env.Env) {
	if !runningHub {
		h.registerHandlers(env)
		go h.RunHub()
	}
}

func HandleWebSocketCommand(cmd string, handler func(c *WebsocketConnection, array []string)) {
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	h.handlers[cmd] = WebsocketHandlerFunc(handler)
}

func (h *Hub) RunHub() {
	//TODO error　handler,　parameter　assertion

	if global.Env().IsDebug {

		go func() {
			t := time.NewTicker(time.Duration(30) * time.Second)
			for {
				select {
				case <-t.C:
					h.broadcast <- "testing websocket broadcast"
				}
			}
		}()
	}

	//handle connect, disconnect, broadcast
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			c.WritePrivateMessage(env.GetWelcomeMessage())
			js, _ := json.Marshal(logger.GetLoggingConfig())
			c.WriteMessage(ConfigMessage, string(js))
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.signalChannel)
			}
		case m := <-h.broadcast:
			h.broadcastMessage(m)
		}
	}

}

func (h *Hub) broadcastMessage(msg string) {

	if len(msg) <= 0 {
		return
	}

	for c := range h.connections {
		c.Broadcast(msg)
	}
}

func BroadcastMessage(msg string) {
	select {
	case h.broadcast <- msg:
		stats.Increment("websocket", "sended")
	default:
		stats.Increment("websocket", "dropped")
		fmt.Println("websocket broadcast too busy, msg dropped")
	}
}
