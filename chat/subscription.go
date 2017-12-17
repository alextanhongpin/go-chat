package chat

import (
	"context"
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
func (s *Subscription) Read(ctx context.Context, pubsub *PubSub) {
	c := s.Client
	// Terminate the connection when the server stops

	// Limit the maximum size allowed by the peer
	c.Conn.SetReadLimit(maxMessageSize)
	for {
		// defer func() {

		// }()
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			log.Println("subscription error: websocket closed due to", err)
			break
		}
		log.Println("got message:", msg)
		// TODO: publish it to redis here
		// c.Broadcast(msg)
		if err := pubsub.Publish(s.Room, msg); err != nil {
			log.Println("error publishing:", err.Error())
			break
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
	c.Unsubscribe(s)
	s.Close()
}

// Broadcast will send the message to the subscribed clients
func (s *Subscription) Broadcast(msg Message) {
	s.Client.Broadcast(msg)
}

// Write will write new messages to the client
func (s *Subscription) Write(ctx context.Context) {
	c := s.Client
	// Remember to close the connection once there is no more writes
	// defer s.Close()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Println("closing subscription channel")
				// The Room closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// Write json data
			c.Conn.WriteJSON(message)
		case <-ctx.Done():
			return
			// default:
		}
	}

	// s.Close()
}

// NewSubscription will return a new pointer to the subscription
func NewSubscription(room string, client *Client) *Subscription {
	return &Subscription{
		Client: client,
		Room:   room,
	}
}
