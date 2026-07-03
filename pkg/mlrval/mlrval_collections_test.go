package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveIndexedOnArrayInBounds(t *testing.T) {
	array := FromArray([]*Mlrval{FromInt(10), FromInt(20), FromInt(30)})

	err := array.RemoveIndexed([]*Mlrval{FromInt(2)})
	assert.Nil(t, err)

	elements := array.GetArray()
	assert.Equal(t, 2, len(elements))
	assert.Equal(t, int64(10), elements[0].intf.(int64))
	assert.Equal(t, int64(30), elements[1].intf.(int64))
}

func TestRemoveIndexedOnArrayNegativeIndex(t *testing.T) {
	array := FromArray([]*Mlrval{FromInt(10), FromInt(20), FromInt(30)})

	// Negative indices are aliased from the end, so -1 removes the last element.
	err := array.RemoveIndexed([]*Mlrval{FromInt(-1)})
	assert.Nil(t, err)

	elements := array.GetArray()
	assert.Equal(t, 2, len(elements))
	assert.Equal(t, int64(10), elements[0].intf.(int64))
	assert.Equal(t, int64(20), elements[1].intf.(int64))
}

func TestRemoveIndexedOnArrayOutOfBounds(t *testing.T) {
	array := FromArray([]*Mlrval{FromInt(10), FromInt(20), FromInt(30)})

	err := array.RemoveIndexed([]*Mlrval{FromInt(4)})
	assert.NotNil(t, err)
	assert.Equal(t, 3, len(array.GetArray()))

	err = array.RemoveIndexed([]*Mlrval{FromInt(0)})
	assert.NotNil(t, err)
	assert.Equal(t, 3, len(array.GetArray()))
}
