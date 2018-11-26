package chat

import (
	"fmt"

	"github.com/go-redis/redis"
)

// TableCache represents a table with distributed cache.
type TableCache struct {
	client *redis.Client
}

// NewTableCache returns a new TableCache.
func NewTableCache(client *redis.Client) *TableCache {
	return &TableCache{client}
}

// Add adds both user and room to the cache.
func (t *TableCache) Add(user, room string) error {
	pipe := t.client.Pipeline()
	// Add room to user's list.
	pipe.SAdd(userKey(user), room)
	// Add user to the room.
	pipe.SAdd(roomKey(room), user)
	_, err := pipe.Exec()
	return err
}

// Delete removes the user from the room, and clears the user's room list.
func (t *TableCache) Delete(user string, fn func(string)) error {
	// Get the rooms the user is in.
	rooms := t.client.SMembers(userKey(user)).Val()

	pipe := t.client.Pipeline()
	// For each room, remove the user from the set.
	for _, room := range rooms {
		pipe.SRem(roomKey(room), user)
		fn(room)
	}
	_, err := pipe.Exec()
	return err
}

// GetRooms returns the rooms the user is in.
func (t *TableCache) GetRooms(user string) []string {
	return t.client.SMembers(userKey(user)).Val()
}

// GetUsers returns the users in a room.
func (t *TableCache) GetUsers(room string) []string {
	return t.client.SMembers(roomKey(room)).Val()
}

func userKey(u string) string {
	return fmt.Sprintf("user:%s", u)
}

func roomKey(r string) string {
	return fmt.Sprintf("room:%s", r)
}
