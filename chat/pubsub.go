package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

// PubSub represents the redis pub/sub
type PubSub struct {
	Pool    *redis.Pool
	Channel string
}

// Conn returns a reused connection from the pool
func (ps *PubSub) Conn() redis.Conn {
	return ps.Pool.Get()
}

// Publish a message to a channel
func (ps *PubSub) Publish(msg Message) error {

	out, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c := ps.Conn()
	defer c.Close()

	// TODO: Make channel a variable
	if _, err := c.Do("PUBLISH", "chat", string(out)); err != nil {
		return errors.Wrap(err, "unable to publish message to Redis")
	}
	if err := c.Flush(); err != nil {
		return errors.Wrap(err, "unable to flush published message to Redis")
	}

	// LPUSH and LTRIM, LRANGE 0 10

	return nil
}

// Subscribe to a redis channel
func (ps *PubSub) Subscribe(room *Room) {
	c := ps.Conn()
	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe("chat")

	for c.Err() == nil {
		switch v := psc.Receive().(type) {
		case redis.Message:
			var msg Message
			if err := json.Unmarshal(v.Data, &msg); err != nil {
				log.Printf("error unmarshalling redis published data: %s\n", err.Error())
				continue
			}
			room.Broadcast <- msg
		case redis.Subscription:
			log.Printf("message is %#v %s %s %d", v, v.Channel, v.Kind, v.Count)
		case error:
			return
		}
	}
	c.Close()
}

// NewPool returns a redis pool that allows a connection to be reused
func NewPool(port, channel string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   5,
		MaxActive: 10,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", port)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// NewPubSub returns a pointer to the PubSub struct
func NewPubSub(port, channel string) *PubSub {
	pool := NewPool(port, channel)
	conn := pool.Get()
	defer conn.Close()

	v, err := conn.Do("PING")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("PING redis, got: %v\n", v)

	return &PubSub{
		Pool: pool,
	}
}
