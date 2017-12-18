package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Subscription contains the Client Connection and Room
type Subscription struct {
	Client *Client
	Room   string
}

// Close will terminate the client websocket connection
func (s *Subscription) Close() {
	s.Client.Conn.Close()
}

// Read will proceed to read the messages published by the client
func (s *Subscription) Read(pubsub *PubSub, room *Room) {
	log.Println("reading from subscription")
	c := s.Client

	c.Conn.SetReadLimit(maxMessageSize)
	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			log.Println("subscription error: websocket closed due to", err)
			break
		}

		log.Printf("publishing message to %s: %v\n", s.Room, msg)
		if err := pubsub.Publish(msg); err != nil {
			log.Println("error publishing:", err.Error())
			break
		}
	}

	room.Unsubscribe <- s
	s.Close()
}

// Write will write new messages to the client
func (s *Subscription) Write() {
	c := s.Client
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteJSON(message)
		}
	}
}

// NewSubscription will return a new pointer to the subscription
func NewSubscription(room string, client *Client) *Subscription {
	return &Subscription{
		Client: client,
		Room:   room,
	}
}
