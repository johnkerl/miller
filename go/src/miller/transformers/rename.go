package transformers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"miller/cliutil"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameRename = "rename"

var RenameSetup = transforming.TransformerSetup{
	Verb:         verbNameRename,
	ParseCLIFunc: transformerRenameParseCLI,
	UsageFunc:    transformerRenameUsage,
	IgnoresInput: false,
}

func transformerRenameParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerRenameUsage(os.Stdout, true, 0)

		} else {
			transformerRenameUsage(os.Stderr, true, 1)
		}
	}

	// Get the rename field names from the command line
	if argi >= argc {
		transformerRenameUsage(os.Stderr, true, 1)
	}
	names := lib.SplitString(args[argi], ",")
	if len(names)%2 != 0 {
		transformerRenameUsage(os.Stderr, true, 1)
	}
	argi++

	transformer, _ := NewTransformerRename(
		names,
	)

	*pargi = argi
	return transformer
}

func transformerRenameUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {old1,new1,old2,new2,...}\n", lib.MlrExeName(), verbNameRename)
	fmt.Fprintf(o, "Renames specified fields.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerRename struct {
	oldToNewNames *lib.OrderedMap
}

func NewTransformerRename(
	names []string,
) (*TransformerRename, error) {
	if len(names)%2 != 0 {
		return nil, errors.New("Rename: names string must have even length.")
	}

	oldToNewNames := lib.NewOrderedMap()
	n := len(names)
	for i := 0; i < n; i += 2 {
		oldName := names[i]
		newName := names[i+1]
		oldToNewNames.Put(oldName, newName)
	}

	this := &TransformerRename{
		oldToNewNames: oldToNewNames,
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerRename) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if this.oldToNewNames.Has(pe.Key) {
				newName := this.oldToNewNames.Get(pe.Key).(string)
				inrec.Rename(pe.Key, newName)
			}

		}
	}
	outputChannel <- inrecAndContext // including end-of-stream marker
}
