package repository

import "github.com/alextanhongpin/go-chat/entity"

// User represents the DAO for user in the database.
type User interface {
	GetUser(id string) (entity.User, error)
	GetUserByName(name string) (entity.User, error)
}
