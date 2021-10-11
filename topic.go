package pubsub

import (
	"context"
	"sync"
)

type Topic struct {
	cond *sync.Cond
}

func NewTopic() *Topic {
	return &Topic{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (t *Topic) Broadcast() {
	t.cond.L.Lock()
	t.cond.Broadcast()
	t.cond.L.Unlock()
}

func (t *Topic) Subscribe(lifecycleGoroutine func(), onSignal func()) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		lifecycleGoroutine()
		cancel()
		t.Broadcast()
	}()

	t.SubscribeCtx(ctx, onSignal)
}

func (t *Topic) SubscribeCtx(ctx context.Context, onSignal func()) {
	condSignal := make(chan struct{}, 1)
	go func() {
	READING:
		for {
			t.cond.L.Lock()
			t.cond.Wait()
			t.cond.L.Unlock()
			condSignal <- struct{}{}

			//check if lifecycleEnd closed
			select {
			case _, ok := <-ctx.Done():
				if !ok {
					break READING
				}
			default:
			}
		}
		close(condSignal)
		println("reading finished")
	}()

LOOP:
	for {
		select {
		case _, ok := <-condSignal:
			if ok {
				onSignal()
			}
		case <-ctx.Done():
			break LOOP
		}
	}
}
