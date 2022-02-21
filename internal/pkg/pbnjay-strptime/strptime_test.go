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
		"%D %X", // no such format code
		false,
		0,
	},
}

func TestStrptime(t *testing.T) {
	for _, item := range testData {
		tval, err := Parse(item.input, item.format)
		if item.errNil {
			assert.Nil(t, err)
			seconds := tval.Unix()
			assert.Equal(t, seconds, item.output)

		} else {
			assert.NotNil(t, err)
		}
	}
}
