package vector

import (
	"errors"
	"sync"
)

// QueueKind is the type of queue
type QueueKind uint

const (
	// QueueKindFifo is the FIFO queue
	QueueKindFifo QueueKind = iota
	// QueueKindLifo is the LIFO queue
	QueueKindLifo
)

var (
	// ErrEmptyQueue is raised by Dequeue when the queue is empty
	ErrEmptyQueue = errors.New("empty queue")
)

// Queue is an interface of queue
type Queue[T any] interface {
	Len() int
	Empty() bool
	Enqueue(T)
	Dequeue() T
}

// QueueImpl is an implementation of queue
type QueueImpl[T any] struct {
	Vector *Impl[T]
	kind   QueueKind
}

// MakeQueue creates a new QueueImpl with the given QueueKind.
//
// kind: The type of queue (FIFO or LIFO).
// Returns a new QueueImpl.
func MakeQueue[T any](kind QueueKind) QueueImpl[T] {
	return QueueImpl[T]{
		Vector: NewVector[T](),
		kind:   kind,
	}
}

// NewQueue creates a new QueueImpl[T] given a QueueKind.
//
// kind: the QueueKind to create the queue with.
// returns: a pointer to the newly created QueueImpl[T].
func NewQueue[T any](kind QueueKind) *QueueImpl[T] {
	ret := MakeQueue[T](kind)
	return &ret
}

// WithLocker sets the locker for the queue.
//
// locker: the synchronization locker to be used.
// Returns a pointer to the modified queue.
func (q *QueueImpl[T]) WithLocker(locker sync.Locker) *QueueImpl[T] {
	q.Vector.WithLocker(locker)
	return q
}

func (q *QueueImpl[T]) len() int {
	return q.Vector.len()
}

// Len returns the number of elements in the queue.
//
// It doesn't take any parameters.
// The return type is an int.
func (q *QueueImpl[T]) Len() int {
	q.Vector.Locker().Lock()
	defer q.Vector.Locker().Unlock()

	return q.len()
}

func (q *QueueImpl[T]) empty() bool {
	return q.Vector.len() == 0
}

// Empty checks if the queue is empty.
//
// No parameter is needed.
// Returns a boolean indicating if the queue is empty or not.
func (q *QueueImpl[T]) Empty() bool {
	q.Vector.Locker().Lock()
	defer q.Vector.Locker().Unlock()

	return q.empty()
}

// enqueue enqueues a value to the queue
func (q *QueueImpl[T]) enqueue(value T) {
	switch q.kind {
	case QueueKindFifo:
		q.Vector.append(value)
	case QueueKindLifo:
		q.Vector.insert(0, value)
	}
}

// Enqueue adds an element to the back of the queue.
//
// value: the element to be added to the queue.
func (q *QueueImpl[T]) Enqueue(value T) {
	q.Vector.Locker().Lock()
	defer q.Vector.Locker().Unlock()

	q.enqueue(value)
}

// dequeue removes and returns the element from the queue
func (q *QueueImpl[T]) dequeue() (ret T) {
	if q.empty() {
		panic(ErrEmptyQueue)
	}

	ret = q.Vector.first()
	q.Vector.remove(0)
	return
}

// Dequeue removes and returns the first element from the queue.
// Panics if the queue is empty.
//
// T, the type of the queue elements.
func (q *QueueImpl[T]) Dequeue() (ret T) {
	q.Vector.Locker().Lock()
	defer q.Vector.Locker().Unlock()

	return q.dequeue()
}
