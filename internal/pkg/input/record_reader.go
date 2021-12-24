// This file contains the interface for file-format-specific record-readers, as
// well as a collection of utility functions.

package input

import (
	"bufio"
	"container/list"
	"io"
	"regexp"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/types"
)

const CSV_BOM = "\xef\xbb\xbf"

// Since Go is concurrent, the context struct (AWK-like variables such as
// FILENAME, NF, NF, FNR, etc.) needs to be duplicated and passed through the
// channels along with each record. Hence the initial context, which readers
// update on each new file/record, and the channel of types.RecordAndContext
// rather than channel of mlrval.Mlrmap.

type IRecordReader interface {
	Read(
		filenames []string,
		initialContext types.Context,
		readerChannel chan<- *list.List, // list of *types.RecordAndContext
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

// TODO: comment copiously
//
// Lines are written to the channel with their trailing newline (or whatever
// IRS) stripped off. So, callers get "a=1,b=2" rather than "a=1,b=2\n".
func channelizedLineScanner(
	lineScanner *bufio.Scanner,
	linesChannel chan<- *list.List,
	downstreamDoneChannel <-chan bool, // for mlr head
	recordsPerBatch int,
) {
	i := 0
	done := false

	lines := list.New()

	for lineScanner.Scan() {
		i++

		lines.PushBack(lineScanner.Text())

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

// IPairSplitter splits a string into left and right, e.g. for IPS.
// This helps us reuse code for splitting by IPS string, or IPS regex.
type iPairSplitter interface {
	Split(input string) []string
}

// IFieldSplitter splits a string into pieces, e.g. for IFS.
// This helps us reuse code for splitting by IFS string, or IFS regex.
type iFieldSplitter interface {
	Split(input string) []string
}

func newPairSplitter(options *cli.TReaderOptions) iPairSplitter {
	if options.IPSRegex == nil {
		return &tIPSSplitter{ips: options.IPS}
	} else {
		return &tIPSRegexSplitter{ipsRegex: options.IPSRegex}
	}
}

func newFieldSplitter(options *cli.TReaderOptions) iFieldSplitter {
	if options.IFSRegex == nil {
		return &tIFSSplitter{ifs: options.IFS}
	} else {
		return &tIFSRegexSplitter{ifsRegex: options.IFSRegex}
	}
}

type tIPSSplitter struct {
	ips string
}

func (s *tIPSSplitter) Split(input string) []string {
	return strings.SplitN(input, s.ips, 2)
}

type tIPSRegexSplitter struct {
	ipsRegex *regexp.Regexp
}

func (s *tIPSRegexSplitter) Split(input string) []string {
	return lib.RegexSplitString(s.ipsRegex, input, 2)
}

type tIFSSplitter struct {
	ifs string
}

func (s *tIFSSplitter) Split(input string) []string {
	return lib.SplitString(input, s.ifs)
}

type tIFSRegexSplitter struct {
	ifsRegex *regexp.Regexp
}

func (s *tIFSRegexSplitter) Split(input string) []string {
	return lib.RegexSplitString(s.ifsRegex, input, -1)
}
