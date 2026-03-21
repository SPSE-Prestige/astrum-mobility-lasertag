package eventbus

import (
	"sync"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

// InMemoryEventBus provides a simple in-memory pub/sub for game events.
type InMemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan domain.WSMessage]struct{}
}

func New() *InMemoryEventBus {
	return &InMemoryEventBus{
		subscribers: make(map[string]map[chan domain.WSMessage]struct{}),
	}
}

func (b *InMemoryEventBus) Publish(gameID string, event domain.WSMessage) {
	b.mu.RLock()
	subs := b.subscribers[gameID]
	b.mu.RUnlock()

	for ch := range subs {
		select {
		case ch <- event:
		default:
		}
	}
}

func (b *InMemoryEventBus) Subscribe(gameID string) chan domain.WSMessage {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan domain.WSMessage, 256)
	if b.subscribers[gameID] == nil {
		b.subscribers[gameID] = make(map[chan domain.WSMessage]struct{})
	}
	b.subscribers[gameID][ch] = struct{}{}
	return ch
}

func (b *InMemoryEventBus) Unsubscribe(gameID string, ch chan domain.WSMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if subs, ok := b.subscribers[gameID]; ok {
		delete(subs, ch)
		if len(subs) == 0 {
			delete(b.subscribers, gameID)
		}
	}
}
