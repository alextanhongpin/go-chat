package chat

import (
	"encoding/json"
	"log"

	"github.com/garyburd/redigo/redis"
)

type PubSub struct {
	Pool *redis.Pool
}

func (ps *PubSub) Conn() redis.Conn {
	return ps.Pool.Get()
}

func (ps *PubSub) Publish(room string, msg Message) error {
	c := ps.Conn()
	log.Println("publishing to redis", room, msg)
	// Marshal into string
	out, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	c.Do("PUBLISH", room, string(out))
	return nil
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
				var msg Message
				err := json.Unmarshal(v.Data, &msg)
				if err != nil {
					log.Printf("error unmarshalling redis published data: %s\n", err.Error())
					psc.Unsubscribe(room)
					c.Close()
					psc.Close()
					return
				}
				// Write to websocket
				subscription.Broadcast(msg)
			case redis.Subscription:
				log.Printf("message is %#v %s %s %d", v, v.Channel, v.Kind, v.Count)
			case error:
				log.Println("error pub/sub. delivery has stopped")
				c.Close()
				psc.Close()
				psc.Unsubscribe(room)
				return
			}
		}
		log.Println("unsubscribing from redis now")

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
