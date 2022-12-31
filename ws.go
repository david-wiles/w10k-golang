package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

// WebsocketConnection manages a single websocket between two threads, for reading and writing to the websocket.
type WebsocketConnection struct {
	// Channel used to send messages to this connection.
	// The messages will be sent using the websocket's write  thread
	Writes chan WebsocketMsg

	// Unique ID for this connection
	ID uuid.UUID

	// The underlying websocket connection
	conn *websocket.Conn
}

// WebsocketMsg represents a single message that should be written to a websocket.
type WebsocketMsg struct {
	// t is the websocket messageType
	t int

	// bytes are the literal bytes to write to the connection
	bytes []byte
}

// NewWebsocket creates the WebsocketConnection and starts threads
func NewWebsocket(socket *websocket.Conn, unregister chan string) *WebsocketConnection {
	ws := &WebsocketConnection{
		Writes: make(chan WebsocketMsg),
		ID:     uuid.New(),
		conn:   socket,
	}

	go ws.readWorker(unregister)
	go ws.writeWorker()

	return ws
}

// writeWorker will read messages from the Writes channel until closed
func (ws *WebsocketConnection) writeWorker() {
	for msg := range ws.Writes {
		if err := ws.conn.WriteMessage(msg.t, msg.bytes); err != nil {
			// Log errors
			log.Printf("conn %s error: %v\n", ws.ID.String(), err)
		}
	}
}

// readWorker will read messages from the websocket.Conn until the connection is closed. The Hub's Unregister channel
// is also passed to the worker so it is able to unregister itself whenever the websocket sends a close message.
func (ws *WebsocketConnection) readWorker(unregister chan string) {
	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err) {
				// Log errors
				log.Printf("conn %s error: %v\n", ws.ID.String(), err)
			}
			// Close this connection, mark as closed and close channel
			unregister <- ws.ID.String()
			_ = ws.conn.Close()
			return
		}
		// Handle message
		log.Printf("conn %s received %v\n", ws.ID.String(), message)
	}
}
