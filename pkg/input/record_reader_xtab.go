package input

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/types"
)

type iXTABPairSplitter interface {
	Split(input string) (key, value string, err error)
}

type RecordReaderXTAB struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int64 // distinct from readerOptions.RecordsPerBatch for join/repl
	pairSplitter    iXTABPairSplitter

	// Note: XTAB uses two consecutive IFS in place of an IRS; IRS is ignored
}

// tStanza is for the channelized reader which operates (for performance) in
// its own goroutine. An XTAB "stanza" is a collection of lines which will be
// parsed as a Miller record. Also for performance (to reduce
// goroutine-scheduler thrash) stanzas are delivered in batches (nominally max
// 500 or so). This struct helps us keep each stanza's comment lines along with
// the stanza they originated in.
type tStanza struct {
	dataLines    *list.List
	commentLines *list.List
}

func newStanza() *tStanza {
	return &tStanza{
		dataLines:    list.New(),
		commentLines: list.New(),
	}
}

func NewRecordReaderXTAB(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (*RecordReaderXTAB, error) {
	return &RecordReaderXTAB{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
		pairSplitter:    newXTABPairSplitter(readerOptions),
	}, nil
}

func (reader *RecordReaderXTAB) Read(
	filenames []string,
	context types.Context,
	readerChannel chan<- *list.List, // list of *types.RecordAndContext
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	if filenames != nil { // nil for mlr -n
		if len(filenames) == 0 { // read from stdin
			handle, err := lib.OpenStdin(
				reader.readerOptions.Prepipe,
				reader.readerOptions.PrepipeIsRaw,
				reader.readerOptions.FileInputEncoding,
			)
			if err != nil {
				errorChannel <- err
			}
			reader.processHandle(handle, "(stdin)", &context, readerChannel, errorChannel, downstreamDoneChannel)
		} else {
			for _, filename := range filenames {
				handle, err := lib.OpenFileForRead(
					filename,
					reader.readerOptions.Prepipe,
					reader.readerOptions.PrepipeIsRaw,
					reader.readerOptions.FileInputEncoding,
				)
				if err != nil {
					errorChannel <- err
				} else {
					reader.processHandle(handle, filename, &context, readerChannel, errorChannel, downstreamDoneChannel)
					handle.Close()
				}
			}
		}
	}
	readerChannel <- types.NewEndOfStreamMarkerList(&context)
}

func (reader *RecordReaderXTAB) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- *list.List, // list of *types.RecordAndContext
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)
	recordsPerBatch := reader.recordsPerBatch

	// XTAB uses repeated IFS, rather than IRS, to delimit records
	lineScanner := NewLineScanner(handle, reader.readerOptions.IFS)

	stanzasChannel := make(chan *list.List, recordsPerBatch)
	go channelizedStanzaScanner(lineScanner, reader.readerOptions, stanzasChannel, downstreamDoneChannel,
		recordsPerBatch)

	for {
		recordsAndContexts, eof := reader.getRecordBatch(stanzasChannel, context, errorChannel)
		if recordsAndContexts.Len() > 0 {
			readerChannel <- recordsAndContexts
		}
		if eof {
			break
		}
	}
}

// Given input like
//
//	a 1
//	b 2
//	c 3
//
//	a 4
//	b 5
//	c 6
//
// this function reads the input stream a line at a time, then produces
// string-lists one per stanza where a stanza is delimited by blank line, or
// start or end of file. A single stanza, once parsed, will become a single
// record.
func channelizedStanzaScanner(
	lineScanner *bufio.Scanner,
	readerOptions *cli.TReaderOptions,
	stanzasChannel chan<- *list.List, // list of list of string
	downstreamDoneChannel <-chan bool, // for mlr head
	recordsPerBatch int64,
) {
	numStanzasSeen := int64(0)
	inStanza := false
	done := false

	stanzas := list.New()
	stanza := newStanza()

	for lineScanner.Scan() {
		line := lineScanner.Text()

		// Check for comments-in-data feature
		// TODO: function-pointer this away
		if readerOptions.CommentHandling != cli.CommentsAreData {
			if strings.HasPrefix(line, readerOptions.CommentString) {
				if readerOptions.CommentHandling == cli.PassComments {
					stanza.commentLines.PushBack(line)
					continue
				} else if readerOptions.CommentHandling == cli.SkipComments {
					continue
				}
				// else comments are data
			}
		}

		if line == "" {
			// Empty-line handling:
			// 1. First empty line(s) in the stream are ignored.
			// 2. After that, one or more empty lines separate records.
			// 3. At end of file, multiple empty lines are ignored.
			if inStanza {
				inStanza = false
				stanzas.PushBack(stanza)
				numStanzasSeen++
				stanza = newStanza()
			} else {
				continue
			}
		} else {
			if !inStanza {
				inStanza = true
			}
			stanza.dataLines.PushBack(line)
		}

		// See if downstream processors will be ignoring further data (e.g. mlr
		// head).  If so, stop reading. This makes 'mlr head hugefile' exit
		// quickly, as it should.
		if numStanzasSeen%recordsPerBatch == 0 {
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
			stanzasChannel <- stanzas
			stanzas = list.New()
		}

		if done {
			break
		}
	}

	// The last stanza may not have a trailing newline after it. Any lines in the stanza
	// at this point will form the final record in the stream.
	if stanza.dataLines.Len() > 0 || stanza.commentLines.Len() > 0 {
		stanzas.PushBack(stanza)
	}

	stanzasChannel <- stanzas
	close(stanzasChannel) // end-of-stream marker
}

