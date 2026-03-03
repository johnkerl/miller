package dkvpx

import (
	"strings"
	"testing"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/stretchr/testify/assert"
)

func TestWrite_Basic(t *testing.T) {
	var buf strings.Builder
	wr := NewWriter(&buf)
	rec := lib.NewOrderedMap[string]()
	rec.Put("x", "1")
	rec.Put("y", "2")
	rec.Put("z", "3")
	err := wr.Write(rec)
	assert.NoError(t, err)
	wr.Flush()
	assert.NoError(t, wr.Error())
	assert.Equal(t, "x=1,y=2,z=3\n", buf.String())
}

func TestWrite_QuotesWhenNeeded(t *testing.T) {
	var buf strings.Builder
	wr := NewWriter(&buf)
	rec := lib.NewOrderedMap[string]()
	rec.Put("x,y", "a,b,c")
	rec.Put("z", "3")
	err := wr.Write(rec)
	assert.NoError(t, err)
	wr.Flush()
	assert.Equal(t, `"x,y"="a,b,c",z=3`+"\n", buf.String())
}

func TestWrite_ValueWithEquals(t *testing.T) {
	var buf strings.Builder
	wr := NewWriter(&buf)
	rec := lib.NewOrderedMap[string]()
	rec.Put("a", "b=c")
	err := wr.Write(rec)
	assert.NoError(t, err)
	wr.Flush()
	assert.Equal(t, `a="b=c"`+"\n", buf.String())
}

func TestWrite_EscapedQuotes(t *testing.T) {
	var buf strings.Builder
	wr := NewWriter(&buf)
	rec := lib.NewOrderedMap[string]()
	rec.Put("x", `the "word"`)
	err := wr.Write(rec)
	assert.NoError(t, err)
	wr.Flush()
	assert.Equal(t, `x="the ""word"""`+"\n", buf.String())
}

func TestWrite_RoundTrip(t *testing.T) {
	input := `"x,y"="a,b,c",z=3` + "\n"
	rdr := NewReader(strings.NewReader(input))
	rec, err := rdr.Read()
	assert.NoError(t, err)

	var buf strings.Builder
	wr := NewWriter(&buf)
	err = wr.Write(rec)
	assert.NoError(t, err)
	wr.Flush()
	assert.Equal(t, input, buf.String())
}
