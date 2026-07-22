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

// TestBIF_datediff exercises the spreadsheet-DATEDIF-style datediff function
// (https://github.com/johnkerl/miller/issues/708), including the examples from
// the Excel DATEDIF documentation, leap-year and month-end edge cases, sign
// behavior for reversed arguments, and case-insensitivity of the unit.
func TestBIF_datediff(t *testing.T) {
	ymd := func(dateString string) *mlrval.Mlrval {
		return BIF_strptime(mlrval.FromString(dateString), mlrval.FromString("%Y-%m-%d"))
	}

	type testCase struct {
		start string
		end   string
		unit  string
		want  int64
	}

	cases := []testCase{
		// Values as documented for Excel's DATEDIF.
		{"2001-01-01", "2003-01-01", "y", 2},
		{"2001-06-01", "2002-08-15", "d", 440},
		{"2001-06-01", "2002-08-15", "yd", 75},
		{"2001-06-01", "2002-08-15", "md", 14},
		{"2001-06-01", "2002-08-15", "m", 14},
		{"2001-06-01", "2002-08-15", "ym", 2},
		{"2001-06-01", "2002-08-15", "y", 1},

		// Case-insensitive unit.
		{"2001-06-01", "2002-08-15", "YD", 75},
		{"2001-06-01", "2002-08-15", "Y", 1},

		// Same date and same month-and-day.
		{"2020-05-15", "2020-05-15", "d", 0},
		{"2020-05-15", "2020-05-15", "y", 0},
		{"2019-05-15", "2020-05-15", "y", 1},
		{"2019-05-15", "2020-05-15", "yd", 0},

		// Reversed arguments negate.
		{"2002-08-15", "2001-06-01", "d", -440},
		{"2002-08-15", "2001-06-01", "m", -14},
		{"2003-01-01", "2001-01-01", "y", -2},

		// Leap-year edges: Feb 29 start, non-leap end year.
		{"2020-02-29", "2021-02-28", "y", 0},
		{"2020-02-29", "2021-03-01", "y", 1},
		{"2020-02-29", "2021-02-28", "d", 365},

		// Month-end borrow: complete months don't count until the day is
		// reached; "md" can go negative as in the spreadsheet versions.
		{"2020-01-31", "2020-03-01", "m", 1},
		{"2020-01-31", "2020-03-01", "md", -1},
		{"2019-01-31", "2019-03-01", "md", -2},
		{"2020-01-31", "2020-02-29", "m", 0},

		// Year boundary for "ym"/"yd" when end month-day precedes start
		// month-day.
		{"2020-11-15", "2021-02-10", "m", 2},
		{"2020-11-15", "2021-02-10", "ym", 2},
		{"2020-11-15", "2021-02-10", "yd", 87},
	}

	for _, c := range cases {
		output := BIF_datediff(ymd(c.start), ymd(c.end), mlrval.FromString(c.unit))
		value, isInt := output.GetIntValue()
		assert.True(t, isInt, "datediff(%q, %q, %q) produced non-int %v", c.start, c.end, c.unit, output)
		assert.Equal(t, c.want, value, "datediff(%q, %q, %q)", c.start, c.end, c.unit)
	}

	// Time-of-day parts are ignored: one second before midnight vs. one
	// second after is still one calendar day.
	output := BIF_datediff(
		mlrval.FromInt(86399), // 1970-01-01T23:59:59Z
		mlrval.FromInt(86401), // 1970-01-02T00:00:01Z
		mlrval.FromString("d"),
	)
	value, isInt := output.GetIntValue()
	assert.True(t, isInt)
	assert.Equal(t, int64(1), value)

	// Unknown unit and non-numeric inputs are errors.
	assert.True(t, BIF_datediff(ymd("2020-01-01"), ymd("2020-01-02"), mlrval.FromString("w")).IsError())
	assert.True(t, BIF_datediff(mlrval.FromString("not-a-date"), ymd("2020-01-02"), mlrval.FromString("d")).IsError())
}

// TestBIF_strptime_still_unsupported_format_code checks that a format code
// which remains unsupported (%U, week-of-year) still errors, i.e. that adding
// the new codes didn't accidentally cause unsupported codes to be silently
// ignored.
func TestBIF_strptime_still_unsupported_format_code(t *testing.T) {
	output := BIF_strptime(mlrval.FromString("2024-03-04 09"), mlrval.FromString("%Y-%m-%d %U"))
	assert.True(t, output.IsError())
}
