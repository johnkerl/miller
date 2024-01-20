package output

import (
	"bufio"
	"fmt"
	"strings"

	csv "github.com/johnkerl/miller/pkg/go-csv"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/mlrval"
)

type RecordWriterCSV struct {
	writerOptions     *cli.TWriterOptions
	ofs0              byte // Go's CSV library only lets its 'Comma' be a single character
	csvWriter         *csv.Writer
	needToPrintHeader bool
	firstRecordKeys   []string
	firstRecordNF     int64
	quoteAll          bool // For double-quote around all fields
}

func NewRecordWriterCSV(writerOptions *cli.TWriterOptions) (*RecordWriterCSV, error) {
	if len(writerOptions.OFS) != 1 {
		return nil, fmt.Errorf("for CSV, OFS can only be a single character")
	}
	if writerOptions.ORS != "\n" && writerOptions.ORS != "\r\n" {
		return nil, fmt.Errorf("for CSV, ORS cannot be altered")
	}
	writer := &RecordWriterCSV{
		writerOptions:     writerOptions,
		csvWriter:         nil, // will be set on first Write() wherein we have the output stream
		needToPrintHeader: !writerOptions.HeaderlessOutput,
		firstRecordKeys:   nil,
		firstRecordNF:     -1,
		quoteAll:          writerOptions.CSVQuoteAll,
	}
	return writer, nil
}

func (writer *RecordWriterCSV) Write(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return nil
	}

	if writer.csvWriter == nil {
		writer.csvWriter = csv.NewWriter(bufferedOutputStream)
		writer.csvWriter.Comma = rune(writer.writerOptions.OFS[0]) // xxx temp
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
		writer.WriteCSVRecordMaybeColorized(fields, bufferedOutputStream, outputIsStdout, true, writer.quoteAll)
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
				"CSV schema change: first keys \"%s\"; current keys \"%s\"",
				strings.Join(writer.firstRecordKeys, writer.writerOptions.OFS),
				strings.Join(outrec.GetKeys(), writer.writerOptions.OFS),
			)
		}
		fields[i] = pe.Value.String()
		i++
	}

	for ; i < outputNF; i++ {
		fields[i] = ""
	}

	writer.WriteCSVRecordMaybeColorized(fields, bufferedOutputStream, outputIsStdout, false, writer.quoteAll)

	return nil
}
