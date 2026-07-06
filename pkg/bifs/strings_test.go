package bifs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func TestBIF_format(t *testing.T) {
	// No arguments at all
	output := BIF_format([]*mlrval.Mlrval{})
	assert.True(t, output.IsVoid())

	// Non-string format string
	output = BIF_format([]*mlrval.Mlrval{mlrval.FromInt(1)})
	assert.True(t, output.IsError())

	// Sequential placeholders
	cases := []struct {
		formatString string
		args         []string
		expected     string
	}{
		// Plain "{}" placeholders (pre-existing behavior)
		{"", []string{}, ""},
		{"abc", []string{}, "abc"},
		{"{}", []string{}, ""},
		{"{}", []string{"1"}, "1"},
		{"{}", []string{"1", "2"}, "1"},
		{"{}:{}:{}", []string{"1", "2"}, "1:2:"},
		{"{}:{}:{}", []string{"1", "2", "3"}, "1:2:3"},
		{"{}:{}:{}", []string{"1", "2", "3", "4"}, "1:2:3"},
		{"<{}:{}>", []string{"abc"}, "<abc:>"},
		{"<{}:{}>", []string{"abc", "def"}, "<abc:def>"},

		// Positional "{N}" placeholders
		{"{1}", []string{"a"}, "a"},
		{"{1}:{1}", []string{"a"}, "a:a"},
		{"{2}:{1}", []string{"a", "b"}, "b:a"},
		{"{1}/{2}/{1}_{3}.ext", []string{"a", "b", "c"}, "a/b/a_c.ext"},
		// Out-of-range positional index interpolates the empty string,
		// consistent with too-few arguments for "{}"
		{"{3}", []string{"a"}, ""},
		{"<{4}>", []string{"a", "b"}, "<>"},
		// Leading zeros are accepted
		{"{01}", []string{"a"}, "a"},

		// Mixing: "{}" consumes arguments sequentially, independently of
		// positional placeholders
		{"{2}{}:{1}{}", []string{"3", "4"}, "43:34"},
		{"{}{1}{}", []string{"a", "b"}, "aab"},

		// Non-placeholder braces are left alone
		{"{ }", []string{"a"}, "{ }"},
		{"{x}", []string{"a"}, "{x}"},
		{"{", []string{"a"}, "{"},
		{"}", []string{"a"}, "}"},
	}

	for _, c := range cases {
		mlrvals := []*mlrval.Mlrval{mlrval.FromString(c.formatString)}
		for _, arg := range c.args {
			mlrvals = append(mlrvals, mlrval.FromString(arg))
		}
		output := BIF_format(mlrvals)
		stringval, ok := output.GetStringValue()
		assert.True(t, ok, "format(%q, %v)", c.formatString, c.args)
		assert.Equal(t, c.expected, stringval, "format(%q, %v)", c.formatString, c.args)
	}

	// "{0}" is an error, since positional placeholders are 1-based
	output = BIF_format([]*mlrval.Mlrval{mlrval.FromString("{0}"), mlrval.FromString("a")})
	assert.True(t, output.IsError())
}
