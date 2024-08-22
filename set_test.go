package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet_NewSet(t *testing.T) {
	s := NewSet[int, CompareFunc[int]](CompareNumber[int]).WithLocker(NewLockerStub())
	assert.True(t, s.Empty())
}

func TestSet_Add(t *testing.T) {

	s := NewSet[int, CompareFunc[int]](CompareNumber[int])
	assert.True(t, s.Empty())

	assert.Equal(t, 4, s.Add(12, 223, 3456, 456))
	assert.Equal(t, 0, s.Add(12, 223, 3456, 456))
	assert.Equal(t, 2, s.Add(13, 223, 3457, 456))
	assert.False(t, s.Empty())
}

func TestSet_Remove(t *testing.T) {

	s := NewSet[int, CompareFunc[int]](CompareNumber[int])
	assert.True(t, s.Empty())

	assert.Equal(t, 4, s.Add(12, 223, 3456, 456))
	assert.Equal(t, 2, s.Remove(12, 3456))
	assert.Equal(t, 1, s.Remove(456))
	assert.Equal(t, 1, s.Remove(223))
	assert.True(t, s.Empty())
}

func TestSet_Union(t *testing.T) {
	s0 := NewSet[int, CompareFunc[int]](CompareNumber[int])
	s1 := NewSet[int, CompareFunc[int]](CompareNumber[int])

	s0.Add(1, 2, 3, 4, 5)
	s1.Add(4, 5, 6, 7, 8)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, s0.Union(s1).Data())
}

func TestSet_LeftDifference(t *testing.T) {
	s0 := NewSet[int, CompareFunc[int]](CompareNumber[int])
	s1 := NewSet[int, CompareFunc[int]](CompareNumber[int])

	s0.Add(1, 2, 3, 4, 5)
	s1.Add(4, 5, 6, 7, 8)

	assert.Equal(t, []int{1, 2, 3}, s0.LeftDifference(s1).Data())
	assert.Equal(t, []int{6, 7, 8}, s1.LeftDifference(s0).Data())
}

func TestSet_RightDifference(t *testing.T) {
	s0 := NewSet[int, CompareFunc[int]](CompareNumber[int])
	s1 := NewSet[int, CompareFunc[int]](CompareNumber[int])

	s0.Add(1, 2, 3, 4, 5)
	s1.Add(4, 5, 6, 7, 8)

	assert.Equal(t, []int{6, 7, 8}, s0.RightDifference(s1).Data())
	assert.Equal(t, []int{1, 2, 3}, s1.RightDifference(s0).Data())
}

func TestSet_Intersection(t *testing.T) {
	s0 := NewSet[int, CompareFunc[int]](CompareNumber[int])
	s1 := NewSet[int, CompareFunc[int]](CompareNumber[int])

	s0.Add(1, 2, 3, 4, 5)
	s1.Add(4, 5, 6, 7, 8)

	assert.Equal(t, []int{4, 5}, s0.Intersection(s1).Data())
	assert.Equal(t, []int{4, 5}, s1.Intersection(s0).Data())
}

func TestStringsSet(t *testing.T) {
	s := NewStringsSet([]string{"first", "second"})

	assert.True(t, s.Has("first"))
	assert.True(t, s.Has("second"))
	assert.False(t, s.Has("third"))
}
