package chat

type RoomManager interface {
	Add(user, room string)
	Del(user string)
	GetUsers(room string) []string
	GetRooms(user string) []string
}

type R string
type U string

type RoomManagerImpl struct {
	store map[interface{}]map[interface{}]struct{}
}

func NewRoomManager() *RoomManagerImpl {
	return &RoomManagerImpl{
		store: make(map[interface{}]map[interface{}]struct{}),
	}
}

func (r *RoomManagerImpl) Add(user, room string) {
	r.add(R(room), U(user))
	r.add(U(user), R(room))
}

func (r *RoomManagerImpl) add(a, b interface{}) {
	if _, found := r.store[a]; !found {
		r.store[a] = make(map[interface{}]struct{})
	}
	r.store[a][b] = struct{}{}
}

func (r *RoomManagerImpl) Del(user string) {
	u := U(user)
	for room := range r.store[u] {
		delete(r.store[room], u)
		if len(r.store[room]) == 0 {
			delete(r.store, room)
		}
	}
	delete(r.store, u)
}

func (r *RoomManagerImpl) GetUsers(room string) []string {
	var result []string
	users, found := r.store[R(room)]
	if !found {
		return result
	}
	for u := range users {
		result = append(result, string(u.(U)))
	}
	return result
}

func (r *RoomManagerImpl) GetRooms(user string) []string {
	var result []string
	rooms, found := r.store[U(user)]
	if !found {
		return result
	}
	for v := range rooms {
		result = append(result, string(v.(R)))
	}
	return result
}
