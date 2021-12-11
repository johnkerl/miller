package input

import (
	"bufio"
	"container/list"
	"errors"
	"io"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/types"
)

type RecordReaderXTAB struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int
	// Note: XTAB uses two consecutive IFS in place of an IRS; IRS is ignored
}

func NewRecordReaderXTAB(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int,
) (*RecordReaderXTAB, error) {
	return &RecordReaderXTAB{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
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
	recordsPerBatch := reader.readerOptions.RecordsPerBatch

	// XTAB uses repeated IFS, rather than IRS, to delimit records
	lineScanner := NewLineScanner(handle, reader.readerOptions.IFS)

	stanzasChannel := make(chan *list.List, recordsPerBatch)
	go channelizedStanzaScanner(lineScanner, stanzasChannel, downstreamDoneChannel, recordsPerBatch)

	for {
		recordsAndContexts, eof := reader.getRecordBatch(stanzasChannel, context, errorChannel)
		readerChannel <- recordsAndContexts
		if eof {
			break
		}
	}
}

// Given input like
//
//   a 1
//   b 2
//   c 3
//
//   a 4
//   b 5
//   c 6
//
// this function reads the input stream a line at a time, then produces
// string-lists one per stanza where a stanza is delimited by blank line, or
// start or end of file. A single stanza, once parsed, will become a single
// record.
func channelizedStanzaScanner(
	lineScanner *bufio.Scanner,
	stanzasChannel chan<- *list.List, // list of list of string
	downstreamDoneChannel <-chan bool, // for mlr head
	recordsPerBatch int,
) {
	numStanzasSeen := 0
	inStanza := false
	done := false

	stanzas := list.New()
	stanza := list.New()

	for lineScanner.Scan() {
		line := lineScanner.Text()
		if line == "" {
			// Empty-line handling:
			// 1. First empty line(s) in the stream are ignored.
			// 2. After that, one or more empty lines separate records.
			// 3. At end of file, multiple empty lines are ignored.
			if inStanza {
				inStanza = false
				stanzas.PushBack(stanza)
				numStanzasSeen++
				stanza = list.New()
			} else {
				continue
			}
		} else {
			if !inStanza {
				inStanza = true
			}
			stanza.PushBack(line)
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
	if stanza.Len() > 0 {
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
		stanza := e.Value.(*list.List)

		//		// TODO: move
		//		// Check for comments-in-data feature
		//		// TODO: function-pointer this away
		//		if reader.readerOptions.CommentHandling != cli.CommentsAreData {
		//			if strings.HasPrefix(line, reader.readerOptions.CommentString) {
		//				if reader.readerOptions.CommentHandling == cli.PassComments {
		//					recordsAndContexts.PushBack(types.NewOutputString(line+reader.readerOptions.IFS, context))
		//					continue
		//				} else if reader.readerOptions.CommentHandling == cli.SkipComments {
		//					continue
		//				}
		//				// else comments are data
		//			}
		//		}

		lib.InternalCodingErrorIf(stanza.Len() == 0)

		record, err := reader.recordFromXTABLines(stanza)
		if err != nil {
			errorChannel <- err
			return
		}
		context.UpdateForInputRecord()
		recordsAndContexts.PushBack(types.NewRecordAndContext(record, context))
	}

	return recordsAndContexts, false
}

func (reader *RecordReaderXTAB) recordFromXTABLines(
	stanza *list.List,
) (*types.Mlrmap, error) {
	record := types.NewMlrmapAsRecord()

	for e := stanza.Front(); e != nil; e = e.Next() {
		line := e.Value.(string)

		var kv []string
		if reader.readerOptions.IPSRegex == nil { // e.g. --no-ips-regex
			kv = strings.SplitN(line, reader.readerOptions.IPS, 2)
		} else {
			kv = lib.RegexSplitString(reader.readerOptions.IPSRegex, line, 2)
		}
		if len(kv) < 1 {
			return nil, errors.New("mlr: internal coding error in XTAB reader")
		}

		key := kv[0]
		if len(kv) == 1 {
			value := types.MLRVAL_VOID
			record.PutReference(key, value)
		} else {
			value := types.MlrvalFromInferredTypeForDataFiles(kv[1])
			record.PutReference(key, value)
		}
	}

	return record, nil
}
