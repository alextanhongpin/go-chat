package chat

type UserID string

func (u UserID) String() string {
	return string(u)
}

type RoomID string

func (r RoomID) String() string {
	return string(r)
}

type SessionID string

func (s SessionID) String() string {
	return string(s)
}

type Tabler interface {
	Get(id interface{}) []string
	Add(a, b interface{}) error
	Delete(a interface{}, fn func(b string)) error
}
