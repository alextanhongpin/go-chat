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

// Server represents the chat server
type Server struct {
	Room   *Room
	PubSub *PubSub
}

// Run starts the server
func (s *Server) Run() {
	go s.Room.Run()
	go s.PubSub.Subscribe(s.Room)
}

// ServeWS returns a handler function for the websocket connection
func (s *Server) ServeWS() http.HandlerFunc {
	log.Println("serving websocket")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		room := s.Room

		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()
		roomID := query.Get("room")

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		client := NewClient(ws)
		subscription := NewSubscription(roomID, client)
		room.Subscribe <- subscription

		go subscription.Read(s.PubSub, room)
		subscription.Write()
	})
}

// NewServer returns a pointer to the server
func NewServer(port, channel string) *Server {
	return &Server{
		Room:   NewRoom(),
		PubSub: NewPubSub(port, channel),
	}
}
