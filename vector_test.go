package vector

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVector_NewVector(t *testing.T) {
	v := NewVector[int]()
	assert.Equal(t, 0, v.Len())
	assert.Equal(t, 0, v.len())
}

func TestVector_WithLocker(t *testing.T) {
	l := &sync.Mutex{}
	v := NewVector[int]().WithLocker(l)
	assert.Equal(t, l, v.locker)
}

func TestVector_FirstLast(t *testing.T) {
	v := NewVector[int]()
	v.Append(123, 2346, 678, 85)
	assert.Equal(t, 123, v.First())
	assert.Equal(t, 123, v.first())
	assert.Equal(t, 85, v.Last())
	assert.Equal(t, 85, v.last())
}

func TestVector_Get(t *testing.T) {
	v := NewVector[int]()
	v.Append(123, 2346, 678, 85)
	assert.Equal(t, 123, v.Get(0))
	assert.Equal(t, 123, v.get(0))
	assert.Equal(t, 2346, v.Get(1))
	assert.Equal(t, 2346, v.get(1))
	assert.Equal(t, 678, v.Get(2))
	assert.Equal(t, 678, v.get(2))
	assert.Equal(t, 85, v.Get(3))
	assert.Equal(t, 85, v.get(3))
}

func TestVector_Insert(t *testing.T) {

	type testCase struct {
		name     string
		data     []int
		index    uint
		toInsert []int
		expected []int
	}

	testCases := []testCase{
		{
			name:     "single insert at 0",
			data:     []int{1, 2, 3, 4},
			index:    0,
			toInsert: []int{5},
			expected: []int{5, 1, 2, 3, 4},
		},
		{
			name:     "many insert at 0",
			data:     []int{1, 2, 3, 4},
			index:    0,
			toInsert: []int{5, 6},
			expected: []int{5, 6, 1, 2, 3, 4},
		},
		{
			name:     "single insert at middle",
			data:     []int{1, 2, 3, 4},
			index:    2,
			toInsert: []int{5},
			expected: []int{1, 2, 5, 3, 4},
		},
		{
			name:     "many insert at middle",
			data:     []int{1, 2, 3, 4},
			index:    2,
			toInsert: []int{5, 6},
			expected: []int{1, 2, 5, 6, 3, 4},
		},
		{
			name:     "single insert at end",
			data:     []int{1, 2, 3, 4},
			index:    3,
			toInsert: []int{5},
			expected: []int{1, 2, 3, 5, 4},
		},
		{
			name:     "many insert at end",
			data:     []int{1, 2, 3, 4},
			index:    3,
			toInsert: []int{5, 6},
			expected: []int{1, 2, 3, 5, 6, 4},
		},
		{
			name:     "single insert at end+1",
			data:     []int{1, 2, 3, 4},
			index:    4,
			toInsert: []int{5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "many insert at end+1",
			data:     []int{1, 2, 3, 4},
			index:    4,
			toInsert: []int{5, 6},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			v := NewVector[int]()

			v.append(testCase.data...)
			v.insert(testCase.index, testCase.toInsert...)

			assert.Equal(tt, testCase.expected, v.Data())

			v = NewVector[int]()

			v.Append(testCase.data...)
			v.Insert(testCase.index, testCase.toInsert...)

			assert.Equal(tt, testCase.expected, v.Data())
		})
	}
}

func TestVector_Insert_ErrIndexOutOfRange(t *testing.T) {
	v := NewVector[int]()

	v.append(1, 2, 3, 4)
	assert.Panics(t, func() {
		v.insert(5, 1234)
	})

	v = NewVector[int]()

	v.Append(1, 2, 3, 4)
	assert.Panics(t, func() {
		v.Insert(5, 1234)
	})
}

func TestVector_Remove(t *testing.T) {

	type testCase struct {
		name         string
		data         []int
		index        uint
		removeResult int
		expected     []int
	}

	testCases := []testCase{
		{
			name:         "remove at 0",
			data:         []int{1, 2, 3, 4},
			index:        0,
			removeResult: 1,
			expected:     []int{2, 3, 4},
		},
		{
			name:         "remove in middle",
			data:         []int{1, 2, 3, 4},
			index:        2,
			removeResult: 3,
			expected:     []int{1, 2, 4},
		},
		{
			name:         "remove at end",
			data:         []int{1, 2, 3, 4},
			index:        3,
			removeResult: 4,
			expected:     []int{1, 2, 3},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			v := NewVector[int]()

			v.append(testCase.data...)
			assert.Equal(tt, testCase.removeResult, v.remove(testCase.index))
			assert.Equal(tt, testCase.expected, v.Data())

			v = NewVector[int]()

			v.append(testCase.data...)
			assert.Equal(tt, testCase.removeResult, v.Remove(testCase.index))
			assert.Equal(tt, testCase.expected, v.Data())
		})
	}
}

func TestVector_Remove_ErrIndexOutOfRange(t *testing.T) {
	v := NewVector[int]()

	v.append(1, 2, 3, 4)
	assert.Panics(t, func() {
		v.remove(5)
	})

	v = NewVector[int]()

	v.Append(1, 2, 3, 4)
	assert.Panics(t, func() {
		v.Remove(5)
	})
}

func TestVector_Reversed(t *testing.T) {

	v := NewVector[int]()

	v.Append(1, 2, 3, 4)

	assert.Equal(t, []int{4, 3, 2, 1}, v.Reversed().Data())
	assert.Equal(t, []int{1, 2, 3, 4}, v.Data())
}

func TestVector_Range(t *testing.T) {

	v := NewVector[int]()
	v.Append(33)

	counter := 0

	v.Range(func(index int, value int) error {
		assert.Equal(t, 0, index)
		assert.Equal(t, 33, value)
		counter++
		return nil
	})

	v.xrange(func(index int, value int) error {
		assert.Equal(t, 0, index)
		assert.Equal(t, 33, value)
		counter++
		return errors.New("error")
	})

	assert.Equal(t, 2, counter)
}
