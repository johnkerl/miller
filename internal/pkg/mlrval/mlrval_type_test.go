package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeNames(t *testing.T) {
	assert.Equal(t, "int", TYPE_NAMES[MT_INT])
	assert.Equal(t, "float", TYPE_NAMES[MT_FLOAT])
	assert.Equal(t, "bool", TYPE_NAMES[MT_BOOL])
	assert.Equal(t, "empty", TYPE_NAMES[MT_VOID])
	assert.Equal(t, "string", TYPE_NAMES[MT_STRING])
	assert.Equal(t, "array", TYPE_NAMES[MT_ARRAY])
	assert.Equal(t, "map", TYPE_NAMES[MT_MAP])
	assert.Equal(t, "funct", TYPE_NAMES[MT_FUNC])
	assert.Equal(t, "error", TYPE_NAMES[MT_ERROR])
	assert.Equal(t, "null", TYPE_NAMES[MT_NULL])
	assert.Equal(t, "absent", TYPE_NAMES[MT_ABSENT])
	assert.Equal(t, MT_DIM, len(TYPE_NAMES))
}
