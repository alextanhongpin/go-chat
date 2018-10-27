package repository

import "time"

type UserCache interface {
	SetUser(key, value string, duration time.Duration) error
	HasUser(key string) bool
	AddUser(user, room string) error
	GetUsers(room string) []string
	GetRooms(user string) []string
	RemoveUser(user string) error

	Set(key, value string) error
	Get(key string) (string, error)
}
