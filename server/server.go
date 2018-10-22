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
	To    string `json:"to"`
	From  string `json:"from"`
	user  string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	broadcast chan Message
	clients   map[string]*websocket.Conn
	rooms     RoomManager
	quit      chan struct{}
	db        *database.Conn
}

func New(db *database.Conn) *Server {
	s := Server{
		broadcast: make(chan Message),
		clients:   make(map[string]*websocket.Conn),
		quit:      make(chan struct{}),
		rooms:     NewRoomManager(), // TODO: Add bloom filter (?).
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
			s.rooms.Del(to)
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
			case "presence":
				clients := s.rooms.GetUsers(msg.Room)

				// Send only to clients in the particular room.
				for _, peer := range clients {
					log.Println("server: broadcasting message to peer", peer, msg)
					// This could be executed in a goroutine if the
					// users have many friends. Fanout operation.
					s.Broadcast(peer, msg)
				}
			default:
				log.Println("server: receive msg", msg)
				s.rooms.Add(msg.user, msg.Room)

				// Store the conversation in a database. It
				// might be a better idea to use a queue rather
				// than writing directly to the datastore.
				_, err := s.db.CreateConversationReply(msg.user, msg.Room, msg.Data)
				if err != nil {
					log.Printf("error: conversation create error, %v\n", err)
					continue
				}

				// Get the list of peers it can send message to.
				clients := s.rooms.GetUsers(msg.Room)

				// Send only to clients in the particular room.
				for _, peer := range clients {
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
		cleanup := func() {
			log.Println("server: remove session", user)
			// Remove client from the session.
			delete(s.clients, user)

			// Remove client from the listening peers.
			log.Println("server: delete relationships", user)

			// Notify other users that the user went offline.
			rooms := s.rooms.GetRooms(user)
			for _, room := range rooms {
				msg := Message{
					Room: room,
					// Data: fmt.Sprintf("%s went offline", user),
					Data: "0",
					Type: "presence",
				}
				s.broadcast <- msg
			}
			s.rooms.Del(user)
		}
		defer cleanup()

		// Notify other party that this user is online.
		rooms, err := s.db.GetRoom(user)
		if err == nil {
			log.Println("adding users into room")
			for _, r := range rooms {
				room := strconv.Itoa(int(r))
				s.rooms.Add(user, room)
				msg := Message{
					Room: room,
					// Data: fmt.Sprintf("%s went online", user),
					Data: "1",
					From: user,
					Type: "presence",
				}
				s.broadcast <- msg
			}
		}

		// Read messages.
		ws.SetReadLimit(maxMessageSize)

		for {
			// Override the decision here.
			var msg Message
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
