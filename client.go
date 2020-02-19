package pubsub

import (
	"errors"

	"github.com/StevenZack/tools/strToolkit"
)

type Client struct {
	channelID string
	server    *Server
	id        int64
	receiver  chan interface{}
}

var autoIncrementId = &strToolkit.AutoIncrementID{}

func NewClient(s *Server) *Client {
	return &Client{
		server:   s,
		id:       autoIncrementId.Generate(),
		receiver: make(chan interface{}),
	}
}

func (c *Client) Sub(chanId string, listener func(interface{})) error {
	c.channelID = chanId
	select {
	case c.server.entering <- c:
	default:
		return errors.New("server not running")
	}
	for msg := range c.receiver {
		listener(msg)
	}
	return nil
}

func (c *Client) UnSub() error {
	select {
	case c.server.leaving <- c:
		return nil
	default:
		return errors.New("server not running")
	}
}
