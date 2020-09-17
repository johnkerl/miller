package mappers

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var RenameSetup = mapping.MapperSetup{
	Verb:         "rename",
	ParseCLIFunc: mapperRenameParseCLI,
	IgnoresInput: false,
}

func mapperRenameParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {

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
		mapperRenameUsage(ostream, args[0], verb, flagSet)
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

	mapper, _ := NewMapperRename(
		names,
	)

	*pargi = argi
	return mapper
}

func mapperRenameUsage(
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
type MapperRename struct {
	oldToNewNames *lib.OrderedMap
}

func NewMapperRename(
	names []string,
) (*MapperRename, error) {
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

	this := &MapperRename{
		oldToNewNames: oldToNewNames,
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperRename) Map(
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
