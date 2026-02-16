package input

import (
	"io"
	"strings"

	"pault.ag/go/debian/control"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// DCF list-valued field names (comma-separated in the format). These are
// exposed as Miller arrays; all other fields remain strings.
var dcfListFieldNames = map[string]bool{
	"Depends":               true,
	"Pre-Depends":           true,
	"Recommends":            true,
	"Suggests":              true,
	"Enhances":              true,
	"Breaks":                true,
	"Conflicts":             true,
	"Replaces":              true,
	"Build-Depends":         true,
	"Build-Depends-Indep":   true,
	"Build-Conflicts":       true,
	"Build-Conflicts-Indep": true,
	"Built-Using":           true,
}

type RecordReaderDCF struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int64
}

func NewRecordReaderDCF(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (*RecordReaderDCF, error) {
	return &RecordReaderDCF{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
	}, nil
}

func (reader *RecordReaderDCF) Read(
	filenames []string,
	context types.Context,
	readerChannel chan<- []*types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool,
) {
	if filenames != nil {
		if len(filenames) == 0 {
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
					handle.Close()
				}
			}
		}
	}
	readerChannel <- types.NewEndOfStreamMarkerList(&context)
}

func (reader *RecordReaderDCF) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- []*types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool,
) {
	context.UpdateForStartOfFile(filename)
	recordsPerBatch := reader.recordsPerBatch

	pr, err := control.NewParagraphReader(handle, nil)
	if err != nil {
		errorChannel <- err
		return
	}

	recordsAndContexts := make([]*types.RecordAndContext, 0, recordsPerBatch)
	i := int64(0)

	for {
		i++
		if i%recordsPerBatch == 0 {
			select {
			case <-downstreamDoneChannel:
				goto flush
			default:
			}
		}

		para, err := pr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			errorChannel <- err
			return
		}
		if para == nil {
			break
		}

		rec := dcfParagraphToRecord(para, reader.readerOptions.DedupeFieldNames)
		if rec == nil {
			// dedupe or other error already sent
			return
		}

		context.UpdateForInputRecord()
		recordsAndContexts = append(recordsAndContexts, types.NewRecordAndContext(rec, context))
		if int64(len(recordsAndContexts)) >= recordsPerBatch {
			readerChannel <- recordsAndContexts
			recordsAndContexts = make([]*types.RecordAndContext, 0, recordsPerBatch)
		}
	}

flush:
	if len(recordsAndContexts) > 0 {
		readerChannel <- recordsAndContexts
	}
}

// splitCommaList splits a DCF comma-separated value and returns Miller Mlrval
// array elements (trimmed, skipping empty).
func splitCommaList(s string) []*mlrval.Mlrval {
	parts := strings.Split(s, ",")
	out := make([]*mlrval.Mlrval, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, mlrval.FromString(p))
	}
	return out
}

// dcfParagraphToRecord converts a control paragraph to a Miller record. List
// fields (Depends, Build-Depends, etc.) become arrays; all other fields stay
// strings. Uses para.Order to preserve field order. Returns nil on error (caller
// must have sent on errorChannel).
func dcfParagraphToRecord(para *control.Paragraph, dedupeFieldNames bool) *mlrval.Mlrmap {
	rec := mlrval.NewMlrmapAsRecord()

	// Preserve order from the DCF file when present.
	keys := para.Order
	if len(keys) == 0 {
		for k := range para.Values {
			keys = append(keys, k)
		}
	}

	for _, k := range keys {
		v, ok := para.Values[k]
		if !ok {
			continue
		}

		var mv *mlrval.Mlrval
		if dcfListFieldNames[k] {
			mv = mlrval.FromArray(splitCommaList(v))
		} else {
			mv = mlrval.FromString(v)
		}

		if _, err := rec.PutReferenceMaybeDedupe(k, mv, dedupeFieldNames); err != nil {
			return nil
		}
	}

	return rec
}
