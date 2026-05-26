// Package stream provides real-time event streaming for SSH sessions,
// allowing consumers to subscribe to session events as they are captured.
package stream

import (
	"errors"
	"sync"

	"github.com/sshtrace/sshtrace/internal/session"
)

// Handler is a function that receives a session event.
type Handler func(sess *session.Session, event session.Event)

// Broker manages subscriptions and broadcasts session events to subscribers.
type Broker struct {
	mu          sync.RWMutex
	subscribers map[string]Handler
}

// New creates a new Broker with no subscribers.
func New() *Broker {
	return &Broker{
		subscribers: make(map[string]Handler),
	}
}

// Subscribe registers a named handler to receive events.
// Returns an error if the name is empty or already registered.
func (b *Broker) Subscribe(name string, h Handler) error {
	if name == "" {
		return errors.New("stream: subscriber name must not be empty")
	}
	if h == nil {
		return errors.New("stream: handler must not be nil")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.subscribers[name]; exists {
		return errors.New("stream: subscriber already registered: " + name)
	}
	b.subscribers[name] = h
	return nil
}

// Unsubscribe removes the handler with the given name.
// Returns an error if the name is not registered.
func (b *Broker) Unsubscribe(name string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.subscribers[name]; !exists {
		return errors.New("stream: subscriber not found: " + name)
	}
	delete(b.subscribers, name)
	return nil
}

// Publish sends the event to all registered subscribers.
// Each handler is called synchronously in an unspecified order.
func (b *Broker) Publish(sess *session.Session, event session.Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.subscribers {
		h(sess, event)
	}
}

// Count returns the number of active subscribers.
func (b *Broker) Count() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subscribers)
}
