package pubsub

import (
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

func (c *Client) Sub(chanId string, listener func(interface{})) {
	c.channelID = chanId
	select {
	case c.server.entering <- c:
	default:
	}
	for msg := range c.receiver {
		listener(msg)
	}
}

func (c *Client) UnSub() {
	select {
	case c.server.leaving <- c:
	default:
	}
}
