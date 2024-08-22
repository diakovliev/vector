package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFifoQueue(t *testing.T) {
	fifo := NewQueue[int](QueueKindFifo)

	assert.True(t, fifo.Empty())

	for i := 0; i < 5; i++ {
		fifo.Enqueue(i)
	}

	for i := 0; i < 5; i++ {
		assert.Equal(t, i, fifo.Dequeue())
	}
}

func TestLifoQueue(t *testing.T) {
	lifo := NewQueue[int](QueueKindLifo)

	assert.True(t, lifo.Empty())

	for i := 0; i < 5; i++ {
		lifo.Enqueue(i)
	}

	for i := 4; i >= 0; i-- {
		assert.Equal(t, i, lifo.Dequeue())
	}
}
