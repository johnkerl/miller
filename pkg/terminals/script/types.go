// Data types for the Miller script terminal.

package script

import (
	"bufio"
	"os"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/dsl/cst"
	"github.com/johnkerl/miller/v6/pkg/input"
	"github.com/johnkerl/miller/v6/pkg/output"
	"github.com/johnkerl/miller/v6/pkg/runtime"
	"github.com/johnkerl/miller/v6/pkg/types"
)

type Script struct {
	exeName string
	name    string

	doWarnings  bool
	cstRootNode *cst.RootNode

	options *cli.TOptions

	readerChannel         chan []*types.RecordAndContext
	errorChannel          chan error
	downstreamDoneChannel chan bool
	recordReader          input.IRecordReader
	recordWriter          output.IRecordWriter

	recordOutputFileName       string
	recordOutputStream         *os.File
	bufferedRecordOutputStream *bufio.Writer

	runtimeState *runtime.State
}
