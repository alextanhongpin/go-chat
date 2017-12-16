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

	// Holds the context to the websocket connection
	ws *websocket.Conn

	Send chan Message
}

// NewClient returns a new chat client
func NewClient(ws *websocket.Conn, r *Room) *Client {
	return &Client{
		ws:   ws,
		Room: r,
		Send: make(chan Message),
	}
}
