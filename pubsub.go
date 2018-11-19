package pubsub

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

type PubSub struct {
	subers map[string]chan string
}

func NewPubSub() *PubSub {
	ps := &PubSub{}
	ps.subers = make(map[string]chan string)
	return ps
}
func (ps *PubSub) PubAll(data []byte) {
	for _, v := range ps.subers {
		v <- string(data)
	}
}
func (p *PubSub) Pub(chanId string, data []byte) bool {
	if suber, ok := p.subers[chanId]; ok {
		suber <- string(data)
		return true
	}
	return false
}
func (ps *PubSub) Sub(f func([]byte), chanId string) {
	c := make(chan string, 1)
	var mchanId = chanId
	if mchanId == "" {
		mchanId = NewToken()
	}
	ps.subers[mchanId] = c
	for a := range c {
		f([]byte(a))
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
func (ps *PubSub) Close() {
	for k, v := range ps.subers {
		close(v)
		delete(ps.subers, k)
	}
}
func NewToken() string {
	ct := time.Now().UnixNano()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(ct, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}
