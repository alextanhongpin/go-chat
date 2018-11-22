package chat

type UserID string
type RoomID string
type SessionID string

type Tabler interface {
	Get(id interface{}) []string
	Add(a, b interface{}) error
	Delete(a interface{}, fn func(b string)) error
}
