package event

import (
	"sync"
)

type HandlerFunc func(payload any)

type Bus struct {
	handlers map[string][]HandlerFunc
	mu       sync.RWMutex
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]HandlerFunc),
	}
}

func (b *Bus) Subscribe(event string, handler HandlerFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[event] = append(b.handlers[event], handler)
}

func (b *Bus) Publish(event string, payload any) {
	b.mu.RLock()
	handlers := b.handlers[event]
	b.mu.RUnlock()

	for _, h := range handlers {
		go h(payload)
	}
}

// Event name constants
const (
	TransactionCreated   = "transaction.created"
	TransactionUpdated   = "transaction.updated"
	TransactionDeleted   = "transaction.deleted"
	BudgetCreated        = "budget.created"
	BudgetUpdated        = "budget.updated"
	BudgetDeleted        = "budget.deleted"
	BillUpdated          = "bill.updated"
	UserLoggedIn         = "user.logged_in"
	UserRegistered       = "user.registered"
	WebhookMessageReady  = "webhook.message.ready"
)
