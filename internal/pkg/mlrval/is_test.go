// ================================================================
// Tests mlrval constructors.
// ================================================================

package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLegit(t *testing.T) {
	assert.False(t, MLRVAL_ERROR.IsLegit())
	assert.False(t, MLRVAL_ABSENT.IsLegit())
	assert.False(t, MLRVAL_NULL.IsLegit())
	assert.True(t, FromString("").IsLegit())
	assert.True(t, FromString("abc").IsLegit())
	assert.True(t, FromInt(123).IsLegit())
	assert.True(t, FromFloat(123.4).IsLegit())
	assert.True(t, FromBool(true).IsLegit())
	assert.True(t, FromArray("test data").IsLegit())
	assert.True(t, FromMap("test data").IsLegit())
}

func TestIsErrorOrAbsent(t *testing.T) {
	assert.True(t, MLRVAL_ERROR.IsErrorOrAbsent())
	assert.True(t, MLRVAL_ABSENT.IsErrorOrAbsent())
	assert.False(t, MLRVAL_NULL.IsErrorOrAbsent())
	assert.False(t, FromString("").IsErrorOrAbsent())
}

func TestIsError(t *testing.T) {
	assert.True(t, MLRVAL_ERROR.IsError())
	assert.False(t, MLRVAL_ABSENT.IsError())
	assert.False(t, MLRVAL_NULL.IsError())
	assert.False(t, FromString("").IsError())
}

func TestIsAbsent(t *testing.T) {
	assert.False(t, MLRVAL_ERROR.IsAbsent())
	assert.True(t, MLRVAL_ABSENT.IsAbsent())
	assert.False(t, MLRVAL_NULL.IsAbsent())
	assert.False(t, FromString("").IsAbsent())
}

func TestIsNull(t *testing.T) {
	assert.False(t, MLRVAL_ERROR.IsNull())
	assert.False(t, MLRVAL_ABSENT.IsNull())
	assert.True(t, MLRVAL_NULL.IsNull())
	assert.False(t, FromString("").IsNull())
}

func TestIsVoid(t *testing.T) {
	assert.False(t, MLRVAL_ERROR.IsVoid())
	assert.False(t, MLRVAL_ABSENT.IsVoid())
	assert.False(t, MLRVAL_NULL.IsVoid())
	assert.True(t, FromString("").IsVoid())
	assert.True(t, FromDeferredType("").IsVoid())
	assert.True(t, FromInferredType("").IsVoid())
	assert.False(t, FromDeferredType("abc").IsVoid())
	assert.False(t, FromInferredType("abc").IsVoid())
}

func TestIsEmptyString(t *testing.T) {
	assert.False(t, MLRVAL_ERROR.IsEmptyString())
	assert.False(t, MLRVAL_ABSENT.IsEmptyString())
	assert.False(t, MLRVAL_NULL.IsEmptyString())
	assert.True(t, FromString("").IsEmptyString())
	assert.True(t, FromDeferredType("").IsEmptyString())
	assert.True(t, FromInferredType("").IsEmptyString())
	assert.False(t, FromDeferredType("abc").IsEmptyString())
	assert.False(t, FromInferredType("abc").IsEmptyString())
}

func TestIsString(t *testing.T) {
	assert.False(t, MLRVAL_ERROR.IsString())
	assert.False(t, MLRVAL_ABSENT.IsString())
	assert.False(t, MLRVAL_NULL.IsString())
	assert.False(t, FromString("").IsString())
	assert.False(t, FromDeferredType("").IsString())
	assert.False(t, FromInferredType("").IsString())
	assert.True(t, FromDeferredType("abc").IsString())
	assert.True(t, FromInferredType("abc").IsString())
}

func TestIsStringOrVoid(t *testing.T) {
	assert.False(t, MLRVAL_ERROR.IsStringOrVoid())
	assert.False(t, MLRVAL_ABSENT.IsStringOrVoid())
	assert.False(t, MLRVAL_NULL.IsStringOrVoid())
	assert.True(t, FromString("").IsStringOrVoid())
	assert.True(t, FromDeferredType("").IsStringOrVoid())
	assert.True(t, FromInferredType("").IsStringOrVoid())
	assert.True(t, FromDeferredType("abc").IsStringOrVoid())
	assert.True(t, FromInferredType("abc").IsStringOrVoid())
}

func TestIsInt(t *testing.T) {
	assert.True(t, FromDeferredType("123").IsInt())
	assert.True(t, FromInferredType("123").IsInt())
	assert.False(t, FromDeferredType("123.4").IsInt())
	assert.False(t, FromInferredType("123.4").IsInt())
	assert.False(t, FromDeferredType("abc").IsInt())
	assert.False(t, FromInferredType("abc").IsInt())
}

func TestIsIntZero(t *testing.T) {
	assert.True(t, FromDeferredType("0").IsIntZero())
	assert.True(t, FromInferredType("0").IsIntZero())
	assert.True(t, FromDeferredType("-0").IsIntZero())
	assert.True(t, FromInferredType("-0").IsIntZero())
	assert.False(t, FromDeferredType("123").IsIntZero())
	assert.False(t, FromInferredType("123").IsIntZero())
	assert.False(t, FromDeferredType("abc").IsIntZero())
	assert.False(t, FromInferredType("abc").IsIntZero())
}

func TestIsFloat(t *testing.T) {
	assert.False(t, FromDeferredType("123").IsFloat())
	assert.False(t, FromInferredType("123").IsFloat())
	assert.True(t, FromDeferredType("123.4").IsFloat())
	assert.True(t, FromInferredType("123.4").IsFloat())
	assert.False(t, FromDeferredType("abc").IsFloat())
	assert.False(t, FromInferredType("abc").IsFloat())
}

func TestIsNumeric(t *testing.T) {
	assert.True(t, FromDeferredType("123").IsNumeric())
	assert.True(t, FromInferredType("123").IsNumeric())
	assert.True(t, FromDeferredType("123.4").IsNumeric())
	assert.True(t, FromInferredType("123.4").IsNumeric())
	assert.False(t, FromDeferredType("abc").IsNumeric())
	assert.False(t, FromInferredType("abc").IsNumeric())
}

func TestIsBool(t *testing.T) {
	assert.False(t, FromDeferredType("123").IsBool())
	assert.False(t, FromInferredType("123").IsBool())
	assert.False(t, FromDeferredType("123.4").IsBool())
	assert.False(t, FromInferredType("123.4").IsBool())
	assert.False(t, FromDeferredType("abc").IsBool())
	assert.False(t, FromInferredType("abc").IsBool())
	assert.False(t, FromDeferredType("true").IsBool(), "from-data-file \"true\" should infer to string")
	assert.True(t, FromDeferredType("true").IsString(), "from-data-file \"true\" should infer to string")
	assert.True(t, FromInferredType("true").IsBool())
	assert.False(t, FromDeferredType("false").IsBool(), "from-data-file \"false\" should infer to string")
	assert.True(t, FromDeferredType("false").IsString(), "from-data-file \"false\" should infer to string")
	assert.True(t, FromInferredType("false").IsBool())
}
