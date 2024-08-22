package vector

import "sync"

// Stack is an interface of stack
type Stack[T any] interface {
	Empty() bool
	Push(T)
	Top() T
	Pop() T
}

// StackImpl is an implementation of stack
type StackImpl[T any] struct {
	Vector *Impl[T]
}

// MakeStack creates a new StackImpl object that holds values of type T.
//
// It takes no parameters and returns a StackImpl[T] object.
func MakeStack[T any]() StackImpl[T] {
	return StackImpl[T]{
		Vector: NewVector[T](),
	}
}

// NewStack creates and returns a new instance of StackImpl with the given type.
//
// T: any type for the elements in the stack.
// Returns a pointer to the new StackImpl instance.
func NewStack[T any]() *StackImpl[T] {
	ret := MakeStack[T]()
	return &ret
}

// WithLocker sets the locker for the StackImpl[T] to the specified value.
//
// locker: a sync.Locker interface implementation used to synchronize access to the stack.
// *StackImpl[T]: returns the modified stack, allowing for method chaining.
func (s *StackImpl[T]) WithLocker(locker sync.Locker) *StackImpl[T] {
	s.Vector.WithLocker(locker)
	return s
}

// Empty returns true if the stack is empty
func (s *StackImpl[T]) Empty() bool {
	return s.Vector.Len() == 0
}

// Push pushes a value onto the stack
func (s *StackImpl[T]) Push(value T) {
	s.Vector.Insert(0, value)
}

// Top returns the the top element of stack. This
// method is not changes stack content
func (s *StackImpl[T]) Top() T {
	return s.Vector.First()
}

// Pop removes the element at the top of the stack
func (s *StackImpl[T]) Pop() T {
	return s.Vector.Remove(0)
}
