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

func (r *TableCache) Add(user, room interface{}) error {
	if _, ok := user.(UserID); !ok {
		return errors.New("UserID is required")
	}
	if _, ok := user.(RoomID); !ok {
		return errors.New("RoomID is required")
	}

	pipe := r.client.Pipeline()
	pipe.SAdd(userKey(user.(string)), room)
	pipe.SAdd(roomKey(room.(string)), user)
	_, err := pipe.Exec()
	return err
}

func (r *TableCache) Delete(user interface{}, fn func(string)) error {
	if _, ok := user.(UserID); !ok {
		return errors.New("UserID is required")
	}
	pipe := r.client.Pipeline()
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

func (r *TableCache) Get(a interface{}) []string {
	switch a.(type) {
	case UserID:
		user := a.(string)
		return r.client.SMembers(userKey(user)).Val()
	case RoomID:
		room := a.(string)
		return r.client.SMembers(roomKey(room)).Val()
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
