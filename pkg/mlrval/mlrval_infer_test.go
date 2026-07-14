// Tests mlrval constructors.

package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferNormally(t *testing.T) {
	assert.True(t, inferNormally(FromDeferredType("")).IsVoid())

	assert.True(t, inferNormally(FromDeferredType("true")).IsString())
	assert.True(t, inferNormally(FromDeferredType("false")).IsString())

	assert.True(t, inferNormally(FromDeferredType("abc")).IsString())

	assert.True(t, inferNormally(FromDeferredType("0123")).IsString())
	assert.True(t, inferNormally(FromDeferredType("-0123")).IsString())
	assert.True(t, inferNormally(FromDeferredType("0377")).IsString())
	assert.True(t, inferNormally(FromDeferredType("-0377")).IsString())
	assert.True(t, inferNormally(FromDeferredType("0923")).IsString())
	assert.True(t, inferNormally(FromDeferredType("-0923")).IsString())

	assert.True(t, inferNormally(FromDeferredType("123")).IsInt())
	assert.True(t, inferNormally(FromDeferredType("-123")).IsInt())

	assert.True(t, inferNormally(FromDeferredType("0xff")).IsInt())
	assert.True(t, inferNormally(FromDeferredType("-0xff")).IsInt())
	assert.True(t, inferNormally(FromDeferredType("0b1011")).IsInt())
	assert.True(t, inferNormally(FromDeferredType("-0b1011")).IsInt())
	assert.True(t, inferNormally(FromDeferredType("0x7fffffffffffffff")).IsInt())
	assert.True(t, inferNormally(FromDeferredType("0x8000000000000000")).IsInt())
	assert.True(t, inferNormally(FromDeferredType("0xffffffffffffffff")).IsInt())

	assert.True(t, inferNormally(FromDeferredType("12_3")).IsString())
	assert.True(t, inferNormally(FromDeferredType("-12_3")).IsString())
	assert.True(t, inferNormally(FromDeferredType("1_2.3_4")).IsString())
	assert.True(t, inferNormally(FromDeferredType("-1_2.3_4")).IsString())
	assert.True(t, inferNormally(FromDeferredType("0xca_fe")).IsString())
	assert.True(t, inferNormally(FromDeferredType("-0xca_fe")).IsString())
	assert.True(t, inferNormally(FromDeferredType("0b1011_1101")).IsString())
	assert.True(t, inferNormally(FromDeferredType("-0b1011_1101")).IsString())

	assert.True(t, inferNormally(FromDeferredType(".")).IsString())
	assert.True(t, inferNormally(FromDeferredType("-.")).IsString())
	assert.True(t, inferNormally(FromDeferredType("123.")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("-123.")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType(".123")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("-.123")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("123.456")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("-123.456")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("1e2.")).IsString())
	assert.True(t, inferNormally(FromDeferredType("-1e2.")).IsString())
	assert.True(t, inferNormally(FromDeferredType("1e-2.")).IsString())
	assert.True(t, inferNormally(FromDeferredType("-1e-2.")).IsString())
	assert.True(t, inferNormally(FromDeferredType("1.2e3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("-1.2e3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("1.2e-3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("-1.2e-3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("1.e3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("-1.e3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("1.e-3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("-1.e-3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType(".2e3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("-.2e3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType(".2e-3")).IsFloat())
	assert.True(t, inferNormally(FromDeferredType("-.2e-3")).IsFloat())
}

func TestInferWithOctalAsInt(t *testing.T) {
	assert.True(t, inferWithOctalAsInt(FromDeferredType("")).IsVoid())

	assert.True(t, inferWithOctalAsInt(FromDeferredType("true")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("false")).IsString())

	assert.True(t, inferWithOctalAsInt(FromDeferredType("abc")).IsString())

	assert.True(t, inferWithOctalAsInt(FromDeferredType("0123")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-0123")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("0377")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-0377")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("0923")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-0923")).IsInt())

	assert.True(t, inferWithOctalAsInt(FromDeferredType("123")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-123")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("0xff")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-0xff")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("0b1011")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-0b1011")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("0x7fffffffffffffff")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("0x8000000000000000")).IsInt())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("0xffffffffffffffff")).IsInt())

	assert.True(t, inferWithOctalAsInt(FromDeferredType("12_3")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-12_3")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("1_2.3_4")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-1_2.3_4")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("0xca_fe")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-0xca_fe")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("0b1011_1101")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-0b1011_1101")).IsString())

	assert.True(t, inferWithOctalAsInt(FromDeferredType(".")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-.")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("123.")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-123.")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType(".123")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-.123")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("123.456")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-123.456")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("1e2.")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-1e2.")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("1e-2.")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-1e-2.")).IsString())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("1.2e3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-1.2e3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("1.2e-3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-1.2e-3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("1.e3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-1.e3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("1.e-3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-1.e-3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType(".2e3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-.2e3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType(".2e-3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-.2e-3")).IsFloat())
}

func TestInferWithIntAsFloat(t *testing.T) {
	assert.True(t, inferWithIntAsFloat(FromDeferredType("")).IsVoid())

	assert.True(t, inferWithIntAsFloat(FromDeferredType("true")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("false")).IsString())

	assert.True(t, inferWithIntAsFloat(FromDeferredType("abc")).IsString())

	assert.True(t, inferWithIntAsFloat(FromDeferredType("0123")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-0123")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("0377")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-0377")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("0923")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-0923")).IsString())

	assert.True(t, inferWithIntAsFloat(FromDeferredType("123")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-123")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("0xff")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-0xff")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("0b1011")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-0b1011")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("0x7fffffffffffffff")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("0x8000000000000000")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("0xffffffffffffffff")).IsFloat())

	assert.True(t, inferWithIntAsFloat(FromDeferredType("12_3")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-12_3")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("1_2.3_4")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-1_2.3_4")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("0xca_fe")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-0xca_fe")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("0b1011_1101")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-0b1011_1101")).IsString())

	assert.True(t, inferWithIntAsFloat(FromDeferredType(".")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-.")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("123.")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-123.")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType(".123")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-.123")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("123.456")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-123.456")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("1e2.")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-1e2.")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("1e-2.")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-1e-2.")).IsString())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("1.2e3")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-1.2e3")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("1.2e-3")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-1.2e-3")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("1.e3")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-1.e3")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("1.e-3")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-1.e-3")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType(".2e3")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("-.2e3")).IsFloat())
	assert.True(t, inferWithIntAsFloat(FromDeferredType(".2e-3")).IsFloat())
	assert.True(t, inferWithOctalAsInt(FromDeferredType("-.2e-3")).IsFloat())
}

