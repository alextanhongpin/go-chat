package entity

import "time"

type Post struct {
	ID        int
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	UserID    int
}
