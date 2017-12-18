package chat

import (
	"log"
)

// Room holds the clients in a particular room
type Room struct {
	// Registered clients
	Clients map[string]map[*Client]bool

	// Inbound messages from the clients
	Broadcast chan Message

	// Register requests from the clients
	Subscribe chan *Subscription

	// Unrequest requests from the clients
	Unsubscribe chan *Subscription
}

// NewRoom returns a reference to a room
func NewRoom() *Room {
	return &Room{
		Broadcast:   make(chan Message),
		Subscribe:   make(chan *Subscription),
		Unsubscribe: make(chan *Subscription),
		Clients:     make(map[string]map[*Client]bool),
	}
}

// Join adds a new client to a room
func (r *Room) Join(room string, client *Client) {
	clients := r.Clients[room]
	if clients == nil {
		clients := make(map[*Client]bool)
		r.Clients[room] = clients
	}
	r.Clients[room][client] = true
}

// Quit removes a client from a room
func (r *Room) Quit(room string, client *Client) {
	clients := r.Clients[room]
	if clients != nil {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.Send)
			if len(clients) == 0 {
				delete(r.Clients, room)
			}
		}
	}
}

// Emit will broadcast the message to a room
func (r *Room) Emit(msg Message) {
	log.Printf("room.Emit \"%s\" to room \"%s\" from \"%s\"\n", msg.Text, msg.Room, msg.Handle)
	clients := r.Clients[msg.Room]

	for c := range clients {
		select {
		case c.Send <- msg:
			log.Println("send msg:", msg)
		default:
			r.Quit(msg.Room, c)
		}
	}
}

// Run will initialize the room and the corresponding channels
func (r *Room) Run() {
	for {
		select {
		case s := <-r.Subscribe:
			r.Join(s.Room, s.Client)
		case s := <-r.Unsubscribe:
			r.Quit(s.Room, s.Client)
		case m := <-r.Broadcast:
			r.Emit(m)
		default:
		}
	}
}
