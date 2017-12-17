package cache

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

type Client struct {
}

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) Publish(ch, val string) {
	c.Conn.Send(ch, val)
	c.Conn.Flush()
}

func (c *Client) Publish(id string, content interface{}) {
	c.Conn.Do("PUBLISH", id, string(content))
}

func (c *Client) Subscribe(id string) {

	psc = redis.PubSubConn{Conn: c.Conn}
	defer psc.Close()
	if err := psc.Subscribe(id); err != nil {
		log.Println(err)
	}

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			// Get the websocket message here, do the necessary business logic
			// and publish it back to the channel
			// websocket.write
			log.Println("got message", v.Channel, string(v.Data))
		case redis.Subscription:
			log.Printf("message is %#v %s %s %d", v, v.Channel, v.Kind, v.Count)
		case error:
			log.Println("error pub/sub. delivery has stopped")
			return
		}
	}
}

func (c *Client) Pool(addr string) *redis.Pool {
	if addr == "" {
		addr = ":6379"
	}

	return &redis.Pool{
		MaxIdle:   5,
		MaxActive: 10,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

func New() *Client {
	c := Client{}
	pool := c.Pool(":6379")
	defer pool.Close()

	for {
		c := pool.Get()
		psc := redis.PubSubConn{Conn: c}
		psc.Subscribe("example")
		for c.Err() == nil {
			switch v := psc.Receive().(type) {
			case redis.Message:
				log.Println("got message", v.Channel, string(v.Data))
			case redis.Subscription:
				log.Printf("message is %#v %s %s %d", v, v.Channel, v.Kind, v.Count)
			case error:
				log.Println("error pub/sub. delivery has stopped")
				return
			}
		}
		c.Close()
	}

}
