package input

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type RecordReaderYAML struct {
	readerOptions   *cli.TReaderOptions
	recordsPerBatch int64
}

func NewRecordReaderYAML(
	readerOptions *cli.TReaderOptions,
	recordsPerBatch int64,
) (*RecordReaderYAML, error) {
	return &RecordReaderYAML{
		readerOptions:   readerOptions,
		recordsPerBatch: recordsPerBatch,
	}, nil
}

func (reader *RecordReaderYAML) Read(
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

func (reader *RecordReaderYAML) processHandle(
	handle io.Reader,
	filename string,
	context *types.Context,
	readerChannel chan<- []*types.RecordAndContext,
	errorChannel chan error,
	downstreamDoneChannel <-chan bool,
) {
	context.UpdateForStartOfFile(filename)
	recordsPerBatch := reader.recordsPerBatch
	decoder := yaml.NewDecoder(handle)
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

		mlrval, eof, err := mlrval.MlrvalDecodeFromYAML(decoder)
		if eof {
			break
		}
		if err != nil {
			errorChannel <- err
			return
		}

		if mlrval.IsMap() {
			record := mlrval.GetMap()
			if record == nil {
				errorChannel <- fmt.Errorf("internal coding error in YAML record-reader")
				return
			}
			context.UpdateForInputRecord()
			recordsAndContexts = append(recordsAndContexts, types.NewRecordAndContext(record, context))
			if int64(len(recordsAndContexts)) >= recordsPerBatch {
				readerChannel <- recordsAndContexts
				recordsAndContexts = make([]*types.RecordAndContext, 0, recordsPerBatch)
			}
		} else if mlrval.IsArray() {
			records := mlrval.GetArray()
			if records == nil {
				errorChannel <- fmt.Errorf("internal coding error in YAML record-reader")
				return
			}
			for _, mv := range records {
				if !mv.IsMap() {
					errorChannel <- fmt.Errorf(
						"valid but unmillerable YAML: expected map (object); got %s",
						mv.GetTypeName(),
					)
					return
				}
				record := mv.GetMap()
				if record == nil {
					errorChannel <- fmt.Errorf("internal coding error in YAML record-reader")
					return
				}
				context.UpdateForInputRecord()
				recordsAndContexts = append(recordsAndContexts, types.NewRecordAndContext(record, context))
				if int64(len(recordsAndContexts)) >= recordsPerBatch {
					readerChannel <- recordsAndContexts
					recordsAndContexts = make([]*types.RecordAndContext, 0, recordsPerBatch)
				}
			}
		} else {
			errorChannel <- fmt.Errorf(
				"valid but unmillerable YAML: expected map (object); got %s",
				mlrval.GetTypeName(),
			)
			return
		}
	}

flush:
	if len(recordsAndContexts) > 0 {
		readerChannel <- recordsAndContexts
	}
}
