package chat

import (
	"fmt"
	"log"

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
	pipe.SAdd(userKey(user), room)
	pipe.SAdd(roomKey(room), user)
	log.Println("added", userKey(user), roomKey(room))
	_, err := pipe.Exec()
	return err
}

func (t *TableCache) Delete(user string, fn func(string)) error {
	log.Println("called delete pipeline", user)
	pipe := t.client.Pipeline()
	rooms := t.client.SMembers(userKey(user)).Val()
	// Get the rooms that the user belong to.
	// rooms := pipe.SMembers(userKey(user)).Val()
	log.Println("rooms belonging to", rooms, user)

	// For each room, remove the user from the set.
	for _, room := range rooms {
		pipe.SRem(roomKey(room), user)
		log.Println("table cache", room)
		fn(room)
	}
	_, err := pipe.Exec()
	return err
}

func (t *TableCache) GetRooms(user string) []string {
	return t.client.SMembers(userKey(user)).Val()
}

func (t *TableCache) GetUsers(room string) []string {
	return t.client.SMembers(roomKey(room)).Val()
}

func userKey(u string) string {
	return fmt.Sprintf("user:%s", u)
}

func roomKey(r string) string {
	return fmt.Sprintf("room:%s", r)
}
