package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidJSONNumber(t *testing.T) {
	validCases := []string{
		"0",
		"-0",
		"1",
		"-1",
		"123",
		"-123",
		"0.5",
		"-0.5",
		"4.56",
		"-4.56",
		"123.456",
		"1e2",
		"1E2",
		"1e+2",
		"1e-2",
		"0.5e2",
		"-0.5E-2",
		"1234567890123456789",
	}
	for _, input := range validCases {
		assert.True(t, isValidJSONNumber(input), "expected valid JSON number: %q", input)
	}

	invalidCases := []string{
		"",
		"-",
		"+",
		"+1",       // leading plus
		"+4.56",    // leading plus
		"004.56",   // leading zeros
		"00.56",    // leading zeros
		"-004.56",  // leading zeros
		"007.5",    // leading zeros
		"01",       // leading zero
		".56",      // no digit before decimal point
		"-.5",      // no digit before decimal point
		"4.",       // no digit after decimal point
		"1.e3",     // no digit after decimal point
		".5e2",     // no digit before decimal point
		"00.5e2",   // leading zeros
		"1e",       // no exponent digits
		"1e+",      // no exponent digits
		"1.2.3",    // multiple decimal points
		"0x1f",     // hex is not JSON
		"0b101",    // binary is not JSON
		"NaN",      // not in the JSON grammar
		"+Inf",     // not in the JSON grammar
		"-Inf",     // not in the JSON grammar
		"Infinity", // not in the JSON grammar
		"abc",
		"1abc",
	}
	for _, input := range invalidCases {
		assert.False(t, isValidJSONNumber(input), "expected invalid JSON number: %q", input)
	}
}

// TestMarshalJSONFloatValidity checks that float-typed mlrvals whose original
// string representations are not within the JSON number grammar -- e.g.
// leading zeros like 004.56 -- are re-rendered as valid JSON numbers.
// See https://github.com/johnkerl/miller/issues/1114
// and https://github.com/johnkerl/miller/issues/1293.
func TestMarshalJSONFloatValidity(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		// Invalid-for-JSON originals are re-rendered:
		{"004.56", "4.56"},
		{"-004.56", "-4.56"},
		{"+4.56", "4.56"},
		{".56", "0.56"},
		{"4.", "4"},
		{"1.e3", "1000"},
		{".5e2", "50"},
		{"00.5e2", "50"},
		// Valid-for-JSON originals are passed through as-is:
		{"4.56", "4.56"},
		{"0.5", "0.5"},
		{"1e-2", "1e-2"},
		{"1E3", "1E3"},
		{"123.456", "123.456"},
	}
	for _, c := range cases {
		mv := FromInferredType(c.input)
		assert.True(t, mv.IsFloat(), "expected float inference for %q", c.input)
		actual, err := mv.marshalJSONFloat(false)
		assert.Nil(t, err)
		assert.Equal(t, c.expected, actual, "JSON marshal of %q", c.input)
	}
}
