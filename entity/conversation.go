package entity

import "time"

// Conversation represents the user conversation in the db.
type Conversation struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
}
