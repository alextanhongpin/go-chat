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
func (s *Subscription) Read(pubsub *PubSub) {
	c := s.Client
	// Terminate the connection when the server stops
	defer func() {
		c.Unsubscribe(s)
		s.Close()
	}()

	// Limit the maximum size allowed by the peer
	c.Conn.SetReadLimit(maxMessageSize)
	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			log.Println("websocket closed:", err)
			break
		}
		log.Println("got message:", msg)
		// TODO: publish it to redis here
		// c.Broadcast(msg)
		pubsub.Publish(s.Room, msg)
		// pubsub.Do("PUBLISH", s.Room, string(msg))
	}
}

func (s *Subscription) Broadcast(msg Message) {
	s.Client.Broadcast(msg)
}

// Write will write new messages to the client
func (s *Subscription) Write() {
	c := s.Client
	// Remember to close the connection once there is no more writes
	defer func() {
		s.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The Room closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// Write json data
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
