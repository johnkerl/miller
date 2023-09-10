// ================================================================
// Tests mlrval typed-value extractors
// ================================================================

package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	mv := FromInferredType("234")
	stringval, ok := mv.GetStringValue()
	assert.False(t, ok)

	mv = FromDeferredType("234")
	stringval, ok = mv.GetStringValue()
	assert.False(t, ok)

	mv = FromInferredType("234.5")
	stringval, ok = mv.GetStringValue()
	assert.False(t, ok)

	mv = FromDeferredType("234.5")
	stringval, ok = mv.GetStringValue()
	assert.False(t, ok)

	mv = FromInferredType("abc")
	stringval, ok = mv.GetStringValue()
	assert.Equal(t, "abc", stringval)
	assert.True(t, ok)

	mv = FromDeferredType("abc")
	stringval, ok = mv.GetStringValue()
	assert.Equal(t, "abc", stringval)
	assert.True(t, ok)

	mv = FromInferredType("")
	stringval, ok = mv.GetStringValue()
	assert.Equal(t, "", stringval)
	assert.True(t, ok)

	mv = FromDeferredType("")
	stringval, ok = mv.GetStringValue()
	assert.Equal(t, "", stringval)
	assert.True(t, ok)
}

func TestGetIntValue(t *testing.T) {
	mv := FromInferredType("123")
	intval, ok := mv.GetIntValue()
	assert.Equal(t, int64(123), intval)
	assert.True(t, ok)

	mv = FromDeferredType("123")
	intval, ok = mv.GetIntValue()
	assert.Equal(t, int64(123), intval)
	assert.True(t, ok)

	mv = FromInferredType("123.4")
	intval, ok = mv.GetIntValue()
	assert.False(t, ok)

	mv = FromDeferredType("123.4")
	intval, ok = mv.GetIntValue()
	assert.False(t, ok)

	mv = FromInferredType("abc")
	intval, ok = mv.GetIntValue()
	assert.False(t, ok)

	mv = FromDeferredType("abc")
	intval, ok = mv.GetIntValue()
	assert.False(t, ok)
}

func TestGetFloatValue(t *testing.T) {
	mv := FromInferredType("234")
	floatval, ok := mv.GetFloatValue()
	assert.False(t, ok)

	mv = FromDeferredType("234")
	floatval, ok = mv.GetFloatValue()
	assert.False(t, ok)

	mv = FromInferredType("234.5")
	floatval, ok = mv.GetFloatValue()
	assert.Equal(t, 234.5, floatval)
	assert.True(t, ok)

	mv = FromDeferredType("234.5")
	floatval, ok = mv.GetFloatValue()
	assert.Equal(t, 234.5, floatval)
	assert.True(t, ok)

	mv = FromInferredType("abc")
	floatval, ok = mv.GetFloatValue()
	assert.False(t, ok)

	mv = FromDeferredType("abc")
	floatval, ok = mv.GetFloatValue()
	assert.False(t, ok)
}

func TestGetNumericToFloatValue(t *testing.T) {
	mv := FromInferredType("234")
	floatval, ok := mv.GetNumericToFloatValue()
	assert.Equal(t, 234.0, floatval)
	assert.True(t, ok)

	mv = FromDeferredType("234")
	floatval, ok = mv.GetNumericToFloatValue()
	assert.Equal(t, 234.0, floatval)
	assert.True(t, ok)

	mv = FromInferredType("234.5")
	floatval, ok = mv.GetNumericToFloatValue()
	assert.Equal(t, 234.5, floatval)
	assert.True(t, ok)

	mv = FromDeferredType("234.5")
	floatval, ok = mv.GetNumericToFloatValue()
	assert.Equal(t, 234.5, floatval)
	assert.True(t, ok)

	mv = FromInferredType("abc")
	floatval, ok = mv.GetNumericToFloatValue()
	assert.False(t, ok)

	mv = FromDeferredType("abc")
	floatval, ok = mv.GetNumericToFloatValue()
	assert.False(t, ok)
}

func TestGetBoolValue(t *testing.T) {
	mv := FromInferredType("234")
	boolval, ok := mv.GetBoolValue()
	assert.False(t, ok)

	mv = FromDeferredType("234")
	boolval, ok = mv.GetBoolValue()
	assert.False(t, ok)

	mv = FromInferredType("abc")
	boolval, ok = mv.GetBoolValue()
	assert.False(t, ok)

	mv = FromDeferredType("abc")
	boolval, ok = mv.GetBoolValue()
	assert.False(t, ok)

	mv = FromInferredType("true")
	boolval, ok = mv.GetBoolValue()
	assert.True(t, boolval)
	assert.True(t, ok)

	mv = FromDeferredType("false")
	boolval, ok = mv.GetBoolValue()
	assert.False(t, ok, "from-data-file \"false\" should infer to string")
}

func TestGetTypeName(t *testing.T) {
	mv := FromInferredType("234")
	assert.Equal(t, "int", mv.GetTypeName())

	mv = FromDeferredType("234")
	assert.Equal(t, "int", mv.GetTypeName())

	mv = FromInferredType("234.5")
	assert.Equal(t, "float", mv.GetTypeName())

	mv = FromDeferredType("234.5")
	assert.Equal(t, "float", mv.GetTypeName())

	mv = FromInferredType("abc")
	assert.Equal(t, "string", mv.GetTypeName())

	mv = FromDeferredType("abc")
	assert.Equal(t, "string", mv.GetTypeName())

	mv = FromInferredType("")
	assert.Equal(t, "empty", mv.GetTypeName())

	mv = FromDeferredType("")
	assert.Equal(t, "empty", mv.GetTypeName())
}
