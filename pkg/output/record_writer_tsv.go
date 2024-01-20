package output

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/colorizer"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
)

type RecordWriterTSV struct {
	writerOptions     *cli.TWriterOptions
	needToPrintHeader bool
	firstRecordKeys   []string
	firstRecordNF     int64
}

func NewRecordWriterTSV(writerOptions *cli.TWriterOptions) (*RecordWriterTSV, error) {
	if writerOptions.OFS != "\t" {
		return nil, fmt.Errorf("for TSV, OFS cannot be altered")
	}
	if writerOptions.ORS != "\n" && writerOptions.ORS != "\r\n" {
		return nil, fmt.Errorf("for CSV, ORS cannot be altered")
	}
	return &RecordWriterTSV{
		writerOptions:     writerOptions,
		needToPrintHeader: !writerOptions.HeaderlessOutput,
		firstRecordKeys:   nil,
		firstRecordNF:     -1,
	}, nil
}

func (writer *RecordWriterTSV) Write(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return nil
	}

	if writer.firstRecordKeys == nil {
		writer.firstRecordKeys = outrec.GetKeys()
		writer.firstRecordNF = int64(len(writer.firstRecordKeys))
	}

	if writer.needToPrintHeader {
		fields := make([]string, outrec.FieldCount)
		i := 0
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			fields[i] = pe.Key
			i++
		}
		for pe := outrec.Head; pe != nil; pe = pe.Next {
			bufferedOutputStream.WriteString(
				colorizer.MaybeColorizeKey(
					lib.TSVEncodeField(pe.Key),
					outputIsStdout,
				),
			)

			if pe.Next != nil {
				bufferedOutputStream.WriteString(writer.writerOptions.OFS)
			}
		}

		bufferedOutputStream.WriteString(writer.writerOptions.ORS)

		writer.needToPrintHeader = false
	}

	var outputNF int64 = outrec.FieldCount
	if outputNF < writer.firstRecordNF {
		outputNF = writer.firstRecordNF
	}

	fields := make([]string, outputNF)
	var i int64 = 0
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		if i < writer.firstRecordNF && pe.Key != writer.firstRecordKeys[i] {
			return fmt.Errorf(
				"TSV schema change: first keys \"%s\"; current keys \"%s\"",
				strings.Join(writer.firstRecordKeys, writer.writerOptions.OFS),
				strings.Join(outrec.GetKeys(), writer.writerOptions.OFS),
			)
		}
		fields[i] = colorizer.MaybeColorizeValue(
			lib.TSVEncodeField(pe.Value.String()),
			outputIsStdout,
		)
		i++
	}

	for ; i < outputNF; i++ {
		fields[i] = ""
	}

	for j, field := range fields {
		if j > 0 {
			bufferedOutputStream.WriteString(writer.writerOptions.OFS)
		}
		bufferedOutputStream.WriteString(field)
	}

	bufferedOutputStream.WriteString(writer.writerOptions.ORS)

	return nil
}
