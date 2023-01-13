package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebsocketConnection manages a single websocket between two threads, for reading and writing to the websocket.
type WebsocketConnection struct {
	// Channel used to send messages to this connection.
	// The messages will be sent using the websocket's write  thread
	Writes chan WebsocketMsg

	// Unique ID for this connection
	ID uuid.UUID

	// The underlying websocket connection
	conn *websocket.Conn

	// Reference to the hub managing this connection
	hub *Hub
}

// WebsocketMsg represents a single message that should be written to a websocket.
type WebsocketMsg struct {
	// t is the websocket messageType
	T int

	// bytes are the literal bytes to write to the connection
	Bytes []byte
}

// NewWebsocket creates the WebsocketConnection and starts threads
func NewWebsocket(socket *websocket.Conn, hub *Hub) *WebsocketConnection {
	ws := &WebsocketConnection{
		Writes: make(chan WebsocketMsg),
		ID:     uuid.New(),
		conn:   socket,
		hub:    hub,
	}

	go ws.readWorker()
	go ws.writeWorker()

	return ws
}

// writeWorker will read messages from the Writes channel until closed
func (ws *WebsocketConnection) writeWorker() {
	for msg := range ws.Writes {
		if err := ws.conn.WriteMessage(msg.T, msg.Bytes); err != nil {
			// Log errors
			log.Printf("conn %s error: %v\n", ws.ID.String(), err)
		}
	}
}

// readWorker will read messages from the websocket.Conn until the connection is closed. The Hub's Unregister channel
// is also passed to the worker so it is able to unregister itself whenever the websocket sends a close message.
func (ws *WebsocketConnection) readWorker() {
	for {
		t, message, err := ws.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err) {
				// Log errors
				log.Printf("conn %s error: %v\n", ws.ID.String(), err)
			}
			// Close this connection, mark as closed and close channel
			ws.hub.Unregister <- ws.ID.String()
			_ = ws.conn.Close()
			return
		}

		// Handle message
		ws.hub.Writes <- WebsocketMsg{t, message}
	}
}

// Hub is the manager for all websocket connections. Based on the implementation in gorilla/websocket's chat example
// https://github.com/gorilla/websocket/blob/76ecc29eff79f0cedf70c530605e486fc32131d1/examples/chat/hub.go
//
// The Hub has a channel to register new connections and a channel to deregister connections. This ensures consistency
// of the connections stored in the map during the ping interval. Each connection is assigned a UUID, and that UUID is
// the key of the map.
type Hub struct {
	Register   chan *WebsocketConnection
	Unregister chan string
	Writes     chan WebsocketMsg

	connections map[string]*WebsocketConnection
}

// NewHub creates a Hub
func NewHub() *Hub {
	return &Hub{
		connections: make(map[string]*WebsocketConnection),
		Register:    make(chan *WebsocketConnection),
		Unregister:  make(chan string),
		Writes:      make(chan WebsocketMsg, 1024),
	}
}

// Monitor starts a loop which waits for messages on any channel. This will loop infinitely, and so should not be used
// where it would block any other processing.
func (hub *Hub) Monitor() {
	for {
		select {
		case ws := <-hub.Register:
			hub.connections[ws.ID.String()] = ws
		case id := <-hub.Unregister:
			if ws, ok := hub.connections[id]; ok {
				close(ws.Writes)
				delete(hub.connections, id)
			}
		case msg := <-hub.Writes:
			// Whenever we receieve a message, send it to the corresponding connection
			if msg.T == websocket.TextMessage {
				text := string(msg.Bytes)
				if len(text) > 36 {
					if _, err := uuid.Parse(text[:36]); err == nil {
						if dest, ok := hub.connections[text[:36]]; ok {
							dest.Writes <- WebsocketMsg{
								T:     websocket.TextMessage,
								Bytes: []byte(text[36:]),
							}
						}
					}
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
	ws := NewWebsocket(conn, hub)
	hub.Register <- ws
	return ws.ID.String()
}

// handleWebsocket will upgrade the request to a websocket connection and add the connection to the Hub
func handleWebsocket(hub *Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade error: %v\n", err)
			return
		}

		id := hub.RegisterConnection(conn)
		log.Printf("registering new websocket %s\n", id)
	}
}

func main() {
	mux := http.NewServeMux()
	hub := NewHub()

	go hub.Monitor()

	// All requests are handled as websocket requests
	mux.HandleFunc("/ws", handleWebsocket(hub))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
