package repository

import "github.com/alextanhongpin/go-chat/entity"

type Room interface {
	CreateRoom(users ...string) error
	GetRoom(userID string) ([]int64, error)
	GetRooms(userID string) ([]entity.UserRoom, error)
}
