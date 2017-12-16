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

func (s *Subscription) Read() {
	c := s.Client
	// Terminate the connection when the server stops
	defer func() {
		c.Room.Unregister <- s
		c.Conn.Close()
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
		c.Room.Broadcast <- msg
	}
}

func (s *Subscription) Write() {
	c := s.Client
	// Remember to close the connection once there is no more writes
	defer func() {
		c.Conn.Close()
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
