package vector

import "sync"

// Set is an interface of set
type Set[T any, C CompareFunc[T]] interface {
	Empty() bool
	Add(...T) int
	Remove(...T) int
	Has(T) bool
	HasAny(...T) bool
	HasAll(...T) bool
	Range(func(index int, value T) error) error
	Data() []T
}

// SetImpl is an implementation of set
type SetImpl[T any, C CompareFunc[T]] struct {
	Order   *OrderImpl[T, C]
	compare C
}

// MakeSet returns a new SetImpl with a given compare function.
// It takes in a type T and a CompareFunc C.
// Returns a SetImpl with a new Order based on the compare function.
func MakeSet[T any, C CompareFunc[T]](compareFunc C) SetImpl[T, C] {
	return SetImpl[T, C]{
		Order:   NewOrder[T](compareFunc, OrderKindIncreasing),
		compare: compareFunc,
	}
}

// NewSet returns a new SetImpl instance.
// T is the type of the elements in the set.
// C is the type of the compare function.
// compareFunc is the function used to compare elements.
func NewSet[T any, C CompareFunc[T]](compareFunc C) *SetImpl[T, C] {
	ret := MakeSet[T](compareFunc)
	return &ret
}

// WithLocker sets the locker to be used by set and returns the updated set.
//
// locker: a sync.Locker implementation to be used to synchronize access to the set.
// returns: a pointer to the updated set.
func (s *SetImpl[T, C]) WithLocker(locker sync.Locker) *SetImpl[T, C] {
	s.Order.WithLocker(locker)
	return s
}

// Data returns a set data.
func (s *SetImpl[T, C]) Data() []T {
	return s.Order.Data()
}

// Empty checks if the set is empty.
func (s *SetImpl[T, C]) Empty() bool {
	return s.Order.Vector.Len() == 0
}

// Add elements to the set.
func (s *SetImpl[T, C]) add(values ...T) (count int) {

	for _, v := range values {
		index := s.Order.firstIndexOf(v)
		if index == -1 {
			s.Order.add(v)
			count++
		}
	}

	return count
}

// Add elements to the set.
func (s *SetImpl[T, C]) Add(values ...T) (count int) {
	s.Order.Locker().Lock()
	defer s.Order.Locker().Unlock()

	return s.add(values...)
}

// Remove elements from the set.
func (s *SetImpl[T, C]) remove(values ...T) (count int) {

	for _, v := range values {
		index := s.Order.firstIndexOf(v)
		if index != -1 {
			s.Order.Vector.remove(uint(index))
			count++
		}
	}

	return count
}

// Remove elements from the set.
func (s *SetImpl[T, C]) Remove(values ...T) (count int) {
	s.Order.Locker().Lock()
	defer s.Order.Locker().Unlock()

	return s.remove(values...)
}

// Range enumerates set elements.
func (s *SetImpl[T, C]) Range(callback func(index int, value T) error) error {
	s.Order.Locker().Lock()
	defer s.Order.Locker().Unlock()

	return s.Order.Vector.xrange(callback)
}

// Union constructs a new set of the elements what are available in both original sets.
func (s *SetImpl[T, C]) Union(rhs *SetImpl[T, C]) *SetImpl[T, C] {

	s.Order.Locker().Lock()
	rhs.Order.Locker().Lock()
	defer func() {
		rhs.Order.Locker().Unlock()
		s.Order.Locker().Unlock()
	}()

	ret := NewSet[T](s.compare)

	//ret.add(s.Data()...)
	//ret.add(rhs.Data()...)
	ret.Order = s.Order.combine(rhs.Order)

	return ret
}

// LeftDifference constructs a new set of the elements what are available in set and not available in rhs:
// +------+---+------+
// |xxxxxx|   |      |
// | s xxx|   | rhs  |
// |   xxx|   |      |
// |xxxxxx|   |      |
// +------+---+------+
func (s *SetImpl[T, C]) LeftDifference(rhs *SetImpl[T, C]) *SetImpl[T, C] {

	s.Order.Locker().Lock()
	rhs.Order.Locker().Lock()
	defer func() {
		rhs.Order.Locker().Unlock()
		s.Order.Locker().Unlock()
	}()

	ret := NewSet[T](s.compare)

	var values []T
	for _, v := range s.Data() {
		rhsIndex := rhs.Order.firstIndexOf(v)
		if rhsIndex == -1 {
			values = append(values, v)
		}
	}
	ret.add(values...)

	return ret
}

// RightDifference constructs a new set of the elements what are available in rhs and not available in set:
// +------+---+------+
// |      |   |xxxxxx|
// |  s   |   | rhs x|
// |      |   |     x|
// |      |   |xxxxxx|
// +------+---+------+
func (s *SetImpl[T, C]) RightDifference(rhs *SetImpl[T, C]) *SetImpl[T, C] {

	s.Order.Locker().Lock()
	rhs.Order.Locker().Lock()
	defer func() {
		rhs.Order.Locker().Unlock()
		s.Order.Locker().Unlock()
	}()

	ret := NewSet[T](s.compare)

	var values []T
	for _, v := range rhs.Data() {
		lhsIndex := s.Order.firstIndexOf(v)
		if lhsIndex == -1 {
			values = append(values, v)
		}
	}
	ret.add(values...)

	return ret
}

// Intersection constructs a new set of the elements what are available in both sets:
// +------+---+------+
// |      |xxx|      |
// |  s   |xxx| rhs  |
// |      |xxx|      |
// |      |xxx|      |
// +------+---+------+
func (s *SetImpl[T, C]) Intersection(rhs *SetImpl[T, C]) *SetImpl[T, C] {

	s.Order.Locker().Lock()
	rhs.Order.Locker().Lock()
	defer func() {
		rhs.Order.Locker().Unlock()
		s.Order.Locker().Unlock()
	}()

	ret := NewSet[T](s.compare)

	var values []T
	for _, v := range rhs.Data() {
		lhsIndex := s.Order.firstIndexOf(v)
		if lhsIndex != -1 {
			values = append(values, v)
		}
	}
	for _, v := range s.Data() {
		rhsIndex := rhs.Order.firstIndexOf(v)
		if rhsIndex != -1 {
			values = append(values, v)
		}
	}
	ret.add(values...)

	return ret
}

// Has checks if the set contains the given value.
func (s *SetImpl[T, C]) Has(value T) bool {
	s.Order.Locker().Lock()
	defer s.Order.Locker().Unlock()

	return s.Order.firstIndexOf(value) != -1
}

// HasAny checks if the set contains any of the given values.
func (s *SetImpl[T, C]) HasAny(values ...T) bool {
	s.Order.Locker().Lock()
	defer s.Order.Locker().Unlock()

	var counter int
	for _, v := range values {
		if s.Order.firstIndexOf(v) != -1 {
			counter++
		}
	}

	return counter > 0
}

// HasAll checks if the set contains all of the given values.
func (s *SetImpl[T, C]) HasAll(values ...T) bool {
	s.Order.Locker().Lock()
	defer s.Order.Locker().Unlock()

	var counter int
	for _, v := range values {
		if s.Order.firstIndexOf(v) != -1 {
			counter++
		}
	}

	return len(values) == counter
}
