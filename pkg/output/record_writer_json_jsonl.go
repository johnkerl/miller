package output

import (
	"bufio"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type RecordWriterJSON struct {
	// Parameters:
	writerOptions  *cli.TWriterOptions
	jsonFormatting mlrval.TJSONFormatting
	jvQuoteAll     bool

	// State:
	wroteAnyRecords bool
}

func NewRecordWriterJSON(writerOptions *cli.TWriterOptions) (*RecordWriterJSON, error) {
	var jsonFormatting mlrval.TJSONFormatting = mlrval.JSON_SINGLE_LINE
	if writerOptions.JSONOutputMultiline {
		jsonFormatting = mlrval.JSON_MULTILINE
	}
	return &RecordWriterJSON{
		writerOptions:   writerOptions,
		jsonFormatting:  jsonFormatting,
		jvQuoteAll:      writerOptions.JVQuoteAll,
		wroteAnyRecords: false,
	}, nil
}

func NewRecordWriterJSONLines(writerOptions *cli.TWriterOptions) (*RecordWriterJSON, error) {
	wopt := *writerOptions
	wopt.WrapJSONOutputInOuterList = false
	wopt.JSONOutputMultiline = false
	return &RecordWriterJSON{
		writerOptions:   &wopt,
		jsonFormatting:  mlrval.JSON_SINGLE_LINE,
		jvQuoteAll:      writerOptions.JVQuoteAll,
		wroteAnyRecords: false,
	}, nil
}

func (writer *RecordWriterJSON) Write(
	outrec *mlrval.Mlrmap,
	context *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if outrec != nil && writer.jvQuoteAll {
		outrec.StringifyValuesRecursively()
	}

	if writer.writerOptions.WrapJSONOutputInOuterList {
		return writer.writeWithListWrap(outrec, context, bufferedOutputStream, outputIsStdout)
	}
	return writer.writeWithoutListWrap(outrec, context, bufferedOutputStream, outputIsStdout)
}

func (writer *RecordWriterJSON) writeWithListWrap(
	outrec *mlrval.Mlrmap,
	context *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if outrec != nil { // Not end of record stream
		if !writer.wroteAnyRecords {
			bufferedOutputStream.WriteString("[\n")
		}

		// The Mlrmap FormatAsJSON doesn't include the final newline, so that we
		// can place it neatly with commas here (if the user requested them).
		s, err := outrec.FormatAsJSON(writer.jsonFormatting, outputIsStdout)
		if err != nil {
			return err
		}

		if writer.wroteAnyRecords {
			bufferedOutputStream.WriteString(",\n")
		}

		bufferedOutputStream.WriteString(s)

		writer.wroteAnyRecords = true

	} else { // End of record stream

		if !writer.wroteAnyRecords {
			if context.JSONHadBrackets {
				bufferedOutputStream.WriteString("[")
				bufferedOutputStream.WriteString("\n]\n")
			}
		} else {
			bufferedOutputStream.WriteString("\n]\n")
		}

	}
	return nil
}

func (writer *RecordWriterJSON) writeWithoutListWrap(
	outrec *mlrval.Mlrmap,
	_ *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if outrec == nil {
		// End of record stream
		return nil
	}

	// The Mlrmap FormatAsJSON doesn't include the final newline, so that we
	// can place it neatly with commas here (if the user requested them).
	s, err := outrec.FormatAsJSON(writer.jsonFormatting, outputIsStdout)
	if err != nil {
		return err
	}

	bufferedOutputStream.WriteString(s)
	bufferedOutputStream.WriteString("\n")
	return nil
}
