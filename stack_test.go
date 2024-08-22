package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	stack := NewStack[int]()
	assert.True(t, stack.Empty())
	stack.Push(1234)
	stack.Push(1235)
	assert.False(t, stack.Empty())
	assert.Equal(t, 1235, stack.Top())
	assert.Equal(t, 1235, stack.Pop())
	assert.Equal(t, 1234, stack.Top())
	assert.Equal(t, 1234, stack.Pop())
	assert.True(t, stack.Empty())
}
