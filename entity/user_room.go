package entity

type UserRoom struct {
	RoomID int    `json:"room_id"`
	UserID int    `json:"user_id"`
	Name   string `json:"name,omitempty"`
}
