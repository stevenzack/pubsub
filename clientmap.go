package pubsub

type clientMap map[int64]*Client

func newClientMap() clientMap {
	return make(map[int64]*Client)
}
