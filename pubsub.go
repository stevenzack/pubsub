package pubsub

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

type PubSub struct {
	subers map[string]chan interface{}
}

func NewPubSub() *PubSub {
	ps := &PubSub{}
	ps.subers = make(map[string]chan interface{})
	return ps
}
func (ps *PubSub) Pub(data interface{}) {
	for _, v := range ps.subers {
		v <- data
	}
}
func (ps *PubSub) Sub(f func(interface{}, string)) {
	c := make(chan interface{}, 1)
	chanId := NewToken()
	ps.subers[chanId] = c
	for a := range c {
		f(a, chanId)
	}
}
func (ps *PubSub) UnSub(chanId string) {
	v, ok := ps.subers[chanId]
	if ok {
		delete(ps.subers, chanId)
		close(v)
		return
	}
}
func NewToken() string {
	ct := time.Now().UnixNano()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(ct, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}
