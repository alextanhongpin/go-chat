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
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Server returns a new chat server
func Server(redisPort string) http.HandlerFunc {

	room := NewRoom(redisPort)
	go room.Run()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		log.Println("you are in room:", query.Get("room"))
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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

		go subscription.Read(room.PubSub)
		go subscription.Write()

		// Create a new client here / subscribe to a new redis channel here
		room.PubSub.Subscribe(query.Get("room"), subscription)
	})
}
