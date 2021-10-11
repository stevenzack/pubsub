package pubsub

import (
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
	lifecycleEnd := make(chan struct{})
	condSignal := make(chan struct{}, 1)

	go func() {
		lifecycleGoroutine()
		close(lifecycleEnd)
		t.Broadcast()
	}()

	go func() {
	READING:
		for {
			t.cond.L.Lock()
			t.cond.Wait()
			t.cond.L.Unlock()
			condSignal <- struct{}{}

			//check if lifecycleEnd closed
			select {
			case _, ok := <-lifecycleEnd:
				if !ok {
					break READING
				}
			default:
			}
		}
		close(condSignal)
	}()

LOOP:
	for {
		select {
		case _, ok := <-condSignal:
			if ok {
				onSignal()
			}
		case <-lifecycleEnd:
			break LOOP
		}
	}
}
