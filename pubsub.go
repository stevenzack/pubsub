package pubsub

type Client struct {
	Ch     chan interface{}
	ChanId string
}
type ClientMap map[Client]bool

type Msg struct {
	ChanId string
	Data   interface{}
}

var (
	entering  = make(chan Client, 1)
	leaving   = make(chan Client, 1)
	messaging = make(chan Msg, 1)
	shutdown  = make(chan bool)
)

func init() {
	go broadcaster()
}
func broadcaster() {
	chanIdMap := make(map[string]ClientMap)
	for {
		select {
		case msg := <-messaging:
			if climap, ok := chanIdMap[msg.ChanId]; ok && climap != nil {
				for cli := range climap {
					cli.Ch <- msg.Data
				}
			}
		case cli := <-entering:
			if climap, ok := chanIdMap[cli.ChanId]; !ok || climap == nil {
				chanIdMap[cli.ChanId] = make(map[Client]bool)
			}
			chanIdMap[cli.ChanId][cli] = true
		case cli := <-leaving:
			if climap, ok := chanIdMap[cli.ChanId]; ok && climap != nil {
				delete(chanIdMap[cli.ChanId], cli)
				close(cli.Ch)
			}
		case <-shutdown:
			for _, v := range chanIdMap {
				if v != nil {
					for c := range v {
						close(c.Ch)
					}
				}
			}
			break
		}
	}
}

type PubSub struct {
	c Client
}

func NewPubSub() *PubSub {
	return &PubSub{c: Client{Ch: make(chan interface{}, 1)}}
}
func (ps *PubSub) Sub(chanId string, f func(interface{})) {
	ps.c.ChanId = chanId
	entering <- ps.c
	for data := range ps.c.Ch {
		f(data)
	}
}
func (ps *PubSub) UnSub() {
	leaving <- ps.c
}
func Pub(chanId string, data interface{}) {
	messaging <- Msg{ChanId: chanId, Data: data}
}
func Shutdown() {
	shutdown <- true
}
