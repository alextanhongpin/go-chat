package chat

import (
	"log"
	"net/http"
	"time"

	"github.com/alextanhongpin/go-chat/database"
	"github.com/alextanhongpin/go-chat/ticket"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

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

	defaultBroadcastQueueSize = 10000
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
		// return r.Header.Get("Origin") != "http://"+r.Host
	},
}

// Chat represents the chat application.
type Chat struct {
	// broadcast sends the message to a room.
	broadcast chan Message

	// quit signals termination of the goroutine that is handling the
	// broadcast.
	quit chan struct{}

	// In-memory data structure, fetch from db if it doesn't exist.
	// Table that maps session id -> session struct. One-to-one.
	sessions *Sessions

	// Table that maps session id -> user id and vice versa. Many-to-one.
	lookup *Table

	// Table that maps user id -> room ids and vice versa. Many-to-many.
	rooms  *TableCache
	db     *database.Conn
	logger *zap.Logger
}

func New(db *database.Conn, client *redis.Client, logger *zap.Logger) *Chat {
	c := Chat{
		broadcast: make(chan Message, defaultBroadcastQueueSize),
		quit:      make(chan struct{}),
		sessions:  NewSessions(),
		lookup:    NewTableInMemory(),
		rooms:     NewTableCache(client),
		db:        db,
		logger:    logger.With(zap.String("pkg", "Chat")),
	}

	// Register to pubsub to listen to the server.
	// os.Hostname()

	// log.Println("chat: starting event loop")
	logger.Info("starting event loop")
	go c.eventloop()

	return &c
}

// Close terminates the server goroutines gracefully.
func (c *Chat) Close() {
	// This is preferred, as it will block unlike close.
	c.quit <- struct{}{}
	close(c.quit)
	c.logger.Info("closing")
}

// Broadcast sends a message to a client.
func (c *Chat) Broadcast(msg Message) error {
	logger := c.logger.With(zap.String("method", "Broadcast"))
	logger.Info("got message",
		zap.String("receiver", msg.Receiver),
		zap.String("sender", msg.Sender))
	sender, receiver := msg.Sender, msg.Receiver
	_ = sender

	// Get all users in the same room.
	users := c.rooms.GetUsers(receiver)
	for _, user := range users {
		// sendToSelf: bool
		// If set to true, allow sending to one-self. Else, skip.
		// if user == sender {
		//         logger.Info("skip ownself")
		//         continue
		// }

		// Get the sessions for the user. A user can have more than one
		// session.
		sessions := c.lookup.GetSessions(user)
		for _, sid := range sessions {
			sess := c.sessions.Get(sid)
			// Session happens to be empty, but still in the list.
			if sess == nil {
				// Delete the mapping between session id and
				// user id.
				logger.Info("session does not exist",
					zap.String("sid", sid))
				continue
			}
			err := sess.Conn().WriteJSON(msg)
			if err != nil {
				c.Clear(sess)
				logger.Warn("error sending json", zap.Error(err))
				return err
			}
		}
	}
	return nil
}

