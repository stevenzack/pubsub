package pubsub

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

var (
	box  = make(map[string]map[string]chan string)
	lock sync.Mutex
)

type PubSub struct {
	token string
	mc    chan string
}

func NewPubSub() *PubSub {
	ps := &PubSub{}
	ps.token = NewToken()
	return ps
}
func (ps *PubSub) Sub(chanId string, f func([]byte)) {
	ps.mc = make(chan string, 1)
	lock.Lock()
	if v, ok := box[chanId]; ok {
		if v == nil {
			box[chanId] = make(map[string]chan string)
		}
	} else {
		box[chanId] = make(map[string]chan string)
	}
	box[chanId][ps.token] = ps.mc
	lock.Unlock()
	for str := range ps.mc {
		f([]byte(str))
	}
}
func (ps *PubSub) UnSub(chanId string) {
	if v, ok := box[chanId]; ok && v != nil {
		lock.Lock()
		delete(box[chanId], ps.token)
		lock.Unlock()
		close(ps.mc)
	}
}
func NewToken() string {
	ct := time.Now().UnixNano()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(ct, 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}