func TestInferUserDefinedBooleans(t *testing.T) {
	// This test mutates package-level state; restore it on exit so other
	// tests in this package see default inference behavior.
	defer func() {
		inferrerBooleanTable = nil
	}()

	SetInferrerBooleanStrings([]string{"True", "yes", "on"}, true)
	SetInferrerBooleanStrings([]string{"False", "no", "off"}, false)

	mv := inferNormally(FromDeferredType("True"))
	assert.True(t, mv.IsBool())
	boolval, ok := mv.GetBoolValue()
	assert.True(t, ok)
	assert.True(t, boolval)
	// Original string representation is retained for output.
	assert.Equal(t, "True", mv.String())

	mv = inferNormally(FromDeferredType("no"))
	assert.True(t, mv.IsBool())
	boolval, ok = mv.GetBoolValue()
	assert.True(t, ok)
	assert.False(t, boolval)
	assert.Equal(t, "no", mv.String())

	assert.True(t, inferNormally(FromDeferredType("yes")).IsBool())
	assert.True(t, inferNormally(FromDeferredType("on")).IsBool())
	assert.True(t, inferNormally(FromDeferredType("off")).IsBool())

	// Matching is exact and case-sensitive; unlisted strings are unaffected.
	assert.True(t, inferNormally(FromDeferredType("true")).IsString())
	assert.True(t, inferNormally(FromDeferredType("TRUE")).IsString())
	assert.True(t, inferNormally(FromDeferredType("Yes")).IsString())
	assert.True(t, inferNormally(FromDeferredType("abc")).IsString())

	// Non-boolean inference is unaffected.
	assert.True(t, inferNormally(FromDeferredType("")).IsVoid())
	assert.True(t, inferNormally(FromDeferredType("123")).IsInt())
	assert.True(t, inferNormally(FromDeferredType("123.456")).IsFloat())

	// Also applies with -O and -A inference variants.
	assert.True(t, inferWithOctalAsInt(FromDeferredType("True")).IsBool())
	assert.True(t, inferWithIntAsFloat(FromDeferredType("no")).IsBool())

	// -S bypasses all inference, including user-defined booleans.
	assert.True(t, inferString(FromDeferredType("True")).IsString())
}

