package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alextanhongpin/go-chat/database"
	"github.com/alextanhongpin/go-chat/ticket"
	"github.com/gorilla/websocket"
)

var (
	maxMessageSize int64 = 512
)

// https://github.com/gorilla/websocket/issues/46

type Message struct {
	Data  string `json:"data"`
	Room  string `json:"room"`
	Token string `json:"token"`
	Type  string `json:"type"`
	user  string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	broadcast chan Message
	clients   map[string]*websocket.Conn
	mapper    *Mapper
	quit      chan struct{}
}

func New() *Server {
	s := Server{
		broadcast: make(chan Message),
		clients:   make(map[string]*websocket.Conn),
		mapper:    NewMapper(),
		quit:      make(chan struct{}),
	}

	go s.eventloop()

	return &s
}

// Close terminates the server goroutines gracefully.
func (s *Server) Close() {
	close(s.quit)
}

// Broadcast sends a message to a client.
func (s *Server) Broadcast(to string, msg Message) error {
	// switch msg.Type {
	//         case
	// }
	if client, found := s.clients[to]; found {
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("error: %v\n", err)
			// If the delivery fails, remove the client from the
			// list.
			client.Close()

			// Delete all relationships.
			s.mapper.Delete(to)
			return err
		}
	}
	return nil
}

func (s *Server) eventloop() {
	for {
		select {
		case <-s.quit:
			log.Println("server: quit")
			return
		case msg := <-s.broadcast:
			switch msg.Type {
			// case "presence":
			// s.mapper.Has(msg.Room)
			// case "join_room":
			//         // Add the user into the room.
			//         s.mapper.Add(msg.Room, msg.user)

			default:
				log.Println("server: receive msg", msg)
				s.mapper.Add(msg.Room, msg.user)

				// Get the list of peers it can send message to.
				clients := s.mapper.Get(msg.Room)

				// Send only to clients in the particular room.
				for peer := range clients {
					log.Println("server: broadcasting message to peer", peer, msg)
					// This could be executed in a goroutine if the
					// users have many friends. Fanout operation.
					s.Broadcast(peer, msg)
				}
			}
		}
	}
}

func (s *Server) ServeWS(machine ticket.Dispenser, db database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// WebSocket is a httpGet only endpoint.
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// We can also perform checking of origin here.
		if r.Header.Get("Origin") != "http://"+r.Host {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
			return
		}

		// We can get the querystring parameter from the websocket
		// endpoint. This might be useful for validating parameters.
		q := r.URL.Query()
		token := q.Get("token")

		userID, err := machine.Verify(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := db.GetUser(userID)
		if err != nil || u.ID == 0 {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		user := strconv.Itoa(u.ID)
		// From here, we can get the top15 ranked friends and add them into the list.
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Make sure we close the connection when the function returns.
		defer ws.Close()

		// Add client to the session.
		s.clients[user] = ws
		defer func() {
			log.Println("server: remove session", user)
			// Remove client from the session.
			delete(s.clients, user)

			// Remove client from the listening peers.
			log.Println("server: delete relationships", user)
			s.mapper.Delete(user)

			// Broadcast the message to notify other peers that the user went offline.
		}()

		if user == "john" || user == "jane" {
			s.mapper.Add("room1", user)
		}

		// Read messages.
		ws.SetReadLimit(maxMessageSize)
		var msg Message
		// msg.From = user
		for {
			// Override the decision here.
			// msg.Room = "room1"
			msg.user = user
			if err := ws.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("error: %v, user-agent: %v", err, r.Header.Get("User-Agent"))
				}
				return
			}
			s.broadcast <- msg
		}
	}
}
