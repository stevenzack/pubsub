package pubsub

import "sync"

type clientMap struct {
	m sync.Map // int64 to *Client
}

func newClientMap() *clientMap {
	return &clientMap{}
}

func (c *clientMap) Len() int {
	l := 0
	c.m.Range(func(k, v interface{}) bool {
		l++
		return true
	})
	return l
}

func (c *clientMap) Store(id int64, cli *Client) {
	c.m.Store(id, cli)
}

func (c *clientMap) Range(fn func(int64, *Client) bool) {
	c.m.Range(func(k, v interface{}) bool {
		return fn(k.(int64), v.(*Client))
	})
}

func (c *clientMap) Delete(id int64) {
	c.m.Delete(id)
}
