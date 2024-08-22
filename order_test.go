package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder_FirstIndexOf_IncreasingOrder(t *testing.T) {
	o := NewOrder[int, CompareFunc[int]](CompareNumber[int], OrderKindIncreasing)
	o.Add(33, 36, 34, 34, 34, 35)

	assert.Equal(t, -1, o.FirstIndexOf(32))
	assert.Equal(t, 0, o.FirstIndexOf(33))
	assert.Equal(t, 1, o.FirstIndexOf(34))
	assert.Equal(t, 4, o.FirstIndexOf(35))
	assert.Equal(t, 5, o.FirstIndexOf(36))
	assert.Equal(t, -1, o.FirstIndexOf(37))
}

func TestOrder_FirstIndexOf_DecreasingOrder(t *testing.T) {
	o := NewOrder[int, CompareFunc[int]](CompareNumber[int], OrderKindDecreasing)
	o.Add(33, 36, 34, 34, 34, 35)

	assert.Equal(t, -1, o.FirstIndexOf(32))
	assert.Equal(t, 5, o.FirstIndexOf(33))
	assert.Equal(t, 2, o.FirstIndexOf(34))
	assert.Equal(t, 1, o.FirstIndexOf(35))
	assert.Equal(t, 0, o.FirstIndexOf(36))
	assert.Equal(t, -1, o.FirstIndexOf(37))
}
