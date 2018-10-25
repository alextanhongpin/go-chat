package chat

import (
	"log"
	"net/http"
	"strconv"
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
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
		// rooms:     NewRoomManager(), // TODO: Add bloom filter (?).
		db: db,
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
				// clients := s.rooms.GetUsers(msg.Room)

				// Send only to clients in the particular room.
				for _, peer := range clients {
					log.Println("server: broadcasting message to peer", peer, msg)
					// This could be executed in a goroutine if the
					// users have many friends. Fanout operation.
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

		// We can also perform checking of origin here.
		if r.Header.Get("Origin") != "http://"+r.Host {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
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
		defer delete(s.clients, user)

		// Notify other user in the room that the user went online.
		s.online(user)
		defer s.offline(user)

		ws.SetReadLimit(maxMessageSize)
		for {
			var msg Message
			msg.user = user
			// TODO: Hash the id.
			msg.From = user
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

func (s *Server) online(user string) error {
	rooms, err := s.db.GetRoom(user)
	if err != nil {
		return err
	}
	for _, room := range rooms {
		roomID := strconv.FormatInt(room, 10)

		// Notify other users in the room that the user went online.
		s.broadcast <- Message{
			Type: "presence",
			Room: roomID,
			Data: "1",
		}

		// Add user to the room after broadcasting to the room - to
		// avoid notifying oneself.
		// s.rooms.Add(user, roomID)
		s.cache.AddUser(user, roomID)
	}
	return nil
}

func (s *Server) offline(user string) {
	// rooms := s.rooms.GetRooms(user)
	rooms := s.cache.GetRooms(user)
	for _, room := range rooms {
		msg := Message{
			Type: "presence",
			Room: room,
			Data: "0",
		}
		s.broadcast <- msg
	}
	// s.rooms.Del(user)
	s.cache.RemoveUser(user)
}
