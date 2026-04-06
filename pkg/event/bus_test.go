package event

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBus_PublishSubscribe(t *testing.T) {
	bus := NewBus()
	var called atomic.Int32

	bus.Subscribe("test.event", func(payload any) {
		called.Add(1)
	})

	bus.Publish("test.event", "hello")
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, int32(1), called.Load())
}

func TestBus_MultipleSubscribers(t *testing.T) {
	bus := NewBus()
	var called atomic.Int32

	bus.Subscribe("test.event", func(payload any) { called.Add(1) })
	bus.Subscribe("test.event", func(payload any) { called.Add(1) })

	bus.Publish("test.event", nil)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, int32(2), called.Load())
}

func TestBus_NoSubscribers(t *testing.T) {
	bus := NewBus()
	// Should not panic
	bus.Publish("nonexistent.event", nil)
}
