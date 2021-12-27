package scan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeNames(t *testing.T) {
	assert.Equal(t, TypeNames[scanTypeString], "string")
	assert.Equal(t, TypeNames[scanTypeDecimalInt], "decint")
	assert.Equal(t, TypeNames[scanTypeLeadingZeroDecimalInt], "lzdecint") // e.g. 0899
	assert.Equal(t, TypeNames[scanTypeOctalInt], "octint") // e.g. 0o377
	assert.Equal(t, TypeNames[scanTypeLeadingZeroOctalInt], "lzoctint") // e.g. 0377
	assert.Equal(t, TypeNames[scanTypeHexInt], "hexint") // e.g. 0xcafe
	assert.Equal(t, TypeNames[scanTypeBinaryInt], "binint") // e.g. 0b1011
	assert.Equal(t, TypeNames[scanTypeMaybeFloat], "float?") // characters in [0-9\.-+eE] but needs parse to be sure
}
