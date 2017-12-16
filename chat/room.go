package chat

import "log"

type Room struct {
	// Registered clients
	Clients map[string]map[*Client]bool

	// Inbound messages from the clients
	Broadcast chan Message

	// Register requests from the clients
	Register chan *Subscription

	// Unrequest requests from the clients
	Unregister chan *Subscription
}

func NewRoom() *Room {
	return &Room{
		Broadcast:  make(chan Message),
		Register:   make(chan *Subscription),
		Unregister: make(chan *Subscription),
		Clients:    make(map[string]map[*Client]bool),
	}
}

func (r *Room) Run() {
	for {
		select {
		case s := <-r.Register:
			connections := r.Clients[s.Room]
			if connections == nil {
				connections = make(map[*Client]bool)
				r.Clients[s.Room] = connections
			}
			r.Clients[s.Room][s.Client] = true

		case s := <-r.Unregister:
			connections := r.Clients[s.Room]
			if connections != nil {
				if _, ok := connections[s.Client]; ok {
					delete(connections, s.Client)
					close(s.Client.Send)
					if len(connections) == 0 {
						delete(r.Clients, s.Room)
					}
				}
			}

		case m := <-r.Broadcast:
			log.Printf("len of clients: %d\n room: %s", len(r.Clients), m.Room)
			connections := r.Clients[m.Room]
			log.Println("number of connections in room:", len(connections))
			for c := range connections {
				select {
				case c.Send <- m:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(r.Clients, m.Room)
					}
				}
			}
			// for client := range r.Clients {
			// 	log.Println("message to broadcast:", message)
			// 	select {
			// 	case client.Send <- message:
			// 	default:
			// 		close(client.Send)
			// 		delete(r.Clients, client)
			// 	}
			// }
		}
	}
}
