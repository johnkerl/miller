package input

import (
	"fmt"
	"io"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// RecordReaderREC reads GNU recutils (.rec) format: records are stanzas of
// "FieldName: Value" lines separated by one or more blank lines, with two
// continuation-line mechanisms (trailing-backslash logical-line joining, and
// "+"-prefixed embedded-newline continuation) and "#" comments. Record
// descriptors ("%rec: ..." stanzas) are not given any special
// schema/constraint interpretation -- they are parsed like any other record,
// since Miller has no schema-enforcement concept.
type RecordReaderREC struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int64 // distinct from readerOptions.RecordsPerBatch for join/repl
	// recordArena batch-allocates field entries/values; reset per getRecordBatch.
	recordArena *mlrval.RecordArena

	// Note: recutils uses blank lines to delimit records, not a
	// configurable IRS; and ": " to delimit fields, not a configurable IPS.
}

func NewRecordReaderREC(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (*RecordReaderREC, error) {
	return &RecordReaderREC{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
		recordArena:     mlrval.NewRecordArena(64),
	}, nil
}

func (reader *RecordReaderREC) Read(
	filenames []string,
	context types.Context,
	readerChannel chan<- []*types.RecordAndContext, // list of *types.RecordAndContext
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
			} else {
				reader.processHandle(handle, "(stdin)", &context, readerChannel, errorChannel, downstreamDoneChannel)
			}
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
					_ = handle.Close()
				}
			}
		}
	}
	readerChannel <- types.NewEndOfStreamMarkerList(&context)
}

func (reader *RecordReaderREC) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- []*types.RecordAndContext, // list of *types.RecordAndContext
	errorChannel chan error,
	downstreamDoneChannel <-chan bool, // for mlr head
) {
	context.UpdateForStartOfFile(filename)
	recordsPerBatch := reader.recordsPerBatch

	// recutils records are blank-line-delimited stanzas, same boundary
	// logic as XTAB; RS is fixed at "\n" (unlike XTAB, it's not
	// configurable via IFS).
	lineReader := NewLineReader(handle, "\n")

	stanzasChannel := make(chan []*tStanza, recordsPerBatch)
	go channelizedStanzaScanner(lineReader, reader.readerOptions, stanzasChannel, downstreamDoneChannel,
		recordsPerBatch)

	for {
		recordsAndContexts, eof := reader.getRecordBatch(stanzasChannel, context, errorChannel)
		if len(recordsAndContexts) > 0 {
			readerChannel <- recordsAndContexts
		}
		if eof {
			break
		}
	}
}

func (reader *RecordReaderREC) getRecordBatch(
	stanzasChannel <-chan []*tStanza,
	context *types.Context,
	errorChannel chan error,
) (
	recordsAndContexts []*types.RecordAndContext,
	eof bool,
) {
	recordsAndContexts = []*types.RecordAndContext{}

	stanzas, more := <-stanzasChannel
	if !more {
		return recordsAndContexts, true
	}

	reader.recordArena = mlrval.NewRecordArena(len(stanzas) * 8)

	for _, stanza := range stanzas {

		if len(stanza.commentLines) > 0 {
			for _, line := range stanza.commentLines {
				recordsAndContexts = append(recordsAndContexts, types.NewOutputString(line+"\n", context))
			}
		}

		if len(stanza.dataLines) > 0 {
			record, err := reader.recordFromRECLines(stanza.dataLines)
			if err != nil {
				errorChannel <- err
				return
			}
			context.UpdateForInputRecord()
			recordsAndContexts = append(recordsAndContexts, types.NewRecordAndContext(record, context))
		}
	}

	return recordsAndContexts, false
}

// recordFromRECLines converts one blank-line-delimited stanza's raw lines
// into a Miller record. This happens in three passes, in this order:
//
//  1. Backslash-newline logical-line joining: a line ending in a literal "\"
//     is joined directly (no separator inserted) with the next physical
//     line.
//  2. "+"-continuation folding: a line starting with "+" continues the
//     immediately preceding field's value with an embedded "\n"; a single
//     leading space after the "+" is stripped.
//  3. Splitting each remaining logical line on the first occurrence of
//     ": " into a field name and value.
//
// Malformed input (a "+"-continuation with no preceding field, or a line
// with no ": " separator) is a hard error -- recutils' own spec requires
// exactly ": " as the field separator, and there is no leniency flag for
// this format (unlike e.g. CSV's ragged-input handling).
func (reader *RecordReaderREC) recordFromRECLines(
	stanza []string,
) (*mlrval.Mlrmap, error) {
	record := reader.recordArena.NewRecord()
	dedupeFieldNames := reader.readerOptions.DedupeFieldNames

	joinedLines := joinRECBackslashContinuations(stanza)

	fieldLines, err := foldRECPlusContinuations(joinedLines)
	if err != nil {
		return nil, err
	}

	for _, line := range fieldLines {
		key, value, found := strings.Cut(line, ": ")
		if !found {
			return nil, fmt.Errorf(
				"mlr: recutils: missing \": \" field separator in line %q",
				line,
			)
		}
		reader.recordArena.PutDeferred(record, key, value, dedupeFieldNames)
	}

	return record, nil
}

// joinRECBackslashContinuations implements recutils' backslash-newline
// logical-line continuation: a line whose last character is "\" is joined
// directly with the next physical line, with the backslash removed and no
// separator inserted.
func joinRECBackslashContinuations(lines []string) []string {
	joined := make([]string, 0, len(lines))
	var pending strings.Builder
	havePending := false

	for _, line := range lines {
		if trimmed, ok := strings.CutSuffix(line, `\`); ok {
			pending.WriteString(trimmed)
			havePending = true
			continue
		}
		if havePending {
			pending.WriteString(line)
			joined = append(joined, pending.String())
			pending.Reset()
			havePending = false
		} else {
			joined = append(joined, line)
		}
	}
	// A stanza-final line ending in "\" has nothing left to join with; keep
	// its (backslash-stripped) content as-is.
	if havePending {
		joined = append(joined, pending.String())
	}

	return joined
}

// foldRECPlusContinuations implements recutils' "+"-continuation: a line
// starting with "+" continues the immediately preceding field's value with
// an embedded "\n". A single leading space after the "+" is stripped (a
// bare "+" folds in an empty continuation line).
func foldRECPlusContinuations(lines []string) ([]string, error) {
	folded := make([]string, 0, len(lines))

	for _, line := range lines {
		if strings.HasPrefix(line, "+") {
			if len(folded) == 0 {
				return nil, fmt.Errorf(
					"mlr: recutils: continuation line %q has no preceding field in this record",
					line,
				)
			}
			continuation := strings.TrimPrefix(line, "+")
			continuation = strings.TrimPrefix(continuation, " ")
			folded[len(folded)-1] = folded[len(folded)-1] + "\n" + continuation
		} else {
			folded = append(folded, line)
		}
	}

	return folded, nil
}