func (c *Chat) eventloop() {
	logger := c.logger.With(zap.String("method", "eventloop"))

	logger.Info("started event loop")
	getStatus := func(user string) string {
		sessions := c.Get(UserID(user))
		// User has not sessions in place.
		if len(sessions) == 0 {
			return "0"
		}
		return "1"
	}

loop:
	for {
		select {
		case <-c.quit:
			logger.Info("quit")
			break loop
		case msg, ok := <-c.broadcast:
			if !ok {
				break loop
			}
			logger.Info("processing message",
				zap.String("type", msg.Type),
				zap.String("receiver", msg.Receiver),
				zap.String("sender", msg.Sender),
				zap.String("text", msg.Text))

			switch msg.Type {
			case MessageTypeStatus:
				// Requesting the status of a particular user.
				// msg.Text is the user_id in question.
				msg.Text = getStatus(msg.Text)
				logger.Info("check status user in room",
					zap.String("status", msg.Text),
					zap.String("user", msg.Text),
					zap.String("room", msg.Receiver))
			case MessageTypeAuth:
				msg.Text = msg.Sender
			case MessageTypeMessage:
				// Store the conversation in a database. It
				// might be a better idea to use a queue rather
				// than writing directly to the datastore.
				_, err := c.db.CreateConversationReply(msg.Sender, msg.Receiver, msg.Text)
				if err != nil {
					logger.Warn("error creating reply", zap.Error(err))
					continue
				}
			default:
			}
			c.Broadcast(msg)
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
	logger := c.logger.With(zap.String("method", "Bind"))
	// If there are no sessions associated with the user yet, create the
	// rooms and add users into it.
	if sess := c.Get(uid); len(sess) == 0 {
		logger.Info("joining room")
		c.Join(uid)
	}
	logger.Info("bind user with session",
		zap.String("user", uid.String()),
		zap.String("session", sid.String()))

	// Tie the user to the existing session.
	c.lookup.Add(uid.String(), sid.String())

	return func() {
		// Clear the current session that is tied to the user.
		sessions := c.Get(sid)
		for _, sess := range sessions {
			c.Clear(sess)
		}

		logger.Info("clearing session",
			zap.String("user", uid.String()),
			zap.String("session", sid.String()))

		// If the user does not have any other sessions left, clear the room.
		// One user can have multiple sessions.
		if len(c.Get(uid)) == 0 {
			logger.Info("leaving room",
				zap.String("user", uid.String()),
				zap.String("session", sid.String()))
			// Clear rooms. Only if there are no longer any
			// sessions available for that particular user.
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
			Type:     MessageTypeOffline,
			Sender:   sender,
			Receiver: receiver,
			Text:     "1",
		}

		// Add user to room, and keep track of rooms for user.
		err := c.rooms.Add(uid.String(), room.RoomID)
		if err != nil {
			log.Printf("JoinError: %v\n", err)
		}
		log.Printf("Join: added to rooms %#v\n", room.RoomID)
	}
}

func (c *Chat) Leave(uid UserID) {
	// For each room that the user belong to, remove the user.
	onDelete := func(room string) {
		log.Println("deleting room", room, uid)
		sender, receiver := string(uid), room
		c.broadcast <- Message{
			Type:     MessageTypePresence,
			Sender:   sender,
			Receiver: receiver,
			Text:     "0", // Offline
		}
	}
	// Delete user -> rooms relationship.
	err := c.rooms.Delete(uid.String(), onDelete)
	if err != nil {
		log.Println("error removing from roomz", err)
	}
	c.logger.Info("left room",
		zap.String("method", "Leave"),
		zap.String("user", uid.String()))

}

func (c *Chat) Clear(sess *Session) {
	// Closes the connection, and delete the session.
	sess.Conn().Close()
	sessionID := sess.SessionID()
	c.sessions.Delete(string(sessionID))

	c.logger.Info("clear session",
		zap.String("session", sessionID.String()),
		zap.String("method", "Clear"))
}

// One-to-many relationship between sessions and user.
// One user can have multiple sessions (mobile, browser with multiple tabs etc)
// Get(UserID) will return a slice of sessions.
// Get(SessionID) will return one user.
// To get the session.
func (c *Chat) Get(key interface{}) []*Session {
	switch v := key.(type) {
	case SessionID:
		sess := c.sessions.Get(v.String())
		if sess == nil {
			return nil
		}
		return []*Session{sess}
	case UserID:
		// If the UserID is provided, get the SessionID first in order
		// to retrieve the session.
		// A userID can have multiple sessions (many tabs).
		var result []*Session
		sessions := c.lookup.GetSessions(v.String())
		for _, sess := range sessions {
			session := c.sessions.Get(sess)
			if session != nil {
				result = append(result, session)
			}
		}
		return result
	default:
		return nil
	}
}
