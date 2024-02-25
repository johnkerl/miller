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
	Scan() (bool, error)
	Text() string
}

type TLineReader struct {
	underlying   *bufio.Reader
	full_irs     string
	full_irs_len int
	end_irs      byte
	staged       string
	eof          bool
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
		eof:          false,
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

// XXX ERR RET
func (r *TLineReader) Scan() (bool, error) {
	line, err := r.underlying.ReadString(r.end_irs)
	if err == nil {
		r.staged = line
		// XXX CHOMP
		if strings.HasSuffix(line, r.full_irs) {
			r.staged = line[:len(line)-r.full_irs_len]
		}
		r.eof = false
		return true, nil
	} else if err == io.EOF {
		r.staged = ""
		r.eof = true
		return false, nil
	} else {
		return false, err
	}
}

func (r *TLineReader) Text() string {
	if r.eof {
		panic("PORT ME TOO")
	}
	return r.staged
}

//		if err != nil {
//			fmt.Printf("ERR %v\n", err)
//			break
//		}
//		fmt.Printf("line len %d\n", len(line))
//	}
//	fmt.Println("AFTER LOOP")
//}

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
		eof, err := lineReader.Scan()
		if eof {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %#v\n", err)
			break
		}

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
