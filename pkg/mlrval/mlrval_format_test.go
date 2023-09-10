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
