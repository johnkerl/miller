package transformers

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameRename = "rename"

var RenameSetup = TransformerSetup{
	Verb:         verbNameRename,
	UsageFunc:    transformerRenameUsage,
	ParseCLIFunc: transformerRenameParseCLI,
	IgnoresInput: false,
}

func transformerRenameUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	exeName := "mlr"
	verb := verbNameRename

	fmt.Fprintf(o, "Usage: %s %s [options] {old1,new1,old2,new2,...}\n", "mlr", verbNameRename)
	fmt.Fprintf(o, "Renames specified fields.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-r         Treat old field  names as regular expressions. \"ab\", \"a.*b\"\n")
	fmt.Fprintf(o, "           will match any field name containing the substring \"ab\" or\n")
	fmt.Fprintf(o, "           matching \"a.*b\", respectively; anchors of the form \"^ab$\",\n")
	fmt.Fprintf(o, "           \"^a.*b$\" may be used. New field names may be plain strings,\n")
	fmt.Fprintf(o, "           or may contain capture groups of the form \"\\1\" through\n")
	fmt.Fprintf(o, "           \"\\9\". Wrapping the regex in double quotes is optional, but\n")
	fmt.Fprintf(o, "           is required if you wish to follow it with 'i' to indicate\n")
	fmt.Fprintf(o, "           case-insensitivity.\n")
	fmt.Fprintf(o, "-g         Do global replacement within each field name rather than\n")
	fmt.Fprintf(o, "           first-match replacement.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "%s %s old_name,new_name'\n", exeName, verb)
	fmt.Fprintf(o, "%s %s old_name_1,new_name_1,old_name_2,new_name_2'\n", exeName, verb)
	fmt.Fprintf(o, "%s %s -r 'Date_[0-9]+,Date,'  Rename all such fields to be \"Date\"\n", exeName, verb)
	fmt.Fprintf(o, "%s %s -r '\"Date_[0-9]+\",Date' Same\n", exeName, verb)
	fmt.Fprintf(o, "%s %s -r 'Date_([0-9]+).*,\\1' Rename all such fields to be of the form 20151015\n", exeName, verb)
	fmt.Fprintf(o, "%s %s -r '\"name\"i,Name'       Rename \"name\", \"Name\", \"NAME\", etc. to \"Name\"\n", exeName, verb)

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerRenameParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	doRegexes := false
	doGsub := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerRenameUsage(os.Stdout, true, 0)

		} else if opt == "-r" {
			doRegexes = true

		} else if opt == "-g" {
			doGsub = true

		} else {
			transformerRenameUsage(os.Stderr, true, 1)
		}
	}

	if doGsub {
		doRegexes = true
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

	transformer, err := NewTransformerRename(
		names,
		doRegexes,
		doGsub,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type tRegexAndReplacement struct {
	regex                    *regexp.Regexp
	replacement              string
	replacementCaptureMatrix [][]int // TODO: comment
}

type TransformerRename struct {
	oldToNewNames          *lib.OrderedMap
	regexesAndReplacements *list.List
	doGsub                 bool
	recordTransformerFunc  RecordTransformerFunc
}

func NewTransformerRename(
	names []string,
	doRegexes bool,
	doGsub bool,
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

	tr := &TransformerRename{}

	if !doRegexes {
		tr.oldToNewNames = oldToNewNames
		tr.doGsub = false
		tr.recordTransformerFunc = tr.transformWithoutRegexes
	} else {
		tr.regexesAndReplacements = list.New()
		for pe := oldToNewNames.Head; pe != nil; pe = pe.Next {
			regexString := pe.Key
			regex := lib.CompileMillerRegexOrDie(regexString)
			replacement := pe.Value.(string)
			_, replacementCaptureMatrix := lib.RegexReplacementHasCaptures(replacement)
			regexAndReplacement := tRegexAndReplacement{
				regex:                    regex,
				replacement:              replacement,
				replacementCaptureMatrix: replacementCaptureMatrix,
			}
			tr.regexesAndReplacements.PushBack(&regexAndReplacement)
		}
		tr.doGsub = doGsub
		tr.recordTransformerFunc = tr.transformWithRegexes
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerRename) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerRename) transformWithoutRegexes(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.oldToNewNames.Has(pe.Key) {
				newName := tr.oldToNewNames.Get(pe.Key).(string)
				inrec.Rename(pe.Key, newName)
			}

		}
	}
	outputChannel <- inrecAndContext // including end-of-stream marker
}

// ----------------------------------------------------------------
func (tr *TransformerRename) transformWithRegexes(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for pr := tr.regexesAndReplacements.Front(); pr != nil; pr = pr.Next() {
			regexAndReplacement := pr.Value.(*tRegexAndReplacement)
			regex := regexAndReplacement.regex
			replacement := regexAndReplacement.replacement
			replacementCaptureMatrix := regexAndReplacement.replacementCaptureMatrix

			for pe := inrec.Head; pe != nil; pe = pe.Next {
				oldName := pe.Key
				if tr.doGsub {
					newName := regex.ReplaceAllString(oldName, replacement)
					if newName != oldName {
						inrec.Rename(oldName, newName)
					}
				} else {
					newName := lib.RegexSubCompiled(oldName, regex, replacement, replacementCaptureMatrix)
					if newName != oldName {
						inrec.Rename(oldName, newName)
					}
				}
			}
		}

		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // including end-of-stream marker
	}
}
