package pubsub

type Server struct {
	channelMap map[string]clientMap
	entering   chan *Client
	leaving    chan *Client
	messaging  chan *message
	shutdown   chan bool
}

func NewServer() *Server {
	return &Server{
		channelMap: make(map[string]clientMap),
		entering:   make(chan *Client),
		leaving:    make(chan *Client),
		messaging:  make(chan *message),
		shutdown:   make(chan bool),
	}
}

func (s *Server) Run() {
	for {
		select {
		case cli := <-s.entering:
			if cm, ok := s.channelMap[cli.channelID]; !ok || cm == nil {
				s.channelMap[cli.channelID] = newClientMap()
			}
			s.channelMap[cli.channelID][cli.id] = cli
		case msg := <-s.messaging:
			if clientmap, ok := s.channelMap[msg.channelID]; ok && clientmap != nil {
				for _, cli := range clientmap {
					select {
					case cli.receiver <- msg.data:
					default:
					}
				}
			}
		case cli := <-s.leaving:
			if cm, ok := s.channelMap[cli.channelID]; ok && cm != nil {
				delete(s.channelMap[cli.channelID], cli.id)
				close(cli.receiver)
			}
		case <-s.shutdown:
			for _, cm := range s.channelMap {
				if cm != nil {
					for _, cli := range cm {
						close(cli.receiver)
					}
				}
			}
			break
		}
	}
}

func (s *Server) Start() {
	go s.Run()
}

func (s *Server) Pub(chanID string, data interface{}) {
	select {
	case s.messaging <- &message{channelID: chanID, data: data}:
	default:
	}
}

func (s *Server) Stop() {
	s.shutdown <- true
}
