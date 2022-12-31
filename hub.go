package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// Hub is the manager for all websocket connections. Based on the implementation in gorilla/websocket's chat example
// https://github.com/gorilla/websocket/blob/76ecc29eff79f0cedf70c530605e486fc32131d1/examples/chat/hub.go
//
// The Hub has a channel to register new connections and a channel to deregister connections. This ensures consistency
// of the connections stored in the map during the ping interval. Each connection is assigned a UUID, and that UUID is
// the key of the map.
type Hub struct {
	Register    chan *WebsocketConnection
	Unregister  chan string
	connections map[string]*WebsocketConnection
}

// NewHub creates a Hub
func NewHub() *Hub {
	return &Hub{
		connections: make(map[string]*WebsocketConnection),
		Register:    make(chan *WebsocketConnection),
		Unregister:  make(chan string),
	}
}

// Monitor starts a loop which waits for messages on any channel. This will loop infinitely, and so should not be used
// where it would block any other processing.
func (hub *Hub) Monitor(interval time.Duration) {
	// Broadcast a message to all clients every interval
	ticker := time.NewTicker(interval)

	for {
		select {
		case ws := <-hub.Register:
			hub.connections[ws.ID.String()] = ws
		case id := <-hub.Unregister:
			if ws, ok := hub.connections[id]; ok {
				close(ws.Writes)
				delete(hub.connections, id)
			}
		case t := <-ticker.C:
			msg := []byte(t.Format(time.RFC3339Nano))
			for _, ws := range hub.connections {
				select {
				case ws.Writes <- WebsocketMsg{websocket.TextMessage, msg}:
					// Success
				default:
					close(ws.Writes)
					delete(hub.connections, ws.ID.String())
				}
			}
		default:
			// Nothing to do
		}
	}
}

// RegisterConnection takes a raw websocket connection and creates a new WebsocketConnection object. This also gives
// the Unregister channel to the WebsocketConnection so it is able to unregister itself when needed.
func (hub *Hub) RegisterConnection(conn *websocket.Conn) string {
	ws := NewWebsocket(conn, hub.Unregister)
	hub.Register <- ws
	return ws.ID.String()
}
