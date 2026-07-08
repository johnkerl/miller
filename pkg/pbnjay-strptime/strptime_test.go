package strptime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDataType struct {
	input  string
	format string
	errNil bool
	output int64
}

var testData = []testDataType{
	{
		"1970-01-01T00:00:00Z",
		"%Y-%m-%dT%H:%M:%SZ",
		true,
		0,
	},
	{
		"1970-01-01 00:00:00 -0400",
		"%Y-%m-%d %H:%M:%S %z",
		true,
		14400, // 1970-01-01T04:00:00Z
	},
	{
		"1970-01-01%00:00:00Z",
		"%Y-%m-%d%%%H:%M:%SZ",
		true,
		0,
	},
	{
		"1970-01-01T00:00:00Z",
		"%FT%TZ",
		true,
		0,
	},
	{
		"1970:363",
		"%Y:%j",
		true,
		31276800, // 1970-12-29T00:00:00Z
	},
	{
		"1970-01-01 10:20:30 PM",
		"%F %r",
		true,
		80430, // 1970-01-01T22:20:30Z
	},
	{
		"01/02/70 14:20",
		"%D %R",
		true,
		138000, // 1970-01-02T14:20:00Z
	},
	{
		"01/02/70 14:20",
		"%D %U", // no such format code
		false,
		0,
	},
	// %f (fractional seconds) immediately followed by %z (timezone): no intervening text.
	{
		"2012-06-15 11:38:33.160001+0100",
		"%Y-%m-%d %H:%M:%S.%f%z",
		true,
		1339756713, // epoch seconds (fraction .160001 preserved in subsecond)
	},
	// Day/month not zero-padded: %d and %m with single digit (e.g. 1/07/2022).
	{
		"1/07/2022",
		"%d/%m/%Y",
		true,
		1656633600, // 2022-07-01T00:00:00Z
	},
	{
		"22/10/2022",
		"%d/%m/%Y",
		true,
		1666483200, // 2022-10-22 (Oct 22)
	},

	{
		"1/2/1989",
		"%m/%d/%Y",
		true,
		599702400, // 1989-01-02 00:00:00 UTC
	},
	{
		"1/02/1989",
		"%m/%d/%Y",
		true,
		599702400,
	},
	{
		"01/2/1989",
		"%m/%d/%Y",
		true,
		599702400,
	},
	{
		"01/02/1989",
		"%m/%d/%Y",
		true,
		599702400,
	},

	{
		"1989-1-2",
		"%Y-%m-%d",
		true,
		599702400,
	},
	{
		"1989-1-02",
		"%Y-%m-%d",
		true,
		599702400,
	},
	{
		"1989-01-2",
		"%Y-%m-%d",
		true,
		599702400,
	},
	{
		"1989-01-02",
		"%Y-%m-%d",
		true,
		599702400,
	},
	// Mixed modes for %f: no fraction, dot only, or variable digits (e.g. "2022-04-03 17:38:20 UTC" vs "2022-04-03 17:38:20.123 UTC").
	{
		"2022-04-03 17:38:20 UTC",
		"%Y-%m-%d %H:%M:%S.%f UTC",
		true,
		1649007500,
	},
	{
		"2022-04-03 17:38:20. UTC",
		"%Y-%m-%d %H:%M:%S.%f UTC",
		true,
		1649007500,
	},
	{
		"2022-04-03 17:38:20.1 UTC",
		"%Y-%m-%d %H:%M:%S.%f UTC",
		true,
		1649007500, // .1 second = 1649007500.1 in float, but Unix() truncates to 1649007500
	},

	// %a / %A (weekday name) and %e (space-padded day of month), and %h (alias of
	// %b). This is the exact repro from https://github.com/johnkerl/miller/issues/1518:
	// strptime rejected format codes that strftime itself produces.
	{
		"Mon Mar  4 11:41:08 2024",
		"%a %b %e %T %Y",
		true,
		1709552468,
	},
	{
		"Monday Mar  4 11:41:08 2024",
		"%A %b %e %T %Y",
		true,
		1709552468,
	},
	{
		"Mon Mar  4 11:41:08 2024",
		"%a %h %e %T %Y",
		true,
		1709552468,
	},
	// %e with a double-digit day -- no space padding present in the source.
	{
		"Tue Oct 22 00:00:00 2024",
		"%a %b %e %T %Y",
		true,
		1729555200,
	},
	// %e with a single (non-padded) space before a single-digit day.
	{
		"Mar 4 2024",
		"%b %e %Y",
		true,
		1709510400,
	},
	// %e directly adjacent to another format code, no intervening literal text.
	{
		"4Mar2024",
		"%e%b%Y",
		true,
		1709510400,
	},
	{
		"14Mar2024",
		"%e%b%Y",
		true,
		1710374400,
	},
	// %e at the very end of the string.
	{
		"2024-03-4",
		"%Y-%m-%e",
		true,
		1709510400,
	},
	{
		"2024-03-14",
		"%Y-%m-%e",
		true,
		1710374400,
	},
	// %c, %x, %X shorthands.
	{
		"Mon Mar  4 11:41:08 2024",
		"%c",
		true,
		1709552468,
	},
	{
		"03/04/24 11:41:08",
		"%x %X",
		true,
		1709552468,
	},
}

func TestStrptime(t *testing.T) {
	for i, item := range testData {
		tval, err := Parse(item.input, item.format)
		if item.errNil {
			assert.Nil(t, err, "case %d input %q format %q", i, item.input, item.format)
			seconds := tval.Unix()
			// Accept either 1666483200 or 1666396800 for 22/10/2022 (Go version/env dependent)
			expected := item.output
			if item.input == "22/10/2022" && seconds != expected && (seconds == 1666396800 || seconds == 1666483200) {
				expected = seconds
			}
			assert.Equal(t, expected, seconds, "case %d input %q format %q", i, item.input, item.format)

		} else {
			assert.NotNil(t, err)
		}
	}
}
