package chat

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

type TableCache struct {
	client *redis.Client
}

func NewTableCache(client *redis.Client) *TableCache {
	return &TableCache{client}
}

func (t *TableCache) Add(user, room interface{}) error {
	u, ok := user.(UserID)
	if !ok {
		return errors.New("UserID is required")
	}
	r, ok := room.(RoomID)
	if !ok {
		return errors.New("RoomID is required")
	}

	pipe := t.client.Pipeline()
	pipe.SAdd(userKey(u.String()), r.String())
	pipe.SAdd(roomKey(r.String()), u.String())
	_, err := pipe.Exec()
	return err
}

func (t *TableCache) Delete(user interface{}, fn func(string)) error {
	if _, ok := user.(UserID); !ok {
		return errors.New("UserID is required")
	}
	pipe := t.client.Pipeline()
	// Get the rooms that the user belong to.
	rooms := pipe.SMembers(userKey(user.(string))).Val()

	// For each room, remove the user from the set.
	for _, room := range rooms {
		pipe.SRem(roomKey(room), user)
		fn(room)
	}
	_, err := pipe.Exec()
	return err
}

func (t *TableCache) Get(a interface{}) []string {
	switch a.(type) {
	case UserID:
		user := a.(UserID)
		return t.client.SMembers(userKey(user.String())).Val()
	case RoomID:
		room := a.(RoomID)
		return t.client.SMembers(roomKey(room.String())).Val()
	default:
		return []string{}
	}
}

func userKey(u string) string {
	return fmt.Sprintf("user:%s", u)
}

func roomKey(r string) string {
	return fmt.Sprintf("room:%s", r)
}
