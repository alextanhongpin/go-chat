package repository

import "github.com/alextanhongpin/go-chat/entity"

// Room represents the interface for room repository.
type Room interface {
	CreateRoom(users ...string) (int64, error)
	GetRoom(userID string) ([]string, error)
	GetRooms(userID string) ([]entity.UserRoom, error)
}
