package pubsub

import (
	"context"
	"sync"
)

type Server struct {
	topics map[string]*Topic
	lock   sync.Mutex
}

func NewServer() *Server {
	return &Server{
		topics: make(map[string]*Topic),
	}
}

func (s *Server) SubscribeTopics(topics []string, lifecycleGoroutine func(), onSignal func(topic string)) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		lifecycleGoroutine()
		cancel()
		s.BroadcastMultiple(topics)
	}()

	for _, topicName := range topics {
		topic := s.loadTopic(topicName)
		name := topicName
		go topic.SubscribeCtx(ctx, func() {
			onSignal(name)
		})
	}

	<-ctx.Done()
}

func (s *Server) loadTopic(name string) *Topic {
	s.lock.Lock()
	topic, ok := s.topics[name]
	if !ok {
		topic = NewTopic()
		s.topics[name] = topic
	}
	s.lock.Unlock()
	return topic
}

func (s *Server) BroadcastMultiple(topics []string) {
	for _, topicName := range topics {
		if topic, ok := s.topics[topicName]; ok {
			topic.Broadcast()
		}
	}
}

func (s *Server) Broadcast(topic string) {
	if t, ok := s.topics[topic]; ok {
		t.Broadcast()
	}
}
