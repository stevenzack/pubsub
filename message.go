package pubsub

type message struct {
	channelID string
	data      interface{}
}
