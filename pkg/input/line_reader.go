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

// TLineReader handles reading lines which may be delimited by multi-line separators, e.g.
// "\xe2\x90\x9e" for USV.
type TLineReader struct {
	underlying *bufio.Reader
	irs        string
	irs_len    int
	end_irs    byte
}

func NewLineReader(handle io.Reader, irs string) *TLineReader {
	underlying := bufio.NewReader(handle)

	// Not worth complicating the API by adding an error return.
	// Empty IRS is checked elsewhere.
	if len(irs) < 1 {
		panic("Empty IRS")
	}

	return &TLineReader{
		underlying: underlying,
		irs:        irs,
		irs_len:    len(irs),
		end_irs:    irs[len(irs)-1],
	}
}

// Read returns the string without the final newline (or whatever terminator).
// The error condition io.EOF as non-error "error" case.
func (r *TLineReader) Read() (string, error) {
	line, err := r.underlying.ReadString(r.end_irs)
	if err != nil {
		return "", err
	}

	if strings.HasSuffix(line, r.irs) {
		line = line[:len(line)-r.irs_len]
	}
	return line, nil
}

// channelizedLineReader puts the line reading/splitting into its own goroutine in order to pipeline
// the I/O with regard to further processing. Used by record-readers for multiple file formats.
//
// Lines are written to the channel with their trailing newline (or whatever
// IRS) stripped off. So, callers get "a=1,b=2" rather than "a=1,b=2\n".
func channelizedLineReader(
	lineReader *TLineReader,
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
