package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectOrderPriorityQueue(t *testing.T) {

	pq := NewPriorityQueue[int]().WithOrder(PriorityQueueOrderDirect)

	assert.True(t, pq.Empty())

	pq.Enqueue(1, 20)
	pq.Enqueue(2, 30)

	pq.Enqueue(0, 10)
	pq.Enqueue(0, 11)
	pq.Enqueue(0, 12)

	pq.Enqueue(1, 21)

	pq.Enqueue(2, 31)

	pq.Enqueue(1, 22)

	assert.False(t, pq.Empty())

	// Elements must out in accordance to decreasing priorities
	// (elements with highest priority will out first).
	assert.Equal(t, 30, pq.Dequeue())
	assert.Equal(t, 31, pq.Dequeue())
	assert.Equal(t, 20, pq.Dequeue())
	assert.Equal(t, 21, pq.Dequeue())
	assert.Equal(t, 22, pq.Dequeue())
	assert.Equal(t, 10, pq.Dequeue())
	assert.Equal(t, 11, pq.Dequeue())
	assert.Equal(t, 12, pq.Dequeue())

	assert.True(t, pq.Empty())

	assert.Panics(t, func() {
		pq.Dequeue()
	})
}

func TestReverseOrderPriorityQueue(t *testing.T) {

	pq := NewPriorityQueue[int]().WithOrder(PriorityQueueOrderReverse)

	assert.True(t, pq.Empty())

	pq.Enqueue(1, 20)
	pq.Enqueue(2, 30)

	pq.Enqueue(0, 10)
	pq.Enqueue(0, 11)
	pq.Enqueue(0, 12)

	pq.Enqueue(1, 21)

	pq.Enqueue(2, 31)

	pq.Enqueue(1, 22)

	assert.False(t, pq.Empty())

	// The elements must out in accordance to increasing priorities
	// (elements with lowest priority will out first).
	assert.Equal(t, 10, pq.Dequeue())
	assert.Equal(t, 11, pq.Dequeue())
	assert.Equal(t, 12, pq.Dequeue())
	assert.Equal(t, 20, pq.Dequeue())
	assert.Equal(t, 21, pq.Dequeue())
	assert.Equal(t, 22, pq.Dequeue())
	assert.Equal(t, 30, pq.Dequeue())
	assert.Equal(t, 31, pq.Dequeue())

	assert.True(t, pq.Empty())

	assert.Panics(t, func() {
		pq.Dequeue()
	})
}
