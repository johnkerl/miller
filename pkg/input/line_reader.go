// This file contains the interface for file-format-specific record-readers, as
// well as a collection of utility functions.

package input

import (
	"bufio"
	"container/list"
	"io"
	"strings"

	"github.com/johnkerl/miller/pkg/lib"
)

type ILineReader interface {
	// Read returns the string without the final newline (or whatever terminator).
	// The error condition io.EOF as non-error "error" case.
	// EOF is always returned with empty line: the code here is structured so that
	// we do not return a non-empty line along with an EOF indicator.
	Read() (string, error)
}

type DefaultLineReader struct {
	underlying *bufio.Reader
	eof        bool
}

// SingleIRSLineReader handles reading lines with a single-character terminator.
type SingleIRSLineReader struct {
	underlying *bufio.Reader
	end_irs    byte
	eof        bool
}

// MultiIRSLineReader handles reading lines which may be delimited by multi-line separators, e.g.
// "\xe2\x90\x9e" for USV.
type MultiIRSLineReader struct {
	underlying *bufio.Reader
	irs        string
	irs_len    int
	end_irs    byte
	eof        bool
}

func NewLineReader(handle io.Reader, irs string) ILineReader {
	underlying := bufio.NewReader(handle)

	irs_len := len(irs)

	// Not worth complicating the API by adding an error return.
	// Empty IRS is checked elsewhere.
	if irs_len < 1 {
		panic("Empty IRS")

	} else if irs == "\n" || irs == "\r\n" {
		return &DefaultLineReader{
			underlying: underlying,
		}

	} else if irs_len == 1 {
		return &SingleIRSLineReader{
			underlying: underlying,
			end_irs:    irs[0],
		}

	} else {
		return &MultiIRSLineReader{
			underlying: underlying,
			irs:        irs,
			irs_len:    irs_len,
			end_irs:    irs[irs_len-1],
		}
	}
}

func (r *DefaultLineReader) Read() (string, error) {

	if r.eof {
		return "", io.EOF
	}

	line, err := r.underlying.ReadString('\n')

	// If we have EOF and a non-empty line, defer the EOF return to the next Read call.
	if len(line) > 0 && lib.IsEOF(err) {
		r.eof = true
		err = nil
	}

	n := len(line)
	if strings.HasSuffix(line, "\r\n") {
		line = line[:n-2]
	} else if strings.HasSuffix(line, "\n") {
		line = line[:n-1]
	}

	return line, err
}

func (r *SingleIRSLineReader) Read() (string, error) {

	if r.eof {
		return "", io.EOF
	}

	line, err := r.underlying.ReadString(r.end_irs)

	// If we have EOF and a non-empty line, defer the EOF return to the next Read call.
	if len(line) > 0 && lib.IsEOF(err) {
		r.eof = true
		err = nil
	}

	n := len(line)
	if n > 0 && line[n-1] == r.end_irs {
		line = line[:n-1]
	}

	return line, err
}

func (r *MultiIRSLineReader) Read() (string, error) {

	// bufio.Reader.ReadString supports only a single-character terminator.  So we read lines ending
	// in the final character, until we get a line that ends in the entire sequence or EOF.
	//
	// Note that bufio.Scanner has a very nice bufio.Scanner.Split method which can be overridden to
	// support custom line-ending logic.  Sadly, though, bufio.Scanner _only_ supports a fixed
	// maximum line length, and misbehaves badly when presented with longer lines.  So we cannot use
	// bufio.Scanner.  See also https://github.com/johnkerl/miller/issues/1501.

	if r.eof {
		return "", io.EOF
	}

	line := ""

	for {

		piece, err := r.underlying.ReadString(r.end_irs)

		// If we have EOF and a non-empty line, defer the EOF return to the next Read call.
		if len(piece) > 0 && lib.IsEOF(err) {
			r.eof = true
			err = nil
		}

		if err != nil {
			return line, err // includes io.EOF as a non-error "error" case
		}

		if strings.HasSuffix(piece, r.irs) {
			piece = piece[:len(piece)-r.irs_len]
			line += piece
			break
		}

		if r.eof {
			line += piece
			break
		}

	}

	return line, nil
}

// channelizedLineReader puts the line reading/splitting into its own goroutine in order to pipeline
// the I/O with regard to further processing. Used by record-readers for multiple file formats.
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
		if err != nil {
			if lib.IsEOF(err) {
				done = true
				break
			} else {
				break
			}
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
