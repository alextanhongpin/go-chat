package chat

import (
	"log"
	"time"

	"github.com/alextanhongpin/go-chat/ticket"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Client contains the Client Connection and Room
type Client struct {
	// Client *Client
	Room string

	// Conn represents the websocket connection
	Conn *websocket.Conn

	// Buffered channel to hold the messages
	Send chan Message
}

// // Close will terminate the client websocket connection
// func (s *Subscription) Close() {
// 	s.Conn.Close()
// }

// Read will proceed to read the messages published by the client
func (c *Client) Read(pubsub *PubSub, room *Room) {
	c.Conn.SetReadLimit(maxMessageSize)
	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			log.Println(errors.Wrap(err, "websocket closed"))
			break
		}

		// Carry out validation for different type of message type here
		// e.g. authentication
		if msg.Token == "" {
			log.Println("error: user is not authenticated")
			break
		}

		_, err := ticket.Verify(msg.Token)
		if err != nil {
			log.Println(errors.Wrap(err, "invalid token"))
			break
		}

		if err := pubsub.Publish(msg); err != nil {
			log.Println(errors.Wrap(err, "error publishing"))
			break
		}
	}

	room.Unsubscribe <- c

	c.Conn.Close()
}

// Write will write new messages to the client
func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// Perform cleaning of data here (hiding data that does not want to be exposed such as token)
			message.Room = ""
			message.Token = ""
			c.Conn.WriteJSON(message)
		}
	}
}

// NewClient will return a new pointer to the subscription
func NewClient(conn *websocket.Conn, room string) *Client {
	return &Client{
		Room: room,
		Send: make(chan Message),
		Conn: conn,
	}
}
