package entity

import "time"

type Conversation struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
}
