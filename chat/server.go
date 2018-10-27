package chat

import (
	"log"
	"net/http"
	"time"

	"github.com/alextanhongpin/go-chat/database"
	"github.com/alextanhongpin/go-chat/repository"
	"github.com/alextanhongpin/go-chat/ticket"
	"github.com/gorilla/websocket"
)

var (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer
	maxMessageSize int64 = 512
)

// https://github.com/gorilla/websocket/issues/46

type Message struct {
	Data string `json:"data,omitempty"`
	Room string `json:"room,omitempty"`
	Type string `json:"type,omitempty"`
	// To and From value can be hashed for security purpose. Create another
	// lookup table to map the values back to the original id.
	To   string `json:"to,omitempty"`
	From string `json:"from,omitempty"`
	user string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin: func(r *http.Request) bool {
	//         return r.Header.Get("Origin") != "http://"+r.Host
	// },
}

type Server struct {
	broadcast chan Message
	clients   map[string]*websocket.Conn
	quit      chan struct{}
	db        *database.Conn
	cache     repository.UserCache
}

func New(db *database.Conn) *Server {
	s := Server{
		cache:     NewCache(),
		broadcast: make(chan Message),
		clients:   make(map[string]*websocket.Conn),
		quit:      make(chan struct{}),
		db:        db,
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
	if client, found := s.clients[to]; found {
		err := client.SetWriteDeadline(time.Now().Add(writeWait))
		if err != nil {
			client.Close()
			s.cache.RemoveUser(to)
			return err
		}
		if err := client.WriteJSON(msg); err != nil {
			log.Printf("error: %v\n", err)
			// If the delivery fails, remove the client from the
			// list.
			client.Close()

			// Delete all relationships.
			// s.rooms.Del(to)
			s.cache.RemoveUser(to)
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
			case "is_typing":
				err := s.Broadcast(msg.To, Message{
					// Data: msg.user,
					Room: msg.Room,
					Type: msg.Type,
				})
				if err != nil {
					log.Println("authError:", err)
				}
			case "auth":
				err := s.Broadcast(msg.user, Message{
					Data: msg.user,
					Type: msg.Type,
				})
				if err != nil {
					log.Println("authError:", err)
				}
			case "status":
				// Data is the user_id that we want to check the status of.
				user := msg.Data
				_, found := s.clients[user]
				data := "0"
				if found {
					data = "1"
				}
				s.Broadcast(msg.user, Message{
					Data: data,
					Type: "status",
					Room: msg.Room,
					From: user,
				})
			case "presence":
				clients := s.cache.GetUsers(msg.Room)

				// Send only to clients in the particular room.
				for _, peer := range clients {
					log.Println("server: broadcasting message to peer", peer, msg)
					// This could be executed in a
					// goroutine if the users have many
					// friends. Fanout operation.
					s.Broadcast(peer, msg)
				}
			case "message":
				// s.rooms.Add(msg.user, msg.Room)
				s.cache.AddUser(msg.user, msg.Room)

				// Store the conversation in a database. It
				// might be a better idea to use a queue rather
				// than writing directly to the datastore.
				_, err := s.db.CreateConversationReply(msg.user, msg.Room, msg.Data)
				if err != nil {
					log.Printf("error: conversation create error, %v\n", err)
					continue
				}

				// Get the list of peers it can send message to.
				// clients := s.rooms.GetUsers(msg.Room)
				clients := s.cache.GetUsers(msg.Room)

				// Send only to clients in the particular room.
				for _, peer := range clients {
					log.Println("server: broadcasting message to peer", peer, msg)
					// This could be executed in a goroutine if the
					// users have many friends. Fanout operation.
					s.Broadcast(peer, msg)
				}
			default:
				log.Printf("message type %s not supported\n", msg.Type)
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

		// We can get the querystring parameter from the websocket
		// endpoint. This might be useful for validating parameters.
		token := r.URL.Query().Get("token")
		userID, err := machine.Verify(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := db.GetUser(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		// Get the hashed ID
		// user := u.ID
		// From here, we can get the top15 ranked friends and add them into the list.
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Make sure we close the connection when the function returns.
		defer ws.Close()

		log.Println("authenticated", user.ID)
		// Add client to the session.
		s.clients[user.ID] = ws
		defer delete(s.clients, user.ID)

		// Notify other user in the room that the user went online.
		s.online(user.ID)
		defer s.offline(user.ID)

		ws.SetReadLimit(maxMessageSize)
		for {
			var msg Message
			if err := ws.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("error: %v, user-agent: %v", err, r.Header.Get("User-Agent"))
				}
				return
			}
			msg.From = user.ID
			msg.user = user.ID
			s.broadcast <- msg
		}
	}
}

func (s *Server) online(user string) error {
	rooms, err := s.db.GetRoom(user)
	if err != nil {
		return err
	}
	for _, room := range rooms {

		// Notify other users in the room that the user went online.
		s.broadcast <- Message{
			Type: "presence",
			Room: room,
			Data: "1",
		}

		// Add user to the room after broadcasting to the room - to
		// avoid notifying oneself.
		s.cache.AddUser(user, room)
	}
	return nil
}

func (s *Server) offline(user string) {
	rooms := s.cache.GetRooms(user)
	for _, room := range rooms {
		msg := Message{
			Type: "presence",
			Room: room,
			Data: "0",
		}
		s.broadcast <- msg
	}
	s.cache.RemoveUser(user)
}
