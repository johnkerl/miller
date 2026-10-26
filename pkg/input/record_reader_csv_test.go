package input

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"testing"
)

func TestBOMStrippingReaderGzip(t *testing.T) {
	testCases := []struct {
		name            string
		input           string
		want            string
		corruptChecksum bool
	}{
		{"BOM and quoted header", "\xef\xbb\xbf\"a\"\n1\n", "\"a\"\n1\n", false},
		{"no BOM", "\"a\"\n1\n", "\"a\"\n1\n", false},
		{"BOM only", "\xef\xbb\xbf", "", false},
		{"empty input", "", "", false},
		{"checksum error", "\xef\xbb\xbf\"a\"\n1\n", "\"a\"\n1\n", true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var compressed bytes.Buffer
			writer := gzip.NewWriter(&compressed)
			if _, err := io.WriteString(writer, testCase.input); err != nil {
				t.Fatal(err)
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			var wantErr error
			if testCase.corruptChecksum {
				// Corrupt the CRC32 in the gzip trailer.
				compressed.Bytes()[compressed.Len()-8] ^= 1
				wantErr = gzip.ErrChecksum
			}

			reader, err := gzip.NewReader(&compressed)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := reader.Close(); err != nil {
					t.Error(err)
				}
			}()

			got, err := io.ReadAll(NewBOMStrippingReader(reader))
			if !errors.Is(err, wantErr) {
				t.Fatalf("got error %v, want %v", err, wantErr)
			}
			if string(got) != testCase.want {
				t.Fatalf("got %q, want %q", got, testCase.want)
			}
		})
	}
}
