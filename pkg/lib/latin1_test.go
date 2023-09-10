// ================================================================
// Most Miller tests (thousands of them) are command-line-driven via
// mlr regtest. Here are some cases needing special focus.
// ================================================================

package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type tDataForLatin1 struct {
	input          string
	expectedOutput string
	expectError    bool
}

var dataForLatin1ToUTF8 = []tDataForLatin1{
	{
		"",
		"",
		false,
	},
	{
		"The quick brown fox jumped over the lazy dogs.",
		"The quick brown fox jumped over the lazy dogs.",
		false,
	},
	{
		"a\xe4o\xf6",
		"a\xc3\xa4o\xc3\xb6", // "aäoö" -- showing explicitly here "\u00e4" encodes as "\xc3\xa4"
		false,
	},
	{
		"Victor jagt zw\xf6lf Boxk\xe4mpfer quer \xfcber den gro\xdfen Sylter Deich",
		"Victor jagt zwölf Boxkämpfer quer über den großen Sylter Deich",
		false,
	},
}

var dataForUTF8ToLatin1 = []tDataForLatin1{
	{
		"",
		"",
		false,
	},
	{
		"The quick brown fox jumped over the lazy dogs.",
		"The quick brown fox jumped over the lazy dogs.",
		false,
	},
	{
		"a\xc3\xa4o\xc3\xb6", // "aäoö" -- showing explicitly here "\u00e4" encodes as "\xc3\xa4"
		"a\xe4o\xf6",
		false,
	},
	{
		"Victor jagt zwölf Boxkämpfer quer über den großen Sylter Deich",
		"Victor jagt zw\xf6lf Boxk\xe4mpfer quer \xfcber den gro\xdfen Sylter Deich",
		false,
	},
	{
		"Съешь же ещё этих мягких французских булок да выпей чаю",
		"",
		true,
	},
}

func TestLatin1ToUTF8(t *testing.T) {
	for i, entry := range dataForLatin1ToUTF8 {
		actualOutput, err := TryLatin1ToUTF8(entry.input)
		if entry.expectError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" expected \"%s\" got \"%s\"\n",
				i, entry.input, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestUTF8ToLatin1(t *testing.T) {
	for i, entry := range dataForUTF8ToLatin1 {
		actualOutput, err := TryUTF8ToLatin1(entry.input)
		if entry.expectError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" expected \"%s\" got \"%s\"\n",
				i, entry.input, entry.expectedOutput, actualOutput,
			)
		}
	}
}
