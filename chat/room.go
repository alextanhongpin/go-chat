package chat

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
func (r *Room) Join(s *Subscription) {
	clients := r.Clients[s.Room]
	if clients == nil {
		clients := make(map[*Client]bool)
		r.Clients[s.Room] = clients
	}
	r.Clients[s.Room][s.Client] = true
}

// Quit removes a client from a room
func (r *Room) Quit(s *Subscription) {
	clients := r.Clients[s.Room]
	if clients != nil {
		if _, ok := clients[s.Client]; ok {
			delete(clients, s.Client)
			close(s.Client.Send)
			if len(clients) == 0 {
				delete(r.Clients, s.Room)
			}
		}
	}
}

// Emit will broadcast the message to a room
func (r *Room) Emit(msg Message) {
	clients := r.Clients[msg.Room]

	for c := range clients {
		select {
		case c.Send <- msg:
		default:
			r.Quit(&Subscription{Client: c, Room: msg.Room})
		}
	}
}

// Run will initialize the room and the corresponding channels
func (r *Room) Run() {
	for {
		select {
		case s := <-r.Subscribe:
			r.Join(s)
		case s := <-r.Unsubscribe:
			r.Quit(s)
		case m := <-r.Broadcast:
			r.Emit(m)
		default:
		}
	}
}
