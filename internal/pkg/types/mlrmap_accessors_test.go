package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

func TestIsEmpty(t *testing.T) {
	mlrmap := NewMlrmap()
	assert.Equal(t, true, mlrmap.IsEmpty())
}

func TestPutReference(t *testing.T) {
	mlrmap := NewMlrmap()
	key1 := "a"
	val1 := mlrval.FromInt(1)
	mlrmap.PutReference(key1, val1)

	assert.False(t, mlrmap.IsEmpty())

	assert.True(t, mlrmap.Has("a"))
	assert.False(t, mlrmap.Has("b"))
	assert.Equal(t, 1, mlrmap.FieldCount)

	read := mlrmap.Get("b")
	assert.Nil(t, read)

	read = mlrmap.Get("a")
	assert.NotNil(t, read)
	intval, ok := read.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, 1, intval)

	key2 := "b"
	val2 := mlrval.FromBool(true)
	mlrmap.PutReference(key2, val2)

	assert.True(t, mlrmap.Has("a"))
	assert.True(t, mlrmap.Has("b"))
	assert.Equal(t, 2, mlrmap.FieldCount)

	read = mlrmap.Get("a")
	assert.NotNil(t, read)
	read = mlrmap.Get("b")
	assert.NotNil(t, read)
}

// TODO: TestPrependReference
