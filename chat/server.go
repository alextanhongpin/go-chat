package chat

import (
	"context"
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
func (s *Server) Run(ctx context.Context) {
	s.Room.Run(ctx)
}

// ServeWS returns a handler function for the websocket connection
func (s *Server) ServeWS() http.HandlerFunc {
	// room := NewRoom()
	// go room.Run(context.Background())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		room := s.Room

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

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

		client := NewClient(ws, room)
		// defer client.Conn.Close()
		subscription := NewSubscription(roomID, client)

		client.Subscribe(subscription)

		go subscription.Read(ctx, s.PubSub)
		go subscription.Write(ctx)

		// Create a new client here / subscribe to a new redis channel here
		s.PubSub.Subscribe(ctx, roomID, subscription)
	})
}

// NewServer returns a pointer to the server
func NewServer(port string) *Server {
	return &Server{
		Room:   NewRoom(),
		PubSub: NewPubSub(port),
	}
}
