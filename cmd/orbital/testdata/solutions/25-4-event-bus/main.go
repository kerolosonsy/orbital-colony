package main

import (
	"errors"
	"math"
	"sync"
)

type DockSignal struct{ Topic, Text string }
type dockSubscription struct {
	topic  string
	events chan DockSignal
}
type DockBus struct {
	mu          sync.Mutex
	next        int64
	closed      bool
	subscribers map[int64]dockSubscription
}

func (b *DockBus) Subscribe(topic string, capacity int) (int64, <-chan DockSignal, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if topic == "" || capacity <= 0 || b.closed || b.next == math.MaxInt64 {
		return 0, nil, errors.New("invalid subscription")
	}
	if b.subscribers == nil {
		b.subscribers = make(map[int64]dockSubscription)
	}
	b.next++
	events := make(chan DockSignal, capacity)
	b.subscribers[b.next] = dockSubscription{topic, events}
	return b.next, events, nil
}
func (b *DockBus) Unsubscribe(id int64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	subscription, ok := b.subscribers[id]
	if !ok {
		return false
	}
	delete(b.subscribers, id)
	close(subscription.events)
	return true
}
func (b *DockBus) Publish(signal DockSignal) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed || signal.Topic == "" {
		return 0, errors.New("invalid publication")
	}
	delivered := 0
	for _, subscription := range b.subscribers {
		if subscription.topic != signal.Topic {
			continue
		}
		select {
		case subscription.events <- signal:
			delivered++
		default:
		}
	}
	return delivered, nil
}
func (b *DockBus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for id, subscription := range b.subscribers {
		close(subscription.events)
		delete(b.subscribers, id)
	}
}
func main() {}
