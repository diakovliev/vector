// Package vector implements mutable vector data structures based on it.
package vector

import (
	"errors"
	"sort"
	"sync"
)

var (
	// ErrIndexOutOfRange raised by access methods when the index is out of range
	ErrIndexOutOfRange = errors.New("index out of range")

	// ErrEmptyVector raised by access methods when the vector is empty
	ErrEmptyVector = errors.New("empty vector")
)

// Vector is an interface of a vector
type Vector[T any] interface {
	Data() []T
	Len() int
	First() T
	Last() T
	Set(index uint, value T)
	Get(index uint) T
	Append(value ...T)
	Insert(index uint, value ...T)
	Remove(index uint) T
	Range(func(index int, value T) error) error
}

// Impl is an implementation of a vector
type Impl[T any] struct {
	locker sync.Locker
	data   []T
}

// MakeVector creates and returns a Impl of type T.
//
// The function takes no parameters. It returns a Impl of type T.
func MakeVector[T any]() Impl[T] {
	return Impl[T]{
		locker: NewLockerStub(),
		data:   make([]T, 0),
	}
}

// NewVector returns a new instance of Impl[T].
//
// This function takes no parameters.
// It returns a pointer to a Impl[T] object.
func NewVector[T any]() *Impl[T] {
	ret := MakeVector[T]()
	return &ret
}

// WithLocker sets the locker for the Impl instance and returns a pointer to it.
//
// locker: a sync.Locker to set as the locker for the Impl instance
// *Impl[T]: a pointer to the Impl instance
func (v *Impl[T]) WithLocker(locker sync.Locker) *Impl[T] {
	v.locker = locker
	return v
}

// Locker returns the sync.Locker of the Impl.
//
// No parameters.
// Returns a sync.Locker.
func (v *Impl[T]) Locker() sync.Locker {
	return v.locker
}

// Data retrieves the underlying data of the Impl[T].
//
// No parameters.
// Returns a slice of the type T.
func (v *Impl[T]) Data() []T {
	return v.data
}

// len returns the length of the vector
func (v *Impl[T]) len() int {
	return len(v.data)
}

// Len returns the length of the Impl[T].
//
// This method does not take any parameters.
// It returns an int representing the length of the vector.
func (v *Impl[T]) Len() int {
	v.locker.Lock()
	defer v.locker.Unlock()

	return v.len()
}

// first returns the first element
func (v *Impl[T]) first() T {
	if len(v.data) == 0 {
		panic(ErrEmptyVector)
	}

	return v.data[0]
}

// First returns the first element of the vector.
//
// No parameters.
// Returns the element of type T.
func (v *Impl[T]) First() T {
	v.locker.Lock()
	defer v.locker.Unlock()

	return v.first()
}

// last returns the last element
func (v *Impl[T]) last() T {
	if len(v.data) == 0 {
		panic(ErrEmptyVector)
	}

	return v.data[len(v.data)-1]
}

// Last returns the last element of the vector.
//
// No parameters.
// Returns the type T of the vector.
func (v *Impl[T]) Last() T {
	v.locker.Lock()
	defer v.locker.Unlock()

	return v.last()
}

// set sets the value at the given index
func (v *Impl[T]) set(index uint, value T) {
	if int(index) > len(v.data)-1 {
		panic(ErrIndexOutOfRange)
	}

	v.data[index] = value
}

// Set sets the value of the element at the given index in the Vector.
//
// Parameters:
// - index: the index of the value to be set.
// - value: the new value to be set.
func (v *Impl[T]) Set(index uint, value T) {
	v.locker.Lock()
	defer v.locker.Unlock()

	v.set(index, value)
}

// get returns the element at the given index
func (v *Impl[T]) get(index uint) T {
	if int(index) > len(v.data)-1 {
		panic(ErrIndexOutOfRange)
	}

	return v.data[index]
}

// Get returns the element at the given index in the Impl.
//
// index: the index of the element to be retrieved.
// returns: the element at the given index.
func (v *Impl[T]) Get(index uint) T {
	v.locker.Lock()
	defer v.locker.Unlock()

	return v.get(index)
}

// append appends a value to the
func (v *Impl[T]) append(args ...T) {
	v.data = append(v.data, args...)
}

// Append adds elements to the end of the Impl.
//
// args: The element(s) to be added.
func (v *Impl[T]) Append(args ...T) {
	v.locker.Lock()
	defer v.locker.Unlock()

	v.append(args...)
}

// insert inserts one or more elements to the Impl at the specified index.
func (v *Impl[T]) insert(index uint, args ...T) {
	oldData := v.data
	v.data = make([]T, 0)

	switch {
	case index == 0:
		v.data = append(v.data, args...)
		v.data = append(v.data, oldData...)
	case int(index) == len(oldData):
		v.data = append(v.data, oldData...)
		v.data = append(v.data, args...)
	case int(index) > len(oldData):
		panic(ErrIndexOutOfRange)
	default:
		v.data = append(v.data, oldData[:index]...)
		v.data = append(v.data, args...)
		v.data = append(v.data, oldData[index:]...)
	}
}

// Insert inserts one or more elements to the Impl at the specified index.
//
// index is the position where the elements should be inserted. args is a variable
// number of elements of type T to be inserted.
// There is no return value.
func (v *Impl[T]) Insert(index uint, args ...T) {
	v.locker.Lock()
	defer v.locker.Unlock()

	v.insert(index, args...)
}

// remove removes the element at the given index
func (v *Impl[T]) remove(index uint) (ret T) {
	if len(v.data) == 0 {
		panic(ErrEmptyVector)
	}

	ret = v.data[index]
	oldData := v.data
	v.data = make([]T, 0)

	switch {
	case index == 0:
		v.data = append(v.data, oldData[1:]...)
	case int(index) == len(oldData)-1:
		v.data = append(v.data, oldData[:len(oldData)-1]...)
	case int(index) >= len(oldData):
		panic(ErrIndexOutOfRange)
	default:
		v.data = append(v.data, oldData[:index]...)
		v.data = append(v.data, oldData[index+1:]...)
	}

	return
}

// Remove removes an element from the vector at the given index.
//
// index: the index of the element to be removed.
// ret: the removed element.
func (v *Impl[T]) Remove(index uint) (ret T) {
	v.locker.Lock()
	defer v.locker.Unlock()

	return v.remove(index)
}

// xrange calls the callback for each element
func (v *Impl[T]) xrange(callback func(index int, value T) error) error {
	for index, value := range v.data {
		if err := callback(index, value); err != nil {
			return err
		}
	}
	return nil
}

// Range retrieves all the elements in the Impl[T] object and calls the callback function with the
// index and value of each element. If the callback returns an error, the iteration stops and returns the error.
//
// The callback function takes in an int representing the index and a value of type T representing the element
// in the Impl[T] object. It returns an error.
//
// Returns an error if the underlying xrange function fails.
func (v *Impl[T]) Range(callback func(index int, value T) error) error {
	v.locker.Lock()
	defer v.locker.Unlock()

	return v.xrange(callback)
}

// Reversed returns a new vector with elements in reverse order.
//
// No parameters.
// Returns a pointer to a Impl[T] instance.
func (v *Impl[T]) Reversed() (ret *Impl[T]) {
	newData := make([]T, 0)
	v.Locker().Lock()
	newData = append(newData, v.Data()...)
	v.Locker().Unlock()
	sort.SliceStable(newData, func(i, j int) bool {
		return i > j
	})
	newVector := NewVector[T]()
	newVector.Append(newData...)
	return newVector
}
