package output

import (
	"bufio"

	"gopkg.in/yaml.v3"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type RecordWriterYAML struct {
	writerOptions   *cli.TWriterOptions
	bufferedRecords []*yaml.Node // used when WrapYAMLOutputInOuterList is true
	wroteAnyRecords bool         // for multi-doc: emit "---\n" before 2nd and later docs
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
		return writer.writeWithListWrap(outrec, bufferedOutputStream)
	}
	return writer.writeWithoutListWrap(outrec, bufferedOutputStream)
}

func (writer *RecordWriterYAML) writeWithListWrap(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
) error {
	if outrec != nil {
		if writer.bufferedRecords == nil {
			writer.bufferedRecords = []*yaml.Node{}
		}
		native, err := mlrval.MlrmapToYAMLNative(outrec)
		if err != nil {
			return err
		}
		writer.bufferedRecords = append(writer.bufferedRecords, native)
	} else {
		// End of stream: emit single YAML document as array
		seqNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		seqNode.Content = writer.bufferedRecords
		out, err := yaml.Marshal(seqNode)
		if err != nil {
			return err
		}
		bufferedOutputStream.Write(out)
		writer.bufferedRecords = nil
	}
	return nil
}

func (writer *RecordWriterYAML) writeWithoutListWrap(
	outrec *mlrval.Mlrmap,
	bufferedOutputStream *bufio.Writer,
) error {
	if outrec == nil {
		return nil
	}
	if writer.wroteAnyRecords {
		bufferedOutputStream.WriteString("---\n")
	}
	native, err := mlrval.MlrmapToYAMLNative(outrec)
	if err != nil {
		return err
	}
	out, err := yaml.Marshal(native)
	if err != nil {
		return err
	}
	bufferedOutputStream.Write(out)
	writer.wroteAnyRecords = true
	return nil
}
