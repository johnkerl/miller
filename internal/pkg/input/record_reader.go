package input

import (
	"bufio"
	"io"

	"mlr/internal/pkg/types"
)

const CSV_BOM = "\xef\xbb\xbf"

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NF, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record. Hence the initial context, which readers
// update on each new file/record, and the channel of types.RecordAndContext
// rather than channel of types.Mlrmap.

type IRecordReader interface {
	Read(
		filenames []string,
		initialContext types.Context,
		inputChannel chan<- *types.RecordAndContext,
		errorChannel chan error,
		downstreamDoneChannel <-chan bool, // for mlr head
	)
}

// NewLineScanner handles read lines which may be delimited by multi-line separators,
// e.g. "\xe2\x90\x9e" for USV.
func NewLineScanner(handle io.Reader, irs string) *bufio.Scanner {
	scanner := bufio.NewScanner(handle)

	// Handled by default scanner.
	if irs == "\n" || irs == "\r\n" {
		return scanner
	}

	irsbytes := []byte(irs)
	irslen := len(irsbytes)

	// Custom splitter
	recordSplitter := func(
		data []byte,
		atEOF bool,
	) (
		advance int,
		token []byte,
		err error,
	) {
		datalen := len(data)
		end := datalen - irslen
		for i := 0; i <= end; i++ {
			if data[i] == irsbytes[0] {
				match := true
				for j := 1; j < irslen; j++ {
					if data[i+j] != irsbytes[j] {
						match = false
						break
					}
				}
				if match {
					return i + irslen, data[:i], nil
				}
			}
		}
		if !atEOF {
			return 0, nil, nil
		}
		// There is one final token to be delivered, which may be the empty string.
		// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
		// but does not trigger an error to be returned from Scan itself.
		return 0, data, bufio.ErrFinalToken
	}

	scanner.Split(recordSplitter)

	return scanner
}
