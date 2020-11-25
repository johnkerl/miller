package transformers

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var RenameSetup = transforming.TransformerSetup{
	Verb:         "rename",
	ParseCLIFunc: transformerRenameParseCLI,
	IgnoresInput: false,
}

func transformerRenameParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerRenameUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	// Get the rename field names from the command line
	if argi >= argc {
		flagSet.Usage()
		os.Exit(1)
	}
	names := lib.SplitString(args[argi], ",")
	if len(names)%2 != 0 {
		flagSet.Usage()
		os.Exit(1)
	}

	argi += 1

	transformer, _ := NewTransformerRename(
		names,
	)

	*pargi = argi
	return transformer
}

func transformerRenameUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {old1,new1,old2,new2,...}\n", argv0, verb)
	fmt.Fprintf(o, "Renames specified fields.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
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
func (this *TransformerRename) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if this.oldToNewNames.Has(*pe.Key) {
				newName := this.oldToNewNames.Get(*pe.Key).(string)
				inrec.Rename(pe.Key, &newName)
			}

		}
	}
	outputChannel <- inrecAndContext // end-of-stream marker
}
