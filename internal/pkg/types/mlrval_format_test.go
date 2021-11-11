// ================================================================
// Most Miller tests (thousands of them) are command-line-driven via
// mlr regtest. Here are some cases needing special focus.
// ================================================================

// Invoke as:
// * cd internal/pkg/types
// * go test
// Or:
// * cd go
// * go test mlr/internal/pkg/types/...

package types

import (
	"testing"
)

func TestFormatterToString(t *testing.T) {

	mv := MlrvalFromString("hello")
	formatter := newMlrvalFormatterToString("%s")
	fmv := formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "hello" {
		t.Fatal()
	}

	mv = MLRVAL_NULL
	formatter = newMlrvalFormatterToString("%s")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "null" {
		t.Fatal()
	}

	mv = MLRVAL_TRUE
	formatter = newMlrvalFormatterToString("%s")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "true" {
		t.Fatal()
	}

	mv = MLRVAL_FALSE
	formatter = newMlrvalFormatterToString("%s")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "false" {
		t.Fatal()
	}

	mv = MlrvalFromInt(10)
	formatter = newMlrvalFormatterToString("%s")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "10" {
		t.Fatal()
	}

	mv = MlrvalFromString("hello")
	formatter = newMlrvalFormatterToString("[[%s]]")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "[[hello]]" {
		t.Fatal()
	}

}

func TestFormatterToInt(t *testing.T) {

	mv := MlrvalFromString("hello")
	formatter := newMlrvalFormatterToInt("%d")
	fmv := formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "hello" {
		t.Fatal()
	}

	mv = MLRVAL_NULL
	formatter = newMlrvalFormatterToInt("%d")
	fmv = formatter.Format(mv)
	if !fmv.IsNull() {
		t.Fatal()
	}
	if fmv.String() != "null" {
		t.Fatal()
	}

	mv = MLRVAL_TRUE
	formatter = newMlrvalFormatterToInt("%d")
	fmv = formatter.Format(mv)
	if !fmv.IsBool() {
		t.Fatal()
	}
	if fmv.String() != "true" {
		t.Fatal()
	}

	mv = MLRVAL_FALSE
	formatter = newMlrvalFormatterToInt("%d")
	fmv = formatter.Format(mv)
	if !fmv.IsBool() {
		t.Fatal()
	}
	if fmv.String() != "false" {
		t.Fatal()
	}

	mv = MlrvalFromInt(10)
	formatter = newMlrvalFormatterToInt("%d")
	fmv = formatter.Format(mv)
	if !fmv.IsInt() {
		t.Fatal()
	}
	if fmv.String() != "10" {
		t.Fatal()
	}

	mv = MlrvalFromInt(10)
	formatter = newMlrvalFormatterToInt("[[0x%x]]")
	fmv = formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "[[0xa]]" {
		t.Fatal()
	}

	mv = MlrvalFromFloat64(10.1)
	formatter = newMlrvalFormatterToInt("%d")
	fmv = formatter.Format(mv)
	if !fmv.IsInt() {
		t.Fatal()
	}
	if fmv.String() != "10" {
		t.Fatal()
	}

}

func TestFormatterToFloat(t *testing.T) {

	mv := MlrvalFromString("hello")
	formatter := newMlrvalFormatterToFloat("%.4f")
	fmv := formatter.Format(mv)
	if !fmv.IsString() {
		t.Fatal()
	}
	if fmv.String() != "hello" {
		t.Fatal()
	}

	mv = MLRVAL_NULL
	formatter = newMlrvalFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	if !fmv.IsNull() {
		t.Fatal()
	}
	if fmv.String() != "null" {
		t.Fatal()
	}

	mv = MLRVAL_TRUE
	formatter = newMlrvalFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	if !fmv.IsBool() {
		t.Fatal()
	}
	if fmv.String() != "true" {
		t.Fatal()
	}

	mv = MLRVAL_FALSE
	formatter = newMlrvalFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	if !fmv.IsBool() {
		t.Fatal()
	}
	if fmv.String() != "false" {
		t.Fatal()
	}

	mv = MlrvalFromFloat64(10)
	formatter = newMlrvalFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	if !fmv.IsFloat() {
		t.Fatal()
	}
	if fmv.String() != "10.0000" {
		t.Fatal()
	}

	mv = MlrvalFromFloat64(10.1)
	formatter = newMlrvalFormatterToFloat("%.4f")
	fmv = formatter.Format(mv)
	if !fmv.IsFloat() {
		t.Fatal()
	}
	if fmv.String() != "10.1000" {
		t.Fatal()
	}
}

func TestFormatter(t *testing.T) {

	mv := MlrvalFromString("hello")
	formatter, err := GetMlrvalFormatter("%d")
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

	mv = MLRVAL_NULL
	formatter, err = GetMlrvalFormatter("%d")
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

	mv = MLRVAL_TRUE
	formatter, err = GetMlrvalFormatter("%d")
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

	mv = MLRVAL_FALSE
	formatter, err = GetMlrvalFormatter("%d")
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

	mv = MlrvalFromFloat64(10.123)
	formatter, err = GetMlrvalFormatter("%d")
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

	mv = MlrvalFromFloat64(10.1)
	formatter, err = GetMlrvalFormatter("%d")
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
