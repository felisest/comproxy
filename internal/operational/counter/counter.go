package counter

import (
	"sync/atomic"
)

type AtomicCounter struct {
	value atomic.Int64
}

func NewEventCounter() *AtomicCounter {
	return &AtomicCounter{}
}

func (c *AtomicCounter) Inc() {
	c.value.Add(1)
}

func (c *AtomicCounter) Reset() {
	c.value.Store(0)
}

func (c *AtomicCounter) Store(v int64) {
	c.value.Store(v)
}

func (c *AtomicCounter) Value() int64 {
	return c.value.Load()
}
