package chat

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

// PubSub represents the redis pub/sub
type PubSub struct {
	Pool *redis.Pool
}

// Conn returns a reused connection from the pool
func (ps *PubSub) Conn() redis.Conn {
	return ps.Pool.Get()
}

// Publish a message to a channel
func (ps *PubSub) Publish(room string, msg Message) error {
	c := ps.Conn()

	out, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err := c.Do("PUBLISH", room, string(out)); err != nil {
		return err
	}

	return nil
}

// Subscribe to a redis channel
func (ps *PubSub) Subscribe(ctx context.Context, room string, subscription *Subscription) {
	if room == "" {
		return
	}

	for {
		c := ps.Conn()
		psc := redis.PubSubConn{Conn: c}
		psc.Subscribe(room)

		defer func() {
			c.Close()
			psc.Close()
			psc.Unsubscribe(room)
		}()

		for c.Err() == nil {
			switch v := psc.Receive().(type) {
			case redis.Message:
				log.Println("got message redis:", v.Channel, string(v.Data))
				var msg Message
				if err := json.Unmarshal(v.Data, &msg); err != nil {
					log.Printf("error unmarshalling redis published data: %s\n", err.Error())
					continue
				}
				// Write to websocket
				log.Println("writing message to websocket", msg)
				subscription.Broadcast(msg)
			case redis.Subscription:
				log.Printf("message is %#v %s %s %d", v, v.Channel, v.Kind, v.Count)
			case error:
				log.Println("error pub/sub. delivery has stopped")
				return
			}
		}
		log.Println("unsubscribing from redis now")
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

// NewPool returns a redis pool that allows a connection to be reused
func NewPool(port string) *redis.Pool {
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
func NewPubSub(port string) *PubSub {
	pool := NewPool(port)

	conn := pool.Get()
	v, err := conn.Do("PING")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("PING redis, got: %v\n", v)

	return &PubSub{
		Pool: NewPool(port),
	}
}
