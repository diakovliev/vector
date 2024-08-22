package vector

import (
	"sort"
	"sync"
)

// OrderKind is either increasing or decreasing
type OrderKind int

const (
	// OrderKindIncreasing is the increasing order
	OrderKindIncreasing OrderKind = 1

	// OrderKindDecreasing is the decreasing order
	OrderKindDecreasing OrderKind = -1
)

// Order is an interface of order
type Order[T any, C CompareFunc[T]] interface {
	Empty() bool
	Kind() OrderKind
	Data() []T
	Add(value T) uint
	FirstIndexOf(value T) int
}

// OrderImpl is an implementation of order
type OrderImpl[T any, C CompareFunc[T]] struct {
	Vector  *Impl[T]
	compare C
	kind    OrderKind
}

// MakeOrder creates a new instance of OrderImpl with the given type and compare function
// and returns it.
//
// Args:
//
//	compareFunc: The compare function to use for comparing elements of the order.
//	kind: The order kind to use.
//
// Returns:
//
//	A new instance of OrderImpl.
func MakeOrder[T any, C CompareFunc[T]](compareFunc C, kind OrderKind) OrderImpl[T, C] {
	return OrderImpl[T, C]{
		Vector:  NewVector[T](),
		compare: compareFunc,
		kind:    kind,
	}
}

// NewOrder creates a new instance of OrderImpl with the given type and compare function
// and returns it.
//
// Args:
//
//	compareFunc: The compare function to use for comparing elements of the order.
//	kind: The order kind to use.
//
// Returns:
//
//	A new instance of OrderImpl.
func NewOrder[T any, C CompareFunc[T]](compareFunc C, kind OrderKind) *OrderImpl[T, C] {
	ret := MakeOrder[T](compareFunc, kind)
	return &ret
}

// WithLocker sets the locker to be used by order and returns the updated order.
//
// locker: a sync.Locker implementation to be used to synchronize access to the order.
// returns: a pointer to the updated order.
func (o *OrderImpl[T, C]) WithLocker(locker sync.Locker) *OrderImpl[T, C] {
	o.Vector.WithLocker(locker)
	return o
}

// Locker returns the locker to be used by order
func (o *OrderImpl[T, C]) Locker() sync.Locker {
	return o.Vector.Locker()
}

// Empty checks if container is empty
func (o *OrderImpl[T, C]) Empty() bool {
	return o.Vector.Len() == 0
}

// Kind returns an order kind
func (o *OrderImpl[T, C]) Kind() OrderKind {
	return o.kind
}

// Data returns an order data
func (o *OrderImpl[T, C]) Data() []T {
	return o.Vector.Data()
}

// Add element(s) to order, result is count of added elements
func (o *OrderImpl[T, C]) add(values ...T) (count uint) {

	for _, value := range values {

		index := sort.Search(o.Vector.len(), func(i int) bool {
			switch o.kind {
			case OrderKindIncreasing:
				return o.compare(o.Vector.get(uint(i)), value) > 0
			case OrderKindDecreasing:
				return o.compare(o.Vector.get(uint(i)), value) < 0
			default:
				panic("unsupported order kind")
			}
		})

		if index > o.Vector.len()-1 {
			index = o.Vector.len()
		}

		o.Vector.insert(uint(index), value)

		count++
	}

	return
}

// Add element(s) to order, result is count of added elements
func (o *OrderImpl[T, C]) Add(values ...T) (count uint) {
	o.Vector.Locker().Lock()
	defer o.Vector.Locker().Unlock()

	return o.add(values...)
}

// Find element first occurrence index by value
func (o *OrderImpl[T, C]) firstIndexOf(value T) int {
	if o.Vector.len() == 0 {
		return -1
	}
	index := sort.Search(o.Vector.len(), func(i int) bool {
		switch o.kind {
		case OrderKindIncreasing:
			return o.compare(o.Vector.get(uint(i)), value) >= 0
		case OrderKindDecreasing:
			return o.compare(o.Vector.get(uint(i)), value) <= 0
		default:
			panic("unsupported order kind")
		}
	})
	if index < o.Vector.len() && o.compare(o.Vector.get(uint(index)), value) == 0 {
		return index
	}
	return -1
}

