package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArrayifyArrayInputNotMutated(t *testing.T) {
	inner := NewMlrmap()
	inner.PutReference("1", FromString("a"))
	inner.PutReference("2", FromString("b"))
	input := FromArray([]*Mlrval{FromMap(inner)})
	out := input.Arrayify()
	b, err := out.FormatAsJSON(JSON_SINGLE_LINE, false)
	require.NoError(t, err)
	t.Logf("output=%s", string(b))
	b2, _ := input.FormatAsJSON(JSON_SINGLE_LINE, false)
	t.Logf("input after=%s", string(b2))
	assert.JSONEq(t, `[["a","b"]]`, string(b), "array elements that are convertible maps should be converted")
	assert.JSONEq(t, `[{"1":"a","2":"b"}]`, string(b2), "input must not be mutated")
}

func TestArrayifyMapInputNotMutated(t *testing.T) {
	inner := NewMlrmap()
	inner.PutReference("1", FromString("a"))
	inner.PutReference("2", FromString("b"))
	outer := NewMlrmap()
	outer.PutReference("s", FromMap(inner))
	input := FromMap(outer)
	out := input.Arrayify()
	b, _ := out.FormatAsJSON(JSON_SINGLE_LINE, false)
	b2, _ := input.FormatAsJSON(JSON_SINGLE_LINE, false)
	t.Logf("output=%s input after=%s", string(b), string(b2))
	assert.JSONEq(t, `{"s":{"1":"a","2":"b"}}`, string(b2), "input must not be mutated")
}

func TestArrayifyConvertibleTopLevelMapNotMutated(t *testing.T) {
	inner := NewMlrmap()
	inner.PutReference("1", FromString("a"))
	inner.PutReference("2", FromString("b"))
	outer := NewMlrmap()
	outer.PutReference("1", FromMap(inner))
	input := FromMap(outer)

	output := input.Arrayify()

	ob, _ := output.FormatAsJSON(JSON_SINGLE_LINE, false)
	ib, _ := input.FormatAsJSON(JSON_SINGLE_LINE, false)
	assert.JSONEq(t, `[["a","b"]]`, string(ob))
	assert.JSONEq(t, `{"1":{"1":"a","2":"b"}}`, string(ib),
		"Arrayify must not rewrite the input map when converting it to an array")
}
