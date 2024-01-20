package output

import (
	"bufio"
	"fmt"

	csv "github.com/johnkerl/miller/pkg/go-csv"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/mlrval"
)

type RecordWriterCSV struct {
	writerOptions     *cli.TWriterOptions
	ofs0              byte // Go's CSV library only lets its 'Comma' be a single character
	csvWriter         *csv.Writer
	needToPrintHeader bool
	firstFieldCount   int64
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
		firstFieldCount:   -1,
		quoteAll:          writerOptions.CSVQuoteAll,
	}
	return writer, nil
}

func (writer *RecordWriterCSV) Write(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) {
	// End of record stream: nothing special for this output format
	if outrec == nil {
		return
	}

	if writer.csvWriter == nil {
		writer.csvWriter = csv.NewWriter(bufferedOutputStream)
		writer.csvWriter.Comma = rune(writer.writerOptions.OFS[0]) // xxx temp
	}

	if writer.firstFieldCount == -1 {
		writer.firstFieldCount = outrec.FieldCount
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

	outputNF := outrec.FieldCount
	if outputNF < writer.firstFieldCount {
		outputNF = writer.firstFieldCount
	}

	fields := make([]string, outputNF)
	var i int64 = 0
	for pe := outrec.Head; pe != nil; pe = pe.Next {
		fields[i] = pe.Value.String()
		i++
	}

	for ; i < writer.firstFieldCount; i++ {
		fields[i] = ""
	}

	writer.WriteCSVRecordMaybeColorized(fields, bufferedOutputStream, outputIsStdout, false, writer.quoteAll)
}
