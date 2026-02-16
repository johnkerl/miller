package output

import (
	"bufio"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type RecordWriterYAML struct {
	writerOptions   *cli.TWriterOptions
	bufferedRecords []interface{} // used when WrapYAMLOutputInOuterList is true
	wroteAnyRecords bool          // for multi-doc: emit "---\n" before 2nd and later docs
}

func NewRecordWriterYAML(writerOptions *cli.TWriterOptions) (*RecordWriterYAML, error) {
	return &RecordWriterYAML{
		writerOptions:   writerOptions,
		bufferedRecords: nil,
		wroteAnyRecords: false,
	}, nil
}

func (writer *RecordWriterYAML) Write(
	outrec *mlrval.Mlrmap,
	context *types.Context,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
) error {
	if writer.writerOptions.WrapYAMLOutputInOuterList {
		writer.writeWithListWrap(outrec, bufferedOutputStream)
	} else {
		writer.writeWithoutListWrap(outrec, bufferedOutputStream)
	}
	return nil
}

func (writer *RecordWriterYAML) writeWithListWrap(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
) {
	if outrec != nil {
		if writer.bufferedRecords == nil {
			writer.bufferedRecords = []interface{}{}
		}
		native, err := mlrval.MlrmapToYAMLNative(outrec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}
		writer.bufferedRecords = append(writer.bufferedRecords, native)
	} else {
		// End of stream: emit single YAML document as array
		out, err := yaml.Marshal(writer.bufferedRecords)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}
		bufferedOutputStream.Write(out)
		writer.bufferedRecords = nil
	}
}

func (writer *RecordWriterYAML) writeWithoutListWrap(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
) {
	if outrec == nil {
		return
	}
	if writer.wroteAnyRecords {
		bufferedOutputStream.WriteString("---\n")
	}
	native, err := mlrval.MlrmapToYAMLNative(outrec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
		os.Exit(1)
	}
	out, err := yaml.Marshal(native)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
		os.Exit(1)
	}
	bufferedOutputStream.Write(out)
	writer.wroteAnyRecords = true
}
