package chat

import (
	"log"
	"net/http"
	"time"

	"github.com/alextanhongpin/go-chat/database"
	"github.com/alextanhongpin/go-chat/ticket"
	"github.com/go-redis/redis"

	"github.com/gorilla/websocket"
)

var (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer
	maxMessageSize int64 = 512

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Message struct {
	// The text content of the message.
	Text string `json:"text"`

	// Unique ID for the message created.
	ID string `json:"id"`

	// The type for client to handle polymorphism.
	Type string `json:"type"`

	Timestamp int64 `json:"timestamp"`
	// User might not persist timestamp that long (when messages are loaded from db).
	// k = HMAC(timestamp,id,msg, sk)
	// timestamp|id|(msg|sender)k|HMAC(timestamp|id|sender|msg, sk))
	// Hash string

	// The sender of the message. Usually a user id.
	Sender string `json:"sender"`

	// The receiver in this case is the room id (chat group), rather than
	// the user id. This provides better flexibility, as it allows us to
	// send conversations to more than a person (users in the same group,
	// rather than just direct message to a single user).
	Receiver string `json:"receiver"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
		// return r.Header.Get("Origin") != "http://"+r.Host
	},
}

type Chat struct {
	// broadcast sends the message to a room.
	broadcast chan Message
	quit      chan struct{}

	// In-memory data structure, fetch from db if it doesn't exist.
	// Table that maps session id -> session struct. One-to-one.
	sessions *Sessions

	// Table that maps session id -> user id and vice versa. Many-to-one.
	lookup Tabler
	// Table that maps user id -> room ids and vice versa. Many-to-many.
	rooms Tabler
	db    *database.Conn
}

func New(db *database.Conn, client *redis.Client) *Chat {
	defaultBroadcastQueueSize := 10000
	c := Chat{
		broadcast: make(chan Message, defaultBroadcastQueueSize),
		quit:      make(chan struct{}),
		sessions:  NewSessions(),
		lookup:    NewTableInMemory(),
		rooms:     NewTableCache(client),
		db:        db,
	}

	// Register to pubsub to listen to the server.
	// os.Hostname()

	log.Println("chat: starting event loop")
	go c.eventloop()

	return &c
}

// Close terminates the server goroutines gracefully.
func (c *Chat) Close() {
	// This is preferred, as it will block unlike close.
	c.quit <- struct{}{}
	close(c.quit)
	log.Println("chat: closing")
}

// Broadcast sends a message to a client.
func (c *Chat) Broadcast(msg Message) error {
	log.Printf("broadcast: got message %#v\n", msg)
	sender, receiver := msg.Sender, msg.Receiver

	// Get the other users in the same room.
	users := c.rooms.Get(RoomID(receiver))
	log.Printf("broadcast: got users %#v\n", users)
	for _, user := range users {
		// Can skip this, since the user is already removed from the room.
		if user == sender {
			log.Printf("broadcast: skip ownself %#v\n", user)
			continue
		}
		// Find the sessions for the user.
		sessions := c.lookup.Get(UserID(user))
		log.Printf("broadcast: got sessions %#v\n", sessions)
		for _, sid := range sessions {
			sess := c.sessions.Get(sid)
			err := sess.Conn().WriteJSON(msg)
			if err != nil {
				// Clear session.
				c.Clear(sess)
				log.Printf("broadcastError: %v\n", err)
				return err
			}
		}
	}
	return nil
}

func (c *Chat) eventloop() {
loop:
	for {
		select {
		case <-c.quit:
			log.Println("server: quit")
			break loop
		case msg, ok := <-c.broadcast:
			if !ok {
				break loop
			}
			c.Broadcast(msg)
			switch msg.Type {
			// case TypeTyping:
			// err := s.Broadcast()
			// if err != nil {
			//         log.Println("typingError:", err)
			// }
			// case TypeAuth:
			// err := s.Broadcast(msg.From, Message{
			//         Data: msg.From,
			//         Type: msg.Type,
			// })
			// if err != nil {
			//         log.Println("authError:", err)
			// }
			// case TypeStatus:
			// Data is the user_id that we want to check the status of.
			// user := msg.Data
			// _, found := s.clients[user]
			// data := "0"
			// if found {
			//         data = "1"
			// }
			// s.Broadcast(msg.From, Message{
			//         Data: data,
			//         Type: "status",
			//         Room: msg.Room,
			//         // From: user,
			// })
			// case TypePresence:
			// clients := s.cache.GetUsers(msg.Room)
			//
			// // Send only to clients in the particular room.
			// for _, peer := range clients {
			//         log.Println("server: broadcasting message to peer", peer, msg)
			//         // This could be executed in a
			//         // goroutine if the users have many
			//         // friends. Fanout operation.
			//         s.Broadcast(peer, msg)
			// }
			// case TypeMessage:
			// s.cache.AddUser(msg.From, msg.Room)
			//
			// // Store the conversation in a database. It
			// // might be a better idea to use a queue rather
			// // than writing directly to the datastore.
			// _, err := s.db.CreateConversationReply(msg.From, msg.Room, msg.Data)
			// if err != nil {
			//         log.Printf("error: conversation create error, %v\n", err)
			//         continue
			// }
			//
			// // Get the list of peers it can send message to.
			// // clients := s.rooms.GetUsers(msg.Room)
			// clients := s.cache.GetUsers(msg.Room)
			//
			// // Send only to clients in the particular room.
			// for _, peer := range clients {
			//         log.Println("server: broadcasting message to peer", peer, msg)
			//         // This could be executed in a goroutine if the
			//         // users have many friends. Fanout operation.
			//         s.Broadcast(peer, msg)
			// }
			// default:
			// log.Printf("message type %s not supported\n", msg.Type)
			}
		}
	}
}

func (c *Chat) newSession(ws *websocket.Conn) *Session {
	sess := NewSession(ws)
	c.sessions.Put(sess)
	return sess
}

// Bind ties the user id and session id together. One user might have multiple sessions.
func (c *Chat) Bind(uid UserID, sid SessionID) func() {
	// TODO: Check if the session already exist. If yes, there is no need to add the
	// user into the room.
	// if !c.lookup.Has(uid) {
	// Connect the user to the room.
	c.Join(uid)
	// }

	// Tie the user to the existing session.
	c.lookup.Add(uid, sid)

	return func() {
		// Clear the current session that is tied to the user.
		sessions := c.Get(sid)
		for _, sess := range sessions {
			// Clear session table.
			c.Clear(sess)
		}

		// If the user does not have any other sessions left, clear the room.
		// One user can have multiple sessions.
		if len(c.Get(uid)) == 0 {
			// Clear rooms. Only if there are no longer any sessions available for that particular user.
			c.Leave(uid)
		}
	}
}

func (c *Chat) ServeWS(dispenser ticket.Dispenser, db database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// WebSocket is a httpGet only endpoint.
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Upgrade the websocket connection.
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Make sure we close the connection when the function returns.
		defer ws.Close()

		token := r.URL.Query().Get("token")
		userID, err := dispenser.Verify(token)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "unauthorized"))
			return
		}

		// Check if the user exists in the database.
		_, err = db.GetUser(userID)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "unauthorized"))
			return
		}

		// Create a new session. Why not create the session together with the user id?
		// A user can have several sessions (multiple tabs, different devices etc).
		// We need a way to query the sessions for a particular user.
		session := c.newSession(ws)

		// Check the db and get the user info, then tie them together.
		close := c.Bind(UserID(userID), session.SessionID())
		defer close()

		// Notify other user in the room that the user went online.
		ws.SetReadLimit(maxMessageSize)
		ws.SetReadDeadline(time.Now().Add(pongWait))
		ws.SetPongHandler(func(string) error {
			ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		go ping(ws)
		for {
			var msg Message
			if err := ws.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("error: %v, user-agent: %v", err, r.Header.Get("User-Agent"))
				}
				return
			}
			// Every websocket connection is unique - we can safely
			// inject the user id to the message.
			msg.Sender = userID
			c.broadcast <- msg
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func ping(ws *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Join adds the user to the room.
func (c *Chat) Join(uid UserID) {
	log.Printf("%v join the room\n", uid)
	rooms, err := c.db.GetRooms(string(uid))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Join: rooms %#v\n", rooms)
	for _, room := range rooms {
		// Notify other user in the room first, only add the user once the room receive broadcast to avoid notifying oneself.
		sender, receiver := string(uid), room.RoomID

		// Broadcast to a room.
		c.broadcast <- Message{
			Type:     MessageTypeOnline.String(),
			Sender:   sender,
			Receiver: receiver,
		}

		// Add user to room, and keep track of rooms for user.
		err := c.rooms.Add(uid, RoomID(room.RoomID))
		if err != nil {
			log.Printf("JoinError: %v\n", err)
		}
		log.Printf("Join: added to rooms %#v\n", room.RoomID)
	}
	roomss := c.rooms.Get(uid)
	log.Printf("Join: get rooms %#v\n", roomss)
}

func (c *Chat) Leave(uid UserID) {
	log.Printf("%v left the chat\n", uid)
	// For each room that the user belong to, remove the user.
	onDelete := func(room string) {
		sender, receiver := string(uid), room
		c.broadcast <- Message{
			Type:     MessageTypeOffline.String(),
			Sender:   sender,
			Receiver: receiver,
		}
	}
	// Delete user -> rooms relationship.
	c.rooms.Delete(uid, onDelete)

}

func (c *Chat) Clear(sess *Session) {
	// Closes the connection, and delete the session.
	sess.Conn().Close()
	sessionID := sess.SessionID()
	c.sessions.Delete(string(sessionID))
}

// One-to-many relationship between sessions and user.
// One user can have multiple sessions (mobile, browser with multiple tabs etc)
// Get(UserID) will return a slice of sessions.
// Get(SessionID) will return one user.
// To get the session.
func (c *Chat) Get(key interface{}) []*Session {
	switch key.(type) {
	case SessionID:
		sess := c.sessions.Get(key.(string))
		return []*Session{sess}
	case UserID:
		// If the UserID is provided, get the SessionID first
		// in order to retrieve the session.
		// A userID can have multiple sessions (many tabs)
		sessions := c.lookup.Get(key)
		result := make([]*Session, len(sessions))

		for i, sess := range sessions {
			result[i] = c.sessions.Get(sess)
		}
		return result
	default:
		return nil
	}
}
