package chat

import (
	"net/http"

	"github.com/alextanhongpin/go-chat/ticket"
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
	s.Room.Run()
}

// Subscribe to the room
func (s *Server) Subscribe() {
	s.PubSub.Subscribe(s.Room)
}

// ServeWS returns a handler function for the websocket connection
func (s *Server) ServeWS() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		room := s.Room

		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()
		token := query.Get("ticket")

		// Verify that the token is valid
		tic, err := ticket.Verify(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		roomID := tic.ID
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		client := NewClient(ws, roomID)
		room.Subscribe <- client

		go client.Read(s.PubSub, room)
		client.Write()
	})
}

// NewServer returns a pointer to the server
func NewServer(port, channel string) *Server {
	return &Server{
		Room:   NewRoom(),
		PubSub: NewPubSub(port, channel),
	}
}