// TODO: comment copiously we're trying to handle slow/fast/short/long reads: tail -f, smallfile, bigfile.
func (reader *RecordReaderXTAB) getRecordBatch(
	stanzasChannel <-chan *list.List,
	context *types.Context,
	errorChannel chan error,
) (
	recordsAndContexts *list.List,
	eof bool,
) {
	recordsAndContexts = list.New()

	stanzas, more := <-stanzasChannel
	if !more {
		return recordsAndContexts, true
	}

	for e := stanzas.Front(); e != nil; e = e.Next() {
		stanza := e.Value.(*tStanza)

		if stanza.commentLines.Len() > 0 {
			for f := stanza.commentLines.Front(); f != nil; f = f.Next() {
				line := f.Value.(string)
				recordsAndContexts.PushBack(types.NewOutputString(line+reader.readerOptions.IFS, context))
			}
		}

		if stanza.dataLines.Len() > 0 {
			record, err := reader.recordFromXTABLines(stanza.dataLines)
			if err != nil {
				errorChannel <- err
				return
			}
			context.UpdateForInputRecord()
			recordsAndContexts.PushBack(types.NewRecordAndContext(record, context))
		}
	}

	return recordsAndContexts, false
}

func (reader *RecordReaderXTAB) recordFromXTABLines(
	stanza *list.List,
) (*mlrval.Mlrmap, error) {
	record := mlrval.NewMlrmapAsRecord()
	dedupeFieldNames := reader.readerOptions.DedupeFieldNames

	for e := stanza.Front(); e != nil; e = e.Next() {
		line := e.Value.(string)

		key, value, err := reader.pairSplitter.Split(line)
		if err != nil {
			return nil, err
		}

		_, err = record.PutReferenceMaybeDedupe(key, mlrval.FromDeferredType(value), dedupeFieldNames)
		if err != nil {
			return nil, err
		}
	}

	return record, nil
}

// IPairSplitter splits a string into left and right, e.g. for IPS.
// This is similar to the general one for multiple formats; the exception
// is that for XTAB we always allow repeat IPS.
func newXTABPairSplitter(options *cli.TReaderOptions) iXTABPairSplitter {
	if options.IPSRegex == nil {
		return &tXTABIPSSplitter{ips: options.IPS, ipslen: len(options.IPS)}
	} else {
		return &tXTABIPSRegexSplitter{ipsRegex: options.IPSRegex}
	}
}

type tXTABIPSSplitter struct {
	ips    string
	ipslen int
}

// This is a splitter for XTAB lines, like 'abc      123'.  It's not quite the same as the
// field/pair-splitter functions shared by DKVP, NIDX, and CSV-lite. XTAB is the omly format for
// which we need to produce just a pair of items -- a key and a value -- delimited by one or more
// IPS. For exaemple, with IPS being a space, in 'abc     123' we need to get key 'abc' and value
// '123'; for 'abc    123 456' we need key 'abc' and value '123 456'.  It's super-elegant to simply
// regex-split the line like 'kv = lib.RegexSplitString(reader.readerOptions.IPSRegex, line, 2)' --
// however, that's 3x slower than the current implementation. It turns out regexes are great
// but we should use them only when we must, since they are expensive.
func (s *tXTABIPSSplitter) Split(input string) (key, value string, err error) {
	// Empty string is a length-0 return value.
	n := len(input)
	if n == 0 {
		return "", "", fmt.Errorf("internal coding error in XTAB reader")
	}

	// '   abc 123' splits as key '', value 'abc 123'.
	if strings.HasPrefix(input, s.ips) {
		keyStart := 0
		for keyStart < n && strings.HasPrefix(input[keyStart:], s.ips) {
			keyStart += s.ipslen
		}
		return "", input[keyStart:n], nil
	}

	// Find the first IPS, if any. If there isn't any in the input line then there is no value, only key:
	// e.g. the line is 'abc'.
	var keyEnd, valueStart int
	foundIPS := false
	for keyEnd = 1; keyEnd <= n; keyEnd++ {
		if strings.HasPrefix(input[keyEnd:], s.ips) {
			foundIPS = true
			break
		}
	}
	if !foundIPS {
		return input, "", nil
	}

	// Find the first non-IPS character after last-found IPS, if any. If there isn't any in the input
	// line then there is no value, only key: e.g. the line is 'abc   '.
	foundValue := false
	for valueStart = keyEnd + s.ipslen; valueStart <= n; valueStart++ {
		if !strings.HasPrefix(input[valueStart:], s.ips) {
			foundValue = true
			break
		}
	}
	if !foundValue {
		return input[0:keyEnd], "", nil
	}

	return input[0:keyEnd], input[valueStart:n], nil
}

type tXTABIPSRegexSplitter struct {
	ipsRegex *regexp.Regexp
}

func (s *tXTABIPSRegexSplitter) Split(input string) (key, value string, err error) {
	kv := lib.RegexSplitString(s.ipsRegex, input, 2)
	if len(kv) == 0 {
		return "", "", fmt.Errorf("internal coding error in XTAB reader")
	} else if len(kv) == 1 {
		return kv[0], "", nil
	} else if len(kv) == 2 {
		return kv[0], kv[1], nil
	} else {
		return "", "", fmt.Errorf("internal coding error in XTAB reader")
	}
}
