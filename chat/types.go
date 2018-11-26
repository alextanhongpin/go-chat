package chat

// UserID represents the user id.
type UserID string

func (u UserID) String() string {
	return string(u)
}

// RoomID represents the room id.
type RoomID string

func (r RoomID) String() string {
	return string(r)
}

// SessionID represents the session id.
type SessionID string

func (s SessionID) String() string {
	return string(s)
}
