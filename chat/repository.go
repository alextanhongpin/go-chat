package chat

import (
	"fmt"

	"github.com/go-redis/redis"
)

// type Repository interface {
//         Add(user, room string) error
//         Remove(user string)
//         GetUsers(RoomKey) []string
//         GetRooms(UserKey) []string
// }

type Repository struct {
	client *redis.Client
}

func NewRepository() *Repository {
	return nil
}

func (r *Repository) Add(user, room string) error {
	pipe := r.client.Pipeline()
	pipe.SAdd(userKey(user), room)
	pipe.SAdd(roomKey(room), user)
	_, err := pipe.Exec()
	return err
}

func (r *Repository) Remove(user string) error {
	pipe := r.client.Pipeline()
	// Get the rooms that the user belong to.
	rooms := pipe.SMembers(userKey(user)).Val()

	// For each room, remove the user from the set.
	for _, room := range rooms {
		pipe.SRem(roomKey(room), user)
	}
	_, err := pipe.Exec()
	return err
}

func (r *Repository) GetRooms(user string) []string {
	return r.client.SMembers(userKey(user)).Val()
}

func (r *Repository) GetUsers(room string) []string {
	return r.client.SMembers(roomKey(room)).Val()
}

func userKey(u string) string {
	return fmt.Sprintf("user:%s", u)
}

func roomKey(r string) string {
	return fmt.Sprintf("room:%s", r)
}
