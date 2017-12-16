package chat

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Client represents the chat client
type Client struct {
	// Room represent the rooms for each websocket connection
	Room *Room

	// Conn represents the websocket connection
	Conn *websocket.Conn

	// Buffered channel to hold the messages
	Send chan Message
}

// Subscribe adds a subscription to the room
func (c *Client) Subscribe(s *Subscription) {
	c.Room.Register <- s
}

// Unsubscribe removes a subscription from a room
func (c *Client) Unsubscribe(s *Subscription) {
	c.Room.Unregister <- s
}

// NewClient returns a new chat server
func NewClient(ws *websocket.Conn, r *Room) *Client {
	return &Client{
		Conn: ws,
		Room: r,
		Send: make(chan Message),
	}
}