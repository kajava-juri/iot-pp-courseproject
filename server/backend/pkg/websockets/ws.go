package websockets

import (
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var hub = NewWsHub()
var mutex = &sync.Mutex{}

func handler(hub *WsHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	client := hub.NewClient(conn)
	go client.WriteMessages() // Start the write goroutine
	hub.register <- client    // Register with hub instead of direct map access
	defer func() {
		hub.unregister <- client // Unregister when function exits
	}()
	mutex.Unlock()

	// Handle WebSocket connection

	for {

		_, msg, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err) {
				log.Println("Error reading message:", err)
			} else {
				log.Printf("WebSocket connection closed: %s", conn.RemoteAddr())
			}
			hub.unregister <- client
			break
		}
		log.Println("Received message:", string(msg))

		hub.handleClientMessage(client, msg)

		//hub.broadcast <- msg
	}
}

func RegisterRoutes(r chi.Router, hub *WsHub) {
	r.Get("/ws", func(w http.ResponseWriter, req *http.Request) {
		handler(hub, w, req)
	})
}
