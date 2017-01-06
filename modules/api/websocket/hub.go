// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	ws "github.com/gorilla/websocket"
	"github.com/medcl/gopa/core/env"
	"strings"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	connections map[*WebsocketConnection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *WebsocketConnection

	// Unregister requests from connections.
	unregister chan *WebsocketConnection

	// Command handlers
	handlers map[string]WebsocketHandlerFunc
}

type WebsocketHandlerFunc func(c *WebsocketConnection, array []string)

var h = Hub{
	broadcast:   make(chan []byte),
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
	//TOD error　handler,　parameter　assertion

	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			c.write(ws.TextMessage, []byte(env.GetWelcomeMessage()))
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}
