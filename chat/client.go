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
	// Conn represents the websocket connection
	Conn *websocket.Conn

	// Buffered channel to hold the messages
	Send chan Message
}

// NewClient returns a new chat server
func NewClient(ws *websocket.Conn) *Client {
	return &Client{
		Conn: ws,
		Send: make(chan Message),
	}
}
