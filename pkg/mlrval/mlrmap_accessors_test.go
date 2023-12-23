package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	mlrmap := NewMlrmap()
	assert.Equal(t, true, mlrmap.IsEmpty())
}

func TestPutReference(t *testing.T) {
	mlrmap := NewMlrmap()
	key1 := "a"
	val1 := FromInt(1)
	mlrmap.PutReference(key1, val1)

	assert.False(t, mlrmap.IsEmpty())

	assert.True(t, mlrmap.Has("a"))
	assert.False(t, mlrmap.Has("b"))
	assert.Equal(t, int64(1), mlrmap.FieldCount)

	read := mlrmap.Get("b")
	assert.Nil(t, read)

	read = mlrmap.Get("a")
	assert.NotNil(t, read)
	intval, ok := read.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(1), intval)

	key2 := "b"
	val2 := FromBool(true)
	mlrmap.PutReference(key2, val2)

	assert.True(t, mlrmap.Has("a"))
	assert.True(t, mlrmap.Has("b"))
	assert.Equal(t, int64(2), mlrmap.FieldCount)

	read = mlrmap.Get("a")
	assert.NotNil(t, read)
	read = mlrmap.Get("b")
	assert.NotNil(t, read)
}

// TODO: TestPrependReference

func TestGetKeysExcept(t *testing.T) {
	mlrmap := NewMlrmap()
	mlrmap.PutReference("a", FromInt(1))
	mlrmap.PutReference("b", FromInt(2))

	exceptions := make(map[string]bool)
	exceptions["x"] = true
	exceptions["y"] = true

	assert.Equal(t, mlrmap.GetKeys(), []string{"a", "b"})
	assert.Equal(t, mlrmap.GetKeysExcept(exceptions), []string{"a", "b"})

	exceptions["a"] = true
	assert.Equal(t, mlrmap.GetKeysExcept(exceptions), []string{"b"})

	exceptions["b"] = true
	assert.Equal(t, mlrmap.GetKeysExcept(exceptions), []string{})
}
