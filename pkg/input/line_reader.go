// This file contains the interface for file-format-specific record-readers, as
// well as a collection of utility functions.

package input

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"strings"
)

type ILineReader interface {
	Read() (string, error) // includes io.EOF as non-error "error" case
}

type TLineReader struct {
	underlying   *bufio.Reader
	full_irs     string
	full_irs_len int
	end_irs      byte
}

// NewLineReader handles reading lines which may be delimited by multi-line separators,
// e.g. "\xe2\x90\x9e" for USV.
func NewLineReader(handle io.Reader, irs string) *TLineReader {
	underlying := bufio.NewReader(handle)

	// XXX TEMP
	return &TLineReader{
		underlying:   underlying,
		full_irs:     irs,
		full_irs_len: len(irs),
		end_irs:      irs[0],
	}
}

// XXX TO DO: handle the splitter which bufio.Scanner has but bufio.Reader lacks
//	if irs == "\n" || irs == "\r\n" {
//		// Handled by default scanner.
//	} else {
//		irsbytes := []byte(irs)
//		irslen := len(irsbytes)
//
//		// Custom splitter
//		recordSplitter := func(
//			data []byte,
//			atEOF bool,
//		) (
//			advance int,
//			token []byte,
//			err error,
//		) {
//			datalen := len(data)
//			end := datalen - irslen
//			for i := 0; i <= end; i++ {
//				if data[i] == irsbytes[0] {
//					match := true
//					for j := 1; j < irslen; j++ {
//						if data[i+j] != irsbytes[j] {
//							match = false
//							break
//						}
//					}
//					if match {
//						return i + irslen, data[:i], nil
//					}
//				}
//			}
//			if !atEOF {
//				return 0, nil, nil
//			}
//			// There is one final token to be delivered, which may be the empty string.
//			// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
//			// but does not trigger an error to be returned from Scan itself.
//			return 0, data, bufio.ErrFinalToken
//		}
//
//		underlying.Split(recordSplitter)
//	}

func (r *TLineReader) Read() (string, error) {
	line, err := r.underlying.ReadString(r.end_irs)
	if err != nil {
		return "", err
	}

	if strings.HasSuffix(line, r.full_irs) {
		line = line[:len(line)-r.full_irs_len]
	}
	return line, nil
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

	for {
		line, err := lineReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %#v\n", err)
			break
		}

		i++

		lines.PushBack(line)

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
