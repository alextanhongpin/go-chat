package chat

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

type PubSub struct {
	Pool *redis.Pool
}

func (ps *PubSub) Conn() redis.Conn {
	return ps.Pool.Get()
}

func (ps *PubSub) Publish(room string, msg Message) {
	c := ps.Conn()
	log.Println("publishing to redis", room, msg)
	c.Do("PUBLISH", room, "publish:msg")
}

func (ps *PubSub) Subscribe(room string, subscription *Subscription) {
	log.Println("subscribing to redis:", room)
	if room == "" {
		return
	}
	for {
		c := ps.Conn()
		psc := redis.PubSubConn{Conn: c}

		// Subscribe to the room
		psc.Subscribe(room)

		// Listen to the room for new messages
		for c.Err() == nil {
			switch v := psc.Receive().(type) {
			case redis.Message:
				log.Println("got message", v.Channel, string(v.Data))
				// Write to websocket
				subscription.Broadcast(Message{
					Text:   "this is from redis!",
					Handle: "redis",
					Room:   "tech",
				})
			case redis.Subscription:
				log.Printf("message is %#v %s %s %d", v, v.Channel, v.Kind, v.Count)
			case error:
				log.Println("error pub/sub. delivery has stopped")
				return
			}
		}
		log.Println("unsubscribing from redis now")
		c.Close()
		psc.Close()
	}
}

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
	}
}

func NewPubSub(port string) *PubSub {
	return &PubSub{
		Pool: NewPool(port),
	}
}