// FirstIndexOf finds an element first occurrence index by value
func (o *OrderImpl[T, C]) FirstIndexOf(value T) int {
	o.Vector.Locker().Lock()
	defer o.Vector.Locker().Unlock()

	return o.firstIndexOf(value)
}

// Merge orders
func (o *OrderImpl[T, C]) merge(rhs *OrderImpl[T, C]) *OrderImpl[T, C] {
	ret := NewOrder[T](o.compare, o.kind)

	lhsIndex := 0
	rhsIndex := 0

	for lhsIndex < o.Vector.len() && rhsIndex < rhs.Vector.len() {
		lhsValue := o.Vector.get(uint(lhsIndex))
		rhsValue := rhs.Vector.get(uint(rhsIndex))

		compareRes := o.compare(lhsValue, rhsValue)

		if (o.kind == OrderKindIncreasing && compareRes <= 0) || (o.kind == OrderKindDecreasing && compareRes >= 0) {
			ret.Vector.append(lhsValue)
			lhsIndex++
		} else {
			ret.Vector.append(rhsValue)
			rhsIndex++
		}
	}

	if lhsIndex < o.Vector.len() {
		ret.Vector.append(o.Vector.Data()[lhsIndex:]...)
	} else if rhsIndex < rhs.Vector.len() {
		ret.Vector.append(rhs.Vector.Data()[rhsIndex:]...)
	}

	return ret
}

// Merge orders
func (o *OrderImpl[T, C]) Merge(rhs *OrderImpl[T, C]) *OrderImpl[T, C] {
	o.Vector.Locker().Lock()
	rhs.Vector.Locker().Lock()
	defer func() {
		rhs.Vector.Locker().Unlock()
		o.Vector.Locker().Unlock()
	}()

	return o.merge(rhs)
}

// Merge orders and omit non unique elements in resulting order
func (o *OrderImpl[T, C]) combine(rhs *OrderImpl[T, C]) *OrderImpl[T, C] {
	ret := NewOrder[T](o.compare, o.kind)
	lhsIndex := 0
	rhsIndex := 0
	var prevValue *T

	for lhsIndex < o.Vector.len() && rhsIndex < rhs.Vector.len() {
		lhsValue := o.Vector.get(uint(lhsIndex))
		rhsValue := rhs.Vector.get(uint(rhsIndex))

		compareRes := o.compare(lhsValue, rhsValue)

		if (o.kind == OrderKindIncreasing && compareRes <= 0) || (o.kind == OrderKindDecreasing && compareRes >= 0) {
			if prevValue == nil || o.compare(*prevValue, lhsValue) != 0 {
				ret.Vector.append(lhsValue)
				prevValue = &lhsValue
			}
			lhsIndex++
		} else {
			if prevValue == nil || o.compare(*prevValue, rhsValue) != 0 {
				ret.Vector.append(rhsValue)
				prevValue = &rhsValue
			}
			rhsIndex++
		}
	}

	var rest []T
	switch {
	case lhsIndex >= o.Vector.len() && rhsIndex < rhs.Vector.len():
		rest = rhs.Vector.Data()[rhsIndex:]
	case rhsIndex >= rhs.Vector.len() && lhsIndex < o.Vector.len():
		rest = o.Vector.Data()[lhsIndex:]
	}
	for _, restValue := range rest {
		restValue := restValue
		if prevValue == nil || o.compare(*prevValue, restValue) != 0 {
			ret.Vector.append(restValue)
			prevValue = &restValue
		}
	}

	return ret
}

// Combine merges orders and omit non unique elements in resulting order
func (o *OrderImpl[T, C]) Combine(rhs *OrderImpl[T, C]) *OrderImpl[T, C] {
	o.Vector.Locker().Lock()
	rhs.Vector.Locker().Lock()
	defer func() {
		rhs.Vector.Locker().Unlock()
		o.Vector.Locker().Unlock()
	}()

	return o.combine(rhs)
}
