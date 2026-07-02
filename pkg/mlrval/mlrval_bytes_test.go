package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromBytes(t *testing.T) {
	mv := FromBytes([]byte{0xff, 0x01})
	assert.Equal(t, MT_BYTES, mv.Type())
	assert.Equal(t, "bytes", mv.GetTypeName())
	assert.True(t, mv.IsBytes())
	assert.Equal(t, []byte{0xff, 0x01}, mv.AcquireBytesValue())

	empty := FromBytes([]byte{})
	assert.Equal(t, MT_BYTES, empty.Type())
	assert.Equal(t, 0, len(empty.AcquireBytesValue()))
}

func TestBytesStringIsHex(t *testing.T) {
	mv := FromBytes([]byte{0xff, 0x01})
	assert.Equal(t, "ff01", mv.String())

	assert.Equal(t, "", FromBytes([]byte{}).String())

	// StringMaybeQuoted: bytes are unquoted, like numbers
	assert.Equal(t, "ff01", FromBytes([]byte{0xff, 0x01}).StringMaybeQuoted())
}

func TestBytesJSON(t *testing.T) {
	mv := FromBytes([]byte{0xff, 0x01})
	s, err := mv.MarshalJSON(JSON_SINGLE_LINE, false)
	assert.Nil(t, err)
	assert.Equal(t, `"ff01"`, s)
}

func TestBytesCopyIsIndependent(t *testing.T) {
	original := []byte{1, 2, 3}
	mv := FromBytes(original)
	other := mv.Copy()
	original[0] = 99
	assert.Equal(t, []byte{99, 2, 3}, mv.AcquireBytesValue())
	assert.Equal(t, []byte{1, 2, 3}, other.AcquireBytesValue())
}

func TestBytesCmp(t *testing.T) {
	a := FromBytes([]byte{0x01})
	b := FromBytes([]byte{0x02})
	assert.Equal(t, -1, Cmp(a, b))
	assert.Equal(t, 1, Cmp(b, a))
	assert.Equal(t, 0, Cmp(a, FromBytes([]byte{0x01})))
	assert.True(t, Equals(a, FromBytes([]byte{0x01})))
	assert.False(t, Equals(a, b))

	// Mixed-type sort order: string < bytes < array
	assert.Equal(t, -1, Cmp(FromString("zzz"), a))
	assert.Equal(t, 1, Cmp(a, FromString("zzz")))
	assert.Equal(t, -1, Cmp(a, FromArray([]*Mlrval{})))
	assert.Equal(t, 1, Cmp(FromArray([]*Mlrval{}), a))
}
