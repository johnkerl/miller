package bifs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

// TestBIF_strftime_strptime_roundtrip checks that for each of a set of formats
// exercising the newly-added %a %A %e %h %x %X %c strptime codes,
// strptime(strftime(t, fmt), fmt) == t. This is the exact failure mode from
// https://github.com/johnkerl/miller/issues/1518, where strptime accepted fewer
// format codes than strftime produced.
func TestBIF_strftime_strptime_roundtrip(t *testing.T) {
	formats := []string{
		"%a %b %e %T %Y",
		"%A %b %e %T %Y",
		"%a %h %e %T %Y",
		"%A %h %e %T %Y",
		"%c",
		"%x %X",
		"%Y-%m-%d %H:%M:%S",
	}

	epochSecondsValues := []int64{
		0,          // 1970-01-01, single-digit day
		1709552468, // 2024-03-04 11:41:08, single-digit day
		1729555200, // 2024-10-22 00:00:00, double-digit day
		1735732805, // 2025-01-01 12:00:05, single-digit day, new year
	}

	for _, format := range formats {
		for _, epochSeconds := range epochSecondsValues {
			formatted := BIF_strftime(mlrval.FromInt(epochSeconds), mlrval.FromString(format))
			formattedString, isString := formatted.GetStringValue()
			assert.True(t, isString, "strftime(%d, %q) produced non-string %v", epochSeconds, format, formatted)

			parsed := BIF_strptime(mlrval.FromString(formattedString), mlrval.FromString(format))
			parsedSeconds, isNumeric := parsed.GetNumericToFloatValue()
			assert.True(t, isNumeric, "strptime(%q, %q) produced non-numeric %v", formattedString, format, parsed)

			assert.Equal(
				t, float64(epochSeconds), parsedSeconds,
				"round-trip mismatch for format %q: strftime(%d) = %q, strptime(...) = %v",
				format, epochSeconds, formattedString, parsedSeconds,
			)
		}
	}
}

// TestBIF_strptime_new_format_codes exercises the newly-supported %a %A %e %h
// %x %X %c strptime format codes directly, including %e's space-padding edge
// cases (single-digit day with/without padding, double-digit day, %e at the
// end of the string, and %e directly adjacent to another format code).
func TestBIF_strptime_new_format_codes(t *testing.T) {
	type testCase struct {
		input  string
		format string
		want   float64
	}

	cases := []testCase{
		{"Mon Mar  4 11:41:08 2024", "%a %b %e %T %Y", 1709552468},
		{"Monday Mar  4 11:41:08 2024", "%A %b %e %T %Y", 1709552468},
		{"Mon Mar  4 11:41:08 2024", "%a %h %e %T %Y", 1709552468},
		{"Tue Oct 22 00:00:00 2024", "%a %b %e %T %Y", 1729555200},
		{"Mon Mar  4 11:41:08 2024", "%c", 1709552468},
		{"03/04/24 11:41:08", "%x %X", 1709552468},
		// %e with a single (non-padded) space before a single-digit day.
		{"Mar 4 2024", "%b %e %Y", 1709510400},
		// %e with two-digit day.
		{"Mar 14 2024", "%b %e %Y", 1710374400},
		// %e directly adjacent to another format code, no intervening text.
		{"4Mar2024", "%e%b%Y", 1709510400},
		{"14Mar2024", "%e%b%Y", 1710374400},
		// %e at the very end of the string.
		{"2024-03-4", "%Y-%m-%e", 1709510400},
		{"2024-03-14", "%Y-%m-%e", 1710374400},
	}

	for _, c := range cases {
		output := BIF_strptime(mlrval.FromString(c.input), mlrval.FromString(c.format))
		seconds, isNumeric := output.GetNumericToFloatValue()
		assert.True(t, isNumeric, "strptime(%q, %q) produced non-numeric %v", c.input, c.format, output)
		assert.Equal(t, c.want, seconds, "strptime(%q, %q)", c.input, c.format)
	}
}

// TestBIF_strptime_still_unsupported_format_code checks that a format code
// which remains unsupported (%U, week-of-year) still errors, i.e. that adding
// the new codes didn't accidentally cause unsupported codes to be silently
// ignored.
func TestBIF_strptime_still_unsupported_format_code(t *testing.T) {
	output := BIF_strptime(mlrval.FromString("2024-03-04 09"), mlrval.FromString("%Y-%m-%d %U"))
	assert.True(t, output.IsError())
}
