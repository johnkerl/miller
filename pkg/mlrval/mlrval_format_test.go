package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatterToString(t *testing.T) {
	mv := FromString("hello")
	formatter := newFormatterToString("%s")
	fmv := formatter.Format(mv)
	assert.True(t, fmv.IsString())
	assert.Equal(t, "hello", fmv.String())

	mv = NULL
	formatter = newFormatterToString("%s")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsString())
	assert.Equal(t, "null", fmv.String())

	mv = TRUE
	formatter = newFormatterToString("%s")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsString())
	assert.Equal(t, "true", fmv.String())

	mv = FALSE
	formatter = newFormatterToString("%s")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsString())
	assert.Equal(t, "false", fmv.String())

	mv = FromInt(10)
	formatter = newFormatterToString("%s")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsString())
	assert.Equal(t, "10", fmv.String())

	mv = FromString("hello")
	formatter = newFormatterToString("[[%s]]")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsString())
	assert.Equal(t, "[[hello]]", fmv.String())
}

func TestFormatterToInt(t *testing.T) {
	mv := FromString("hello")
	formatter := newFormatterToInt("%d")
	fmv := formatter.Format(mv)
	assert.True(t, fmv.IsString())
	assert.Equal(t, "hello", fmv.String())

	mv = NULL
	formatter = newFormatterToInt("%d")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsNull())
	assert.Equal(t, "null", fmv.String())

	mv = TRUE
	formatter = newFormatterToInt("%d")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsBool())
	assert.Equal(t, "true", fmv.String())

	mv = FALSE
	formatter = newFormatterToInt("%d")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsBool())
	assert.Equal(t, "false", fmv.String())

	mv = FromInt(10)
	formatter = newFormatterToInt("%d")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsInt())
	assert.Equal(t, "10", fmv.String())

	mv = FromInt(10)
	formatter = newFormatterToInt("[[0x%x]]")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsString())
	assert.Equal(t, "[[0xa]]", fmv.String())

	mv = FromFloat(10.1)
	formatter = newFormatterToInt("%d")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsInt())
	assert.Equal(t, "10", fmv.String())
}
func TestFormatterToFloat(t *testing.T) {

	mv := FromString("hello")
	formatter := newFormatterToFloat("%.4f")
	fmv := formatter.Format(mv)
	assert.True(t, fmv.IsString())
	assert.Equal(t, "hello", fmv.String())

	mv = NULL
	formatter = newFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsNull())
	assert.Equal(t, "null", fmv.String())

	mv = TRUE
	formatter = newFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsBool())
	assert.Equal(t, "true", fmv.String())

	mv = FALSE
	formatter = newFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsBool())
	assert.Equal(t, "false", fmv.String())

	mv = FromFloat(10)
	formatter = newFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsFloat())
	assert.Equal(t, "10.0000", fmv.String())

	mv = FromFloat(10.1)
	formatter = newFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsFloat())
	assert.Equal(t, "10.1000", fmv.String())
}

func TestFormatter(t *testing.T) {

	mv := FromString("hello")
	formatter, err := GetFormatter("%d")
	assert.Nil(t, err)
	fmv := formatter.Format(mv)
	assert.True(t, fmv.IsString())
	assert.Equal(t, "hello", fmv.String())

	mv = NULL
	formatter, err = GetFormatter("%d")
	assert.Nil(t, err)
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsNull())
	assert.Equal(t, "null", fmv.String())

	mv = TRUE
	formatter, err = GetFormatter("%d")
	assert.Nil(t, err)
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsBool())
	assert.Equal(t, "true", fmv.String())

	mv = FALSE
	formatter, err = GetFormatter("%d")
	assert.Nil(t, err)
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsBool())
	assert.Equal(t, "false", fmv.String())

	mv = FromFloat(10.123)
	formatter, err = GetFormatter("%d")
	assert.Nil(t, err)
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsInt())
	assert.Equal(t, "10", fmv.String())

	mv = FromFloat(10.1)
	formatter, err = GetFormatter("%d")
	assert.Nil(t, err)
	fmv = formatter.Format(mv)
	assert.True(t, fmv.IsInt())
	assert.Equal(t, "10", fmv.String())
}

func TestFormatterVerbs(t *testing.T) {
	cases := []struct {
		format   string
		input    *Mlrval
		expected string
	}{
		{"%X", FromInt(255), "FF"},
		{"%08X", FromInt(255), "000000FF"},
		{"%o", FromInt(8), "10"},
		{"%b", FromInt(5), "101"},
		{"%E", FromFloat(2.5), "2.500000E+00"},
		{"%.3G", FromFloat(2.5), "2.5"},
		{"%F", FromFloat(2.5), "2.500000"},
		{"%d items", FromInt(17), "17 items"},
		{"[[0x%x]]", FromInt(10), "[[0xa]]"},
		{"[%.2f]", FromFloat(3.14159), "[3.14]"},
		{"%.3lf units", FromFloat(3.14159), "3.142 units"},
		{"%08lld", FromInt(17), "00000017"},
		{"%_d total", FromInt(1234567), "1,234,567 total"},
		{"%d world", FromInt(5), "5 world"},
		{"X%sX", FromString("hello"), "XhelloX"},
		{"%lx", FromInt(255), "ff"},
		{"%le", FromFloat(2.5), "2.500000e+00"},
		{"%.6_f", FromFloat(1234567.891), "1,234,567.891000"},
		{"%12_d", FromInt(1234567), "   1,234,567"},
	}
	for _, c := range cases {
		formatter, err := GetFormatter(c.format)
		assert.Nil(t, err)
		assert.Equal(t, c.expected, formatter.Format(c.input).String(), c.format)
	}
}

func TestFormatterFormatFloat(t *testing.T) {
	// The --ofmt path goes through FormatFloat rather than Format
	cases := []struct {
		format   string
		input    float64
		expected string
	}{
		{"%.3E", 2.5, "2.500E+00"},
		{"%.3lf", 3.14159, "3.142"},
		{"[%.2f]", 3.14159, "[3.14]"},
		{"%08X", 255.0, "000000FF"},
	}
	for _, c := range cases {
		formatter, err := GetFormatter(c.format)
		assert.Nil(t, err)
		assert.Equal(t, c.expected, formatter.FormatFloat(c.input), c.format)
	}
}
