package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Server returns a new chat server
func Server() http.HandlerFunc {

	room := NewRoom()
	go room.Run()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		log.Println("you are in room:", query.Get("room"))
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		client := NewClient(ws, room)
		subscription := NewSubscription(query.Get("room"), client)
		client.Subscribe(subscription)

		go subscription.Read()
		go subscription.Write()
	})
}
