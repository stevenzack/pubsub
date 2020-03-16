package pubsub

import (
	"errors"
	"fmt"
	"sync"
)

type Server struct {
	log          bool
	channelMap   sync.Map //string to *ClientMap
	entering     chan *Client
	leaving      chan *Client
	messaging    chan *message
	shutdown     chan bool
	maxClientNum int

	isRunning bool
}

func (s *Server) SetMaxClientNum(i int) {
	s.maxClientNum = i
}

func NewServer() *Server {
	return &Server{
		entering:  make(chan *Client, 24),
		leaving:   make(chan *Client, 24),
		messaging: make(chan *message, 24),
		shutdown:  make(chan bool, 24),
	}
}

func StartNewServer() *Server {
	s := NewServer()
	s.Start()
	return s
}

func (s *Server) SetLog(b bool) {
	s.log = b
}

func (s *Server) load(chanID string) (*clientMap, bool) {
	cmv, ok := s.channelMap.Load(chanID)
	if !ok {
		return nil, false
	}
	return cmv.(*clientMap), true
}

func (s *Server) rangeChannelMap(fn func(string, *clientMap) bool) {
	s.channelMap.Range(func(k, v interface{}) bool {
		return fn(k.(string), v.(*clientMap))
	})
}

func (s *Server) loadOrCreate(chanID string) *clientMap {
	cmv, _ := s.channelMap.LoadOrStore(chanID, newClientMap())
	return cmv.(*clientMap)
}

func (s *Server) Run() {
	s.isRunning = true
	for s.isRunning {
		select {
		case cli := <-s.entering:
			if s.log {
				fmt.Println("entering:")
			}

			cm := s.loadOrCreate(cli.channelID)
			if s.maxClientNum > 0 {
				if cm.Len() >= s.maxClientNum {
					pop := s.pop(cli.channelID)
					if pop != nil {
						pop.closeChannel()
					}
				}
			}
			cm.Store(cli.id, cli)
		case msg := <-s.messaging:
			if s.log {
				fmt.Println("messaging:", msg.data)
			}
			if cm, ok := s.load(msg.channelID); ok && cm != nil {
				cm.Range(func(id int64, cli *Client) bool {
					select {
					case cli.receiver <- msg.data:
					default:
						fmt.Println("pubsub error: client not running")
					}
					return true
				})
			}
		case cli := <-s.leaving:
			if s.log {
				fmt.Println("leaving:")
			}
			if cm, ok := s.load(cli.channelID); ok && cm != nil {
				cm.Delete(cli.id)
				cli.closeChannel()
			}
		case b := <-s.shutdown:
			if b {
				if s.log {
					fmt.Println("shutdown:")
				}
				s.rangeChannelMap(func(_ string, cm *clientMap) bool {
					if cm != nil {
						cm.Range(func(_ int64, cli *Client) bool {
							cli.closeChannel()
							return true
						})
					}
					return true
				})
				s.isRunning = false
			}
		}
	}
}

func (s *Server) Start() {
	go s.Run()
}

func (s *Server) Pub(chanID string, data interface{}) {
	s.messaging <- &message{channelID: chanID, data: data}
}

func (s *Server) Stop() error {
	select {
	case s.shutdown <- true:
		return nil
	default:
		return errors.New("not running")
	}
}

func (s *Server) pop(chanID string) *Client {
	cm, ok := s.load(chanID)
	if !ok {
		return nil
	}
	var cli *Client
	cm.Range(func(id int64, c *Client) bool {
		cli = c
		cm.Delete(id)
		return false
	})
	return cli
}
