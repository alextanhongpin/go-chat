package entity

// UserRoom represents the mapping between the user and the room.
type UserRoom struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
	// The room name. Defaults to the other user's name.
	Name string `json:"name,omitempty"`
}
