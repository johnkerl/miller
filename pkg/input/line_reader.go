// This file contains the interface for file-format-specific record-readers, as
// well as a collection of utility functions.

package input

import (
	"bufio"
	"container/list"
	"io"
)

type ILineReader interface {
	Scan() bool
	Text() string
}

type TLineReader struct {
	scanner *bufio.Scanner
}

// NewLineReader handles reading lines which may be delimited by multi-line separators,
// e.g. "\xe2\x90\x9e" for USV.
func NewLineReader(handle io.Reader, irs string) *TLineReader {
	scanner := bufio.NewScanner(handle)

	if irs == "\n" || irs == "\r\n" {
		// Handled by default scanner.
	} else {
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
	}

	return &TLineReader{
		scanner: scanner,
	}
}

func (r *TLineReader) Scan() bool {
	return r.scanner.Scan()
}

func (r *TLineReader) Text() string {
	return r.scanner.Text()
}

// TODO: comment copiously
//
// Lines are written to the channel with their trailing newline (or whatever
// IRS) stripped off. So, callers get "a=1,b=2" rather than "a=1,b=2\n".
func channelizedLineReader(
	lineReader ILineReader,
	linesChannel chan<- *list.List,
	downstreamDoneChannel <-chan bool, // for mlr head
	recordsPerBatch int64,
) {
	i := int64(0)
	done := false

	lines := list.New()

	for lineReader.Scan() {
		i++

		lines.PushBack(lineReader.Text())

		// See if downstream processors will be ignoring further data (e.g. mlr
		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
		// quickly, as it should.
		if i%recordsPerBatch == 0 {
			select {
			case _ = <-downstreamDoneChannel:
				done = true
				break
			default:
				break
			}
			if done {
				break
			}
			linesChannel <- lines
			lines = list.New()
		}

		if done {
			break
		}
	}
	linesChannel <- lines
	close(linesChannel) // end-of-stream marker
}
