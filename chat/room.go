package chat

// Room holds the clients in a particular room
type Room struct {
	// Registered clients
	Clients map[string]map[*Client]bool

	// Inbound messages from the clients
	Broadcast chan Message

	// Register requests from the clients
	Subscribe chan *Client

	// Unrequest requests from the clients
	Unsubscribe chan *Client
}

// NewRoom returns a reference to a room
func NewRoom() *Room {
	return &Room{
		Broadcast:   make(chan Message),
		Subscribe:   make(chan *Client),
		Unsubscribe: make(chan *Client),
		Clients:     make(map[string]map[*Client]bool),
	}
}

// Join adds a new client to a room
func (r *Room) Join(c *Client) {
	clients := r.Clients[c.Room]
	if clients == nil {
		clients := make(map[*Client]bool)
		r.Clients[c.Room] = clients
	}
	r.Clients[c.Room][c] = true
}

// Quit removes a client from a room
func (r *Room) Quit(c *Client) {
	clients := r.Clients[c.Room]
	if clients != nil {
		if _, ok := clients[c]; ok {
			delete(clients, c)
			close(c.Send)
			if len(clients) == 0 {
				delete(r.Clients, c.Room)
			}
		}
	}
}

// Emit will broadcast the message to a room
func (r *Room) Emit(msg Message) {
	clients := r.Clients[msg.Room]

	// Perform business logic for handling different messages here
	for c := range clients {
		select {
		case c.Send <- msg:
		default:
			r.Quit(c)
		}
	}
}

// Run will initialize the room and the corresponding channels
func (r *Room) Run() {
	for {
		select {
		case c := <-r.Subscribe:
			r.Join(c)
		case c := <-r.Unsubscribe:
			r.Quit(c)
		case m := <-r.Broadcast:
			r.Emit(m)
		default:
		}
	}
}
