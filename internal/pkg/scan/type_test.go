package scan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeNames(t *testing.T) {
	assert.Equal(t, TypeNames[scanTypeString], "string")
	assert.Equal(t, TypeNames[scanTypeDecimalInt], "decint")
	assert.Equal(t, TypeNames[scanTypeOctalInt], "octint")
	assert.Equal(t, TypeNames[scanTypeHexInt], "hexint")
	assert.Equal(t, TypeNames[scanTypeBinaryInt], "binint")
	assert.Equal(t, TypeNames[scanTypeMaybeFloat], "float?")
	assert.Equal(t, TypeNames[scanTypeBool], "bool")
}
