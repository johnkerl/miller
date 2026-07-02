package dkvpx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatField_Basic(t *testing.T) {
	assert.Equal(t, "x", FormatField("x"))
	assert.Equal(t, "1", FormatField("1"))
	assert.Equal(t, "abc def", FormatField("abc def"))
}

func TestFormatField_QuotesWhenNeeded(t *testing.T) {
	assert.Equal(t, `"x,y"`, FormatField("x,y"))
	assert.Equal(t, `"a,b,c"`, FormatField("a,b,c"))
}

func TestFormatField_ValueWithEquals(t *testing.T) {
	assert.Equal(t, `"b=c"`, FormatField("b=c"))
}

func TestFormatField_EscapedQuotes(t *testing.T) {
	assert.Equal(t, `"the ""word"""`, FormatField(`the "word"`))
}

func TestFormatField_RoundTrip(t *testing.T) {
	input := `"x,y"="a,b,c",z=3` + "\n"
	rdr := NewReader(strings.NewReader(input))
	rec, err := rdr.Read()
	assert.NoError(t, err)

	var sb strings.Builder
	first := true
	for pe := rec.Head; pe != nil; pe = pe.Next {
		if !first {
			sb.WriteByte(',')
		}
		first = false
		sb.WriteString(FormatField(pe.Key))
		sb.WriteByte('=')
		sb.WriteString(FormatField(pe.Value))
	}
	sb.WriteByte('\n')
	assert.Equal(t, input, sb.String())
}
