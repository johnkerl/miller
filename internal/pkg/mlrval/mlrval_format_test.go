package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	assert.Equal(t, true, true)
}

func TestFormatterToString(t *testing.T) {
	mv := FromString("hello")
	formatter := newFormatterToString("%s")
	fmv := formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "hello" {
		t.Fatal()
	}

	mv = NULL
	formatter = newFormatterToString("%s")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "null" {
		t.Fatal()
	}

	mv = TRUE
	formatter = newFormatterToString("%s")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "true" {
		t.Fatal()
	}

	mv = FALSE
	formatter = newFormatterToString("%s")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "false" {
		t.Fatal()
	}

	mv = FromInt(10)
	formatter = newFormatterToString("%s")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "10" {
		t.Fatal()
	}

	mv = FromString("hello")
	formatter = newFormatterToString("[[%s]]")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "[[hello]]" {
		t.Fatal()
	}

}

func TestFormatterToInt(t *testing.T) {
	mv := FromString("hello")
	formatter := newFormatterToInt("%d")
	fmv := formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "hello" {
		t.Fatal()
	}

	mv = NULL
	formatter = newFormatterToInt("%d")
	fmv = formatter.Format(mv)
	if !fmv.IsNull() {
		t.Fatal()
	}
	if fmv.String() != "null" {
		t.Fatal()
	}

	mv = TRUE
	formatter = newFormatterToInt("%d")
	fmv = formatter.Format(mv)
	if !fmv.IsBool() {
		t.Fatal()
	}
	if fmv.String() != "true" {
		t.Fatal()
	}

	mv = FALSE
	formatter = newFormatterToInt("%d")
	fmv = formatter.Format(mv)
	if !fmv.IsBool() {
		t.Fatal()
	}
	if fmv.String() != "false" {
		t.Fatal()
	}

	mv = FromInt(10)
	formatter = newFormatterToInt("%d")
	fmv = formatter.Format(mv)
	if !fmv.IsInt() {
		t.Fatal()
	}
	if fmv.String() != "10" {
		t.Fatal()
	}

	mv = FromInt(10)
	formatter = newFormatterToInt("[[0x%x]]")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "[[0xa]]" {
		t.Fatal()
	}

	mv = FromFloat(10.1)
	formatter = newFormatterToInt("%d")
	fmv = formatter.Format(mv)
	if !fmv.IsInt() {
		t.Fatal()
	}
	if fmv.String() != "10" {
		t.Fatal()
	}

}

func TestFormatterToFloat(t *testing.T) {

	mv := FromString("hello")
	formatter := newFormatterToFloat("%.4f")
	fmv := formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "hello" {
		t.Fatal()
	}

	mv = NULL
	formatter = newFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	if !fmv.IsNull() {
		t.Fatal()
	}
	if fmv.String() != "null" {
		t.Fatal()
	}

	mv = TRUE
	formatter = newFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	if !fmv.IsBool() {
		t.Fatal()
	}
	if fmv.String() != "true" {
		t.Fatal()
	}

	mv = FALSE
	formatter = newFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	if !fmv.IsBool() {
		t.Fatal()
	}
	if fmv.String() != "false" {
		t.Fatal()
	}

	mv = FromFloat(10)
	formatter = newFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	if !fmv.IsFloat() {
		t.Fatal()
	}
	if fmv.String() != "10.0000" {
		t.Fatal()
	}

	mv = FromFloat(10.1)
	formatter = newFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	if !fmv.IsFloat() {
		t.Fatal()
	}
	if fmv.String() != "10.1000" {
		t.Fatal()
	}
}

func TestFormatter(t *testing.T) {

	mv := FromString("hello")
	formatter, err := GetFormatter("%d")
	if err != nil {
		t.Fatal()
	}
	fmv := formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "hello" {
		t.Fatal()
	}

	mv = NULL
	formatter, err = GetFormatter("%d")
	if err != nil {
		t.Fatal()
	}
	fmv = formatter.Format(mv)
	if !fmv.IsNull() {
		t.Fatal()
	}
	if fmv.String() != "null" {
		t.Fatal()
	}

	mv = TRUE
	formatter, err = GetFormatter("%d")
	if err != nil {
		t.Fatal()
	}
	fmv = formatter.Format(mv)
	if !fmv.IsBool() {
		t.Fatal()
	}
	if fmv.String() != "true" {
		t.Fatal()
	}

	mv = FALSE
	formatter, err = GetFormatter("%d")
	if err != nil {
		t.Fatal()
	}
	fmv = formatter.Format(mv)
	if !fmv.IsBool() {
		t.Fatal()
	}
	if fmv.String() != "false" {
		t.Fatal()
	}

	mv = FromFloat(10.123)
	formatter, err = GetFormatter("%d")
	if err != nil {
		t.Fatal()
	}
	fmv = formatter.Format(mv)
	if !fmv.IsInt() {
		t.Fatal()
	}
	if fmv.String() != "10" {
		t.Fatal()
	}

	mv = FromFloat(10.1)
	formatter, err = GetFormatter("%d")
	if err != nil {
		t.Fatal()
	}
	fmv = formatter.Format(mv)
	if !fmv.IsInt() {
		t.Fatal()
	}
	if fmv.String() != "10" {
		t.Fatal()
	}
}
