package chat

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Cache struct {
	client *redis.Client
}

func NewCache() *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &Cache{client: client}
}

func (c *Cache) Set(key, value string) error {
	return c.client.Set(key, value, 0).Err()
}

func (c *Cache) Get(key string) (string, error) {
	return c.client.Get(key).Result()
}

func (c *Cache) SetUser(key, value string, duration time.Duration) error {
	_, err := c.client.SetNX(key, value, duration).Result()
	return err
}

func (c *Cache) HasUser(key string) bool {
	return c.client.Get(key).Err() == redis.Nil
}

func (c *Cache) AddUser(user, room string) error {
	pipe := c.client.Pipeline()

	pipe.SAdd(c.roomKey(room), user)
	pipe.SAdd(c.userKey(user), room)

	_, err := pipe.Exec()
	return err
}

func (c *Cache) GetUsers(room string) []string {
	return c.client.SMembers(c.roomKey(room)).Val()
}

func (c *Cache) GetRooms(user string) []string {
	return c.client.SMembers(c.userKey(user)).Val()
}

func (c *Cache) RemoveUser(user string) error {
	pipe := c.client.Pipeline()
	rooms := pipe.SMembers(c.userKey(user)).Val()
	for _, room := range rooms {
		pipe.SRem(c.roomKey(room), user)
	}
	_, err := pipe.Exec()
	return err
}

func (c *Cache) roomKey(room string) string {
	return fmt.Sprintf("room:%s", room)
}

func (c *Cache) userKey(user string) string {
	return fmt.Sprintf("user:%s", user)
}
