// Tests for https://github.com/johnkerl/miller/issues/1810 and
// https://github.com/johnkerl/miller/issues/1722: `--ors '\r\n'` (or `--ors
// crlf`) should be honored for CSV/CSV-lite/TSV output. These are unit tests
// rather than regression-test cases since the regression-test harness
// normalizes CR/LF to LF before comparing outputs, and so cannot assert
// byte-exact line endings.

package output

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func makeTestRecord() *mlrval.Mlrmap {
	record := mlrval.NewMlrmap()
	record.PutReference("a", mlrval.FromString("1"))
	record.PutReference("b", mlrval.FromString("2"))
	return record
}

// runWriter drives a record-writer with a single a=1,b=2 record and returns
// the output bytes.
func runWriter(t *testing.T, writer IRecordWriter) string {
	t.Helper()
	var buffer bytes.Buffer
	bufferedOutputStream := bufio.NewWriter(&buffer)
	err := writer.Write(makeTestRecord(), nil, bufferedOutputStream, false)
	if err != nil {
		t.Fatal(err)
	}
	err = writer.Write(nil, nil, bufferedOutputStream, false) // end of stream
	if err != nil {
		t.Fatal(err)
	}
	bufferedOutputStream.Flush()
	return buffer.String()
}

func TestCSVWriterLFDefault(t *testing.T) {
	writer, err := NewRecordWriterCSV(&cli.TWriterOptions{OFS: ",", ORS: "\n"})
	if err != nil {
		t.Fatal(err)
	}
	output := runWriter(t, writer)
	expected := "a,b\n1,2\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func TestCSVWriterCRLF(t *testing.T) {
	writer, err := NewRecordWriterCSV(&cli.TWriterOptions{OFS: ",", ORS: "\r\n"})
	if err != nil {
		t.Fatal(err)
	}
	output := runWriter(t, writer)
	expected := "a,b\r\n1,2\r\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func TestCSVWriterRejectsOtherORS(t *testing.T) {
	_, err := NewRecordWriterCSV(&cli.TWriterOptions{OFS: ",", ORS: ";"})
	if err == nil {
		t.Fatal("expected error for ORS \";\" but got none")
	}
}

func TestCSVLiteWriterCRLF(t *testing.T) {
	writer, err := NewRecordWriterCSVLite(&cli.TWriterOptions{OFS: ",", ORS: "\r\n"})
	if err != nil {
		t.Fatal(err)
	}
	output := runWriter(t, writer)
	expected := "a,b\r\n1,2\r\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func TestTSVWriterCRLF(t *testing.T) {
	writer, err := NewRecordWriterTSV(&cli.TWriterOptions{OFS: "\t", ORS: "\r\n"})
	if err != nil {
		t.Fatal(err)
	}
	output := runWriter(t, writer)
	expected := "a\tb\r\n1\t2\r\n"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}
