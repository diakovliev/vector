package vector

import (
	"errors"
	"sync"
)

var (
	// ErrEmptyPriorityQueue is raised by Dequeue when the priority queue is empty
	ErrEmptyPriorityQueue = errors.New("empty priority queue")
)

// PriorityQueueOrder is a function that compares two priorities
type PriorityQueueOrder func(int, int) bool

var (
	// PriorityQueueOrderDirect is a direct order of priority queue
	PriorityQueueOrderDirect PriorityQueueOrder = func(lhs, rhs int) bool {
		return lhs > rhs
	}

	// PriorityQueueOrderReverse is a reverse order of priority queue
	PriorityQueueOrderReverse PriorityQueueOrder = func(lhs, rhs int) bool {
		return lhs < rhs
	}

	// PriorityQueueOrderDefault is a default order of priority queue (PriorityQueueOrderDirect)
	PriorityQueueOrderDefault = PriorityQueueOrderDirect
)

// PriorityQueueElement is a structure for priority queue element
type PriorityQueueElement[T any] struct {
	Priority int
	Value    T
}

// PriorityQueue is an interface of priority queue
type PriorityQueue[T any] interface {
	Len() int
	Empty() bool
	Enqueue(int, T)
	Dequeue() T
}

// PriorityQueueImpl is an implementation of priority queue
type PriorityQueueImpl[T any] struct {
	Vector               *Impl[PriorityQueueElement[T]]
	prioritiesComparator PriorityQueueOrder
}

// MakePriorityQueue returns a new instance of PriorityQueueImpl[T]. It creates a priority queue
// with default PriorityQueueOrderDefault. PriorityQueueImpl[T] is a struct that contains a Vector
// of PriorityQueueElement[T] and a prioritiesComparator function pointer.
//
// Returns a PriorityQueueImpl[T].
func MakePriorityQueue[T any]() PriorityQueueImpl[T] {
	return PriorityQueueImpl[T]{
		Vector:               NewVector[PriorityQueueElement[T]](),
		prioritiesComparator: PriorityQueueOrderDefault,
	}
}

// NewPriorityQueue returns a new instance of PriorityQueueImpl[T].
//
// This function takes no parameters.
// It returns a pointer to a PriorityQueueImpl[T] instance.
func NewPriorityQueue[T any]() *PriorityQueueImpl[T] {
	ret := MakePriorityQueue[T]()
	return &ret
}

// WithLocker sets the locker to be used by the priority queue and returns the updated priority queue.
//
// locker: a sync.Locker implementation to be used to synchronize access to the priority queue.
// returns: a pointer to the updated priority queue.
func (pq *PriorityQueueImpl[T]) WithLocker(locker sync.Locker) *PriorityQueueImpl[T] {
	pq.Vector.WithLocker(locker)
	return pq
}

// WithOrder creates a new PriorityQueue with the given PriorityQueueOrder and returns it.
//
// order: The PriorityQueueOrder to use for the newly created PriorityQueue.
// PriorityQueue[T]: Returns a new PriorityQueue with the given PriorityQueueOrder.
func (pq *PriorityQueueImpl[T]) WithOrder(
	order PriorityQueueOrder,
) PriorityQueue[T] {
	return &PriorityQueueImpl[T]{
		Vector:               pq.Vector,
		prioritiesComparator: order,
	}
}

// len returns the number of items in the queue
func (pq *PriorityQueueImpl[T]) len() int {
	return pq.Vector.len()
}

// Len returns the number of elements in the PriorityQueue.
//
// No parameters are needed.
// An integer is returned that represents the number of elements in the PriorityQueue.
func (pq *PriorityQueueImpl[T]) Len() int {
	pq.Vector.Locker().Lock()
	defer pq.Vector.Locker().Unlock()

	return pq.len()
}

// empty returns true if there are no items in the queue
func (pq *PriorityQueueImpl[T]) empty() bool {
	return pq.Vector.len() == 0
}

// Empty checks if the priority queue is empty.
//
// pq *PriorityQueueImpl[T]: pointer to a PriorityQueueImpl[T] struct.
// bool: returns true if the priority queue is empty, false otherwise.
func (pq *PriorityQueueImpl[T]) Empty() bool {
	pq.Vector.Locker().Lock()
	defer pq.Vector.Locker().Unlock()

	return pq.empty()
}

// enqueue enqueues a value to the queue
func (pq *PriorityQueueImpl[T]) enqueue(priority int, value T) {
	newElement := PriorityQueueElement[T]{Priority: priority, Value: value}

	for index := 0; index < pq.Vector.len(); index++ {
		if pq.prioritiesComparator(
			newElement.Priority,
			pq.Vector.get(uint(index)).Priority,
		) {
			pq.Vector.insert(uint(index), newElement)
			return
		}
	}

	pq.Vector.append(newElement)
}

// Enqueue adds an element to the priority queue with the given priority.
//
// priority: an integer representing the priority of the element.
// value: the element to be added to the priority queue.
func (pq *PriorityQueueImpl[T]) Enqueue(priority int, value T) {
	pq.Vector.Locker().Lock()
	defer pq.Vector.Locker().Unlock()

	pq.enqueue(priority, value)
}

// dequeue removes and returns the element from the queue
func (pq *PriorityQueueImpl[T]) dequeue() (ret T) {
	if pq.empty() {
		panic(ErrEmptyPriorityQueue)
	}

	ret = pq.Vector.first().Value
	pq.Vector.remove(0)
	return
}

// Dequeue removes and returns the highest priority item from the priority queue.
// Panics if the queue is empty.
//
// It doesn't take any parameters and returns the dequeued item of type T.
func (pq *PriorityQueueImpl[T]) Dequeue() (ret T) {
	pq.Vector.Locker().Lock()
	defer pq.Vector.Locker().Unlock()

	return pq.dequeue()
}
