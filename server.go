package pubsub

import (
	"errors"
	"fmt"
)

type Server struct {
	log          bool
	channelMap   map[string]clientMap
	entering     chan *Client
	leaving      chan *Client
	messaging    chan *message
	shutdown     chan bool
	maxClientNum int
}

func (s *Server) SetMaxClientNum(i int) {
	s.maxClientNum = i
}

func NewServer() *Server {
	return &Server{
		channelMap: make(map[string]clientMap),
		entering:   make(chan *Client, 24),
		leaving:    make(chan *Client, 24),
		messaging:  make(chan *message, 24),
		shutdown:   make(chan bool, 24),
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

func (s *Server) Run() {
	for {
		select {
		case cli := <-s.entering:
			if s.log {
				fmt.Println("entering:")
			}
			if cm, ok := s.channelMap[cli.channelID]; !ok || cm == nil {
				s.channelMap[cli.channelID] = newClientMap()
			}
			if s.maxClientNum > 0 {
				if len(s.channelMap[cli.channelID]) >= s.maxClientNum {
					cli.closeChannel()
					continue
				}
			}
			s.channelMap[cli.channelID][cli.id] = cli
		case msg := <-s.messaging:
			if s.log {
				fmt.Println("messaging:", msg.data)
			}
			if clientmap, ok := s.channelMap[msg.channelID]; ok && clientmap != nil {
				for _, cli := range clientmap {
					select {
					case cli.receiver <- msg.data:
					default:
						fmt.Println("pubsub error: client not running")
					}
				}
			}
		case cli := <-s.leaving:
			if s.log {
				fmt.Println("leaving:")
			}
			if cm, ok := s.channelMap[cli.channelID]; ok && cm != nil {
				delete(s.channelMap[cli.channelID], cli.id)
				cli.closeChannel()
			}
		case b := <-s.shutdown:
			if b {
				if s.log {
					fmt.Println("shutdown:")
				}
				for _, cm := range s.channelMap {
					if cm != nil {
						for _, cli := range cm {
							cli.closeChannel()
						}
					}
				}
				break
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

func (s *Server) GetClientNum(chanID string) int {
	return len(s.channelMap[chanID])
}
