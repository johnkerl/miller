package scan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindScanTypeNameStrings(t *testing.T) {
	assert.Equal(t, typeNameString, findScanTypeName(""))
	assert.Equal(t, typeNameString, findScanTypeName("-"))
	assert.Equal(t, typeNameString, findScanTypeName("abc"))
	assert.Equal(t, typeNameString, findScanTypeName("-abc"))
}

func TestFindScanTypeNameDecimals(t *testing.T) {
	assert.Equal(t, typeNameDecimalInt, findScanTypeName("0"))
	assert.Equal(t, typeNameDecimalInt, findScanTypeName("-0"))
	assert.Equal(t, typeNameDecimalInt, findScanTypeName("1"))
	assert.Equal(t, typeNameDecimalInt, findScanTypeName("-1"))
	assert.Equal(t, typeNameDecimalInt, findScanTypeName("2"))
	assert.Equal(t, typeNameDecimalInt, findScanTypeName("-2"))
	assert.Equal(t, typeNameDecimalInt, findScanTypeName("123"))
	assert.Equal(t, typeNameDecimalInt, findScanTypeName("-123"))
}

func TestFindScanTypeNameFloats(t *testing.T) {
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("1."))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("-1."))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName(".2"))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("-.2"))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("-."))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("1.2"))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("-1.2"))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("12e-2"))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("-12e-2"))

	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("1.2.3"))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("-1.2.3"))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("1e2e3"))
	assert.Equal(t, typeNameMaybeFloat, findScanTypeName("-1e2e3"))

	assert.Equal(t, typeNameString, findScanTypeName("."))
	assert.Equal(t, typeNameString, findScanTypeName("1e2x3"))
	assert.Equal(t, typeNameString, findScanTypeName("-1e2x3"))

	assert.Equal(t, typeNameString, findScanTypeName("inf"))
	assert.Equal(t, typeNameString, findScanTypeName("infinity"))
	assert.Equal(t, typeNameString, findScanTypeName("NaN"))
	assert.Equal(t, typeNameString, findScanTypeName("-inf"))
	assert.Equal(t, typeNameString, findScanTypeName("-infinity"))
	assert.Equal(t, typeNameString, findScanTypeName("-NaN"))
}

func TestFindScanTypeNameHexes(t *testing.T) {
	assert.Equal(t, typeNameHexInt, findScanTypeName("0x0"))
	assert.Equal(t, typeNameHexInt, findScanTypeName("-0x0"))
	assert.Equal(t, typeNameHexInt, findScanTypeName("0xf"))
	assert.Equal(t, typeNameHexInt, findScanTypeName("-0xf"))
	assert.Equal(t, typeNameHexInt, findScanTypeName("0xcafe"))
	assert.Equal(t, typeNameHexInt, findScanTypeName("-0xcafe"))

	assert.Equal(t, typeNameHexInt, findScanTypeName("0x7ffffffffffffffe"))
	assert.Equal(t, typeNameHexInt, findScanTypeName("0x7fffffffffffffff"))
	assert.Equal(t, typeNameHexInt, findScanTypeName("0x8000000000000000"))
	assert.Equal(t, typeNameHexInt, findScanTypeName("0x8000000000000001"))
	assert.Equal(t, typeNameHexInt, findScanTypeName("0xfffffffffffffffe"))
	assert.Equal(t, typeNameHexInt, findScanTypeName("0xffffffffffffffff"))

	assert.Equal(t, typeNameString, findScanTypeName("0x"))
	assert.Equal(t, typeNameString, findScanTypeName("-0x"))
	assert.Equal(t, typeNameString, findScanTypeName("0xcape"))
	assert.Equal(t, typeNameString, findScanTypeName("-0xcape"))
}

func TestFindScanTypeNameOctals(t *testing.T) {
	assert.Equal(t, typeNameLeadingZeroOctalInt, findScanTypeName("00"))
	assert.Equal(t, typeNameLeadingZeroOctalInt, findScanTypeName("-00"))
	assert.Equal(t, typeNameLeadingZeroOctalInt, findScanTypeName("01"))
	assert.Equal(t, typeNameLeadingZeroOctalInt, findScanTypeName("-01"))
	assert.Equal(t, typeNameLeadingZeroOctalInt, findScanTypeName("0377"))
	assert.Equal(t, typeNameLeadingZeroOctalInt, findScanTypeName("-0377"))

	assert.Equal(t, typeNameLeadingZeroDecimalInt, findScanTypeName("08"))
	assert.Equal(t, typeNameLeadingZeroDecimalInt, findScanTypeName("-08"))

	assert.Equal(t, typeNameLeadingZeroDecimalInt, findScanTypeName("06789"))
	assert.Equal(t, typeNameLeadingZeroDecimalInt, findScanTypeName("-06789"))

	assert.Equal(t, typeNameOctalInt, findScanTypeName("0o377"))
	assert.Equal(t, typeNameOctalInt, findScanTypeName("-0o377"))

	assert.Equal(t, typeNameString, findScanTypeName("0o6789"))
	assert.Equal(t, typeNameString, findScanTypeName("-0o6789"))
}

func TestFindScanTypeNameBinaries(t *testing.T) {
	assert.Equal(t, typeNameBinaryInt, findScanTypeName("0b0"))
	assert.Equal(t, typeNameBinaryInt, findScanTypeName("-0b0"))
	assert.Equal(t, typeNameBinaryInt, findScanTypeName("0b1011"))
	assert.Equal(t, typeNameBinaryInt, findScanTypeName("-0b1011"))

	assert.Equal(t, typeNameString, findScanTypeName("0b"))
	assert.Equal(t, typeNameString, findScanTypeName("-0b"))
	assert.Equal(t, typeNameString, findScanTypeName("0b1021"))
	assert.Equal(t, typeNameString, findScanTypeName("-0b1021"))
}

func TestFindScanTypeNameBooleans(t *testing.T) {
	assert.Equal(t, typeNameString, findScanTypeName("true"))
	assert.Equal(t, typeNameString, findScanTypeName("True"))
	assert.Equal(t, typeNameString, findScanTypeName("false"))
	assert.Equal(t, typeNameString, findScanTypeName("False"))
}
