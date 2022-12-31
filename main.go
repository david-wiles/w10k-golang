package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var interval = time.Second

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
	if d, ok := os.LookupEnv("PING_INTERVAL"); ok {
		dur, err := time.ParseDuration(d)
		if err != nil {
			panic("PING_INTERVAL must be a duration")
		}
		interval = dur
	}

	mux := http.NewServeMux()

	hub := NewHub()
	go hub.Monitor(interval)

	// All requests are handled as websocket requests
	mux.HandleFunc("/ws", handleWebsocket(hub))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