func TestInferString(t *testing.T) {
	assert.True(t, inferString(FromDeferredType("")).IsVoid())

	assert.True(t, inferString(FromDeferredType("true")).IsString())
	assert.True(t, inferString(FromDeferredType("false")).IsString())

	assert.True(t, inferString(FromDeferredType("abc")).IsString())

	assert.True(t, inferString(FromDeferredType("0123")).IsString())
	assert.True(t, inferString(FromDeferredType("-0123")).IsString())
	assert.True(t, inferString(FromDeferredType("0377")).IsString())
	assert.True(t, inferString(FromDeferredType("-0377")).IsString())
	assert.True(t, inferString(FromDeferredType("0923")).IsString())
	assert.True(t, inferString(FromDeferredType("-0923")).IsString())

	assert.True(t, inferString(FromDeferredType("123")).IsString())
	assert.True(t, inferString(FromDeferredType("-123")).IsString())
	assert.True(t, inferString(FromDeferredType("0xff")).IsString())
	assert.True(t, inferString(FromDeferredType("-0xff")).IsString())
	assert.True(t, inferString(FromDeferredType("0b1011")).IsString())
	assert.True(t, inferString(FromDeferredType("-0b1011")).IsString())
	assert.True(t, inferString(FromDeferredType("0x7fffffffffffffff")).IsString())
	assert.True(t, inferString(FromDeferredType("0x8000000000000000")).IsString())
	assert.True(t, inferString(FromDeferredType("0xffffffffffffffff")).IsString())

	assert.True(t, inferString(FromDeferredType("12_3")).IsString())
	assert.True(t, inferString(FromDeferredType("-12_3")).IsString())
	assert.True(t, inferString(FromDeferredType("1_2.3_4")).IsString())
	assert.True(t, inferString(FromDeferredType("-1_2.3_4")).IsString())
	assert.True(t, inferString(FromDeferredType("0xca_fe")).IsString())
	assert.True(t, inferString(FromDeferredType("-0xca_fe")).IsString())
	assert.True(t, inferString(FromDeferredType("0b1011_1101")).IsString())
	assert.True(t, inferString(FromDeferredType("-0b1011_1101")).IsString())

	assert.True(t, inferString(FromDeferredType(".")).IsString())
	assert.True(t, inferString(FromDeferredType("-.")).IsString())
	assert.True(t, inferString(FromDeferredType("123.")).IsString())
	assert.True(t, inferString(FromDeferredType("-123.")).IsString())
	assert.True(t, inferString(FromDeferredType(".123")).IsString())
	assert.True(t, inferString(FromDeferredType("-.123")).IsString())
	assert.True(t, inferString(FromDeferredType("123.456")).IsString())
	assert.True(t, inferString(FromDeferredType("-123.456")).IsString())
	assert.True(t, inferString(FromDeferredType("1e2.")).IsString())
	assert.True(t, inferString(FromDeferredType("-1e2.")).IsString())
	assert.True(t, inferString(FromDeferredType("1e-2.")).IsString())
	assert.True(t, inferString(FromDeferredType("-1e-2.")).IsString())
	assert.True(t, inferString(FromDeferredType("1.2e3")).IsString())
	assert.True(t, inferString(FromDeferredType("-1.2e3")).IsString())
	assert.True(t, inferString(FromDeferredType("1.2e-3")).IsString())
	assert.True(t, inferString(FromDeferredType("-1.2e-3")).IsString())
	assert.True(t, inferString(FromDeferredType("1.e3")).IsString())
	assert.True(t, inferString(FromDeferredType("-1.e3")).IsString())
	assert.True(t, inferString(FromDeferredType("1.e-3")).IsString())
	assert.True(t, inferString(FromDeferredType("-1.e-3")).IsString())
	assert.True(t, inferString(FromDeferredType(".2e3")).IsString())
	assert.True(t, inferString(FromDeferredType("-.2e3")).IsString())
	assert.True(t, inferString(FromDeferredType(".2e-3")).IsString())
	assert.True(t, inferString(FromDeferredType("-.2e-3")).IsString())
}
