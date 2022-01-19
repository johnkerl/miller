// ================================================================
// Tests mlrval constructors.
// ================================================================

package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromDeferredType(t *testing.T) {
	mv := FromDeferredType("123")
	assert.Equal(t, MT_PENDING, mv.mvtype)
	assert.Equal(t, "123", mv.printrep)
	assert.True(t, mv.printrepValid)

	mv = FromDeferredType("true")
	assert.Equal(t, MT_PENDING, mv.mvtype)
	assert.Equal(t, "true", mv.printrep)
	assert.True(t, mv.printrepValid)

	mv = FromDeferredType("abc")
	assert.Equal(t, MT_PENDING, mv.mvtype)
	assert.Equal(t, "abc", mv.printrep)
	assert.True(t, mv.printrepValid)

	mv = FromDeferredType("")
	assert.Equal(t, MT_PENDING, mv.mvtype)
	assert.Equal(t, "", mv.printrep)
	assert.True(t, mv.printrepValid)
}

func TestFromInferredType(t *testing.T) {
	mv := FromInferredType("123")
	assert.Equal(t, MT_INT, mv.mvtype)
	assert.Equal(t, "123", mv.printrep)
	assert.True(t, mv.printrepValid)
	assert.Equal(t, mv.intval, 123)

	mv = FromInferredType("true")
	assert.Equal(t, MT_BOOL, mv.mvtype)
	assert.Equal(t, "true", mv.printrep)
	assert.True(t, mv.printrepValid)
	assert.Equal(t, mv.boolval, true)

	mv = FromInferredType("abc")
	assert.Equal(t, MT_STRING, mv.mvtype)
	assert.Equal(t, "abc", mv.printrep)
	assert.True(t, mv.printrepValid)

	mv = FromInferredType("")
	assert.Equal(t, MT_VOID, mv.mvtype)
	assert.Equal(t, "", mv.printrep)
	assert.True(t, mv.printrepValid)
}

func TestFromString(t *testing.T) {
	mv := FromString("123")
	assert.Equal(t, MT_STRING, mv.mvtype)
	assert.Equal(t, "123", mv.printrep)
	assert.True(t, mv.printrepValid)

	mv = FromString("")
	assert.Equal(t, MT_VOID, mv.mvtype)
	assert.Equal(t, "", mv.printrep)
	assert.True(t, mv.printrepValid)
}

func TestFromInt(t *testing.T) {
	mv := FromInt(123)
	assert.Equal(t, MT_INT, mv.mvtype)
	assert.False(t, mv.printrepValid, "printrep should not be computed yet")
}

func TestTryFromIntString(t *testing.T) {
	mv := TryFromIntString("123")
	assert.Equal(t, MT_INT, mv.mvtype)
	assert.True(t, mv.printrepValid, "printrep should be computed")

	mv = TryFromIntString("[123]")
	assert.Equal(t, MT_STRING, mv.mvtype)
	assert.True(t, mv.printrepValid, "printrep should be computed")
}

func TestFromFloat(t *testing.T) {
	mv := FromFloat(123.4)
	assert.Equal(t, MT_FLOAT, mv.mvtype)
	assert.False(t, mv.printrepValid, "printrep should not be computed yet")
}

func TestTryFromFloatString(t *testing.T) {
	mv := TryFromFloatString("123.4")
	assert.Equal(t, MT_FLOAT, mv.mvtype)
	assert.True(t, mv.printrepValid, "printrep should be computed")

	mv = TryFromIntString("[123.4]")
	assert.Equal(t, MT_STRING, mv.mvtype)
	assert.True(t, mv.printrepValid, "printrep should be computed")
}

func TestFromBool(t *testing.T) {
	mv := FromBool(true)
	assert.Equal(t, MT_BOOL, mv.mvtype)
	assert.True(t, mv.printrepValid)

	mv = FromBool(false)
	assert.Equal(t, MT_BOOL, mv.mvtype)
	assert.True(t, mv.printrepValid)
}

func TestFromBoolString(t *testing.T) {
	mv := FromBoolString("true")
	assert.Equal(t, MT_BOOL, mv.mvtype)
	assert.True(t, mv.printrepValid)

	mv = FromBoolString("false")
	assert.Equal(t, MT_BOOL, mv.mvtype)
	assert.True(t, mv.printrepValid)
}

func TestFromFunction(t *testing.T) {
	mv := FromFunction("test data", "f001")
	assert.Equal(t, MT_FUNC, mv.mvtype)
	assert.True(t, mv.printrepValid)
	assert.Equal(t, "test data", mv.funcval.(string))
}

func TestFromArray(t *testing.T) {
	mv := FromArray([]*Mlrval{FromInt(10)})
	assert.Equal(t, MT_ARRAY, mv.mvtype)
	assert.Equal(t, 1, len(mv.arrayval))
}

func TestFromMap(t *testing.T) {
	mv := FromMap(NewMlrmap())
	assert.Equal(t, MT_MAP, mv.mvtype)
	assert.True(t, mv.mapval.IsEmpty())
}
