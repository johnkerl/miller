package transformers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/dsl"
	"miller/src/dsl/cst"
	"miller/src/lib"
	"miller/src/parsing/lexer"
	"miller/src/parsing/parser"
	"miller/src/runtime"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNamePut = "put"

var PutSetup = transforming.TransformerSetup{
	Verb:         verbNamePut,
	UsageFunc:    transformerPutUsage,
	ParseCLIFunc: transformerPutOrFilterParseCLI,
	IgnoresInput: false,
}

const verbNameFilter = "filter"

var FilterSetup = transforming.TransformerSetup{
	Verb:         verbNameFilter,
	UsageFunc:    transformerFilterUsage,
	ParseCLIFunc: transformerPutOrFilterParseCLI,
	IgnoresInput: false,
}

// ----------------------------------------------------------------
func transformerPutUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	transformerPutOrFilterUsage(o, doExit, exitCode, "put")
}

func transformerFilterUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	transformerPutOrFilterUsage(o, doExit, exitCode, "filter")
}

func transformerPutOrFilterUsage(
	o *os.File,
	doExit bool,
	exitCode int,
	verb string,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {DSL expression}\n", lib.MlrExeName(), verb)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o,
		`-f {file name} File containing a DSL expression.

-e {expression} You can use this after -f to add an expression. Example use
   case: define functions/subroutines in a file you specify with -f, then call
   them with an expression you specify with -e.

(If you mix -e and -f then the expressions are evaluated in the order encountered.
Since the expression pieces are simply concatenated, please be sure to use intervening
semicolons to separate expressions.)

-s name=value: Predefines out-of-stream variable @name to have 
    Thus mlr put -s foo=97 '$column += @foo' is like
    mlr put 'begin {@foo = 97} $column += @foo'.
    The value part is subject to type-inferencing.
    May be specified more than once, e.g. -s name1=value1 -s name2=value2.
    Note: the value may be an environment variable, e.g. -s sequence=$SEQUENCE

-x (default false) Prints records for which {expression} evaluates to false, not true,
   i.e. invert the sense of the filter expression.

-q Does not include the modified record in the output stream.
   Useful for when all desired output is in begin and/or end blocks.

-p Prints the expressions's AST (abstract syntax tree), which gives
   full transparency on the precedence and associativity rules of
   Miller's grammar, to stdout.

-d Like -p but uses a parenthesized-expression format for the AST.

-D Like -d but with output all on one line.

-E Echo DSL expression before printing parse-tree

-v Same as -E -p.

-X Exit after parsing but before stream-processing. Useful with -v/-d/-D, if you
   only want to look at parser information.

-S and -F: There are no-ops in Miller 6 and above, since now type-inferencing is done
   by the record-readers before filter/put is executed. Supported as no-op pass-through
   flags for backward compatibility.

-h|--help Show this message.
`)

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
func transformerPutOrFilterParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	mainRecordWriterOptions *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var dslStrings []string = make([]string, 0)
	echoDSLString := false
	printASTAsTree := false
	printASTMultiLine := false
	printASTSingleLine := false
	exitAfterParse := false
	invertFilter := false
	suppressOutputRecord := false
	presets := make([]string, 0)

	// TODO: make sure this is a full nested-struct copy.
	var recordWriterOptions *cliutil.TWriterOptions = nil
	if mainRecordWriterOptions != nil {
		// TODO: make a .Copy() method
		copyThereof := *mainRecordWriterOptions
		recordWriterOptions = &copyThereof
	}

	// Parse local flags.

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerPutUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			// Get a DSL string from the user-specified filename
			filename := cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s %s: cannot load DSL expression from file \"%s\": ",
					lib.MlrExeName(), verb, filename)
				fmt.Println(err)
				return nil
			}
			dslString := string(data)
			dslStrings = append(dslStrings, dslString)

		} else if opt == "-e" {
			dslString := cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			dslStrings = append(dslStrings, dslString)

		} else if opt == "-s" {
			// E.g.
			//   mlr put -s sum=0
			// is like
			//   mlr put -s 'begin {@sum = 0}'
			preset := cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			presets = append(presets, preset)

		} else if opt == "-x" {
			invertFilter = true
		} else if opt == "-q" {
			suppressOutputRecord = true

		} else if opt == "-E" {
			echoDSLString = true
		} else if opt == "-p" {
			printASTAsTree = true
		} else if opt == "-v" {
			echoDSLString = true
			printASTAsTree = true
		} else if opt == "-d" {
			printASTMultiLine = true
		} else if opt == "-D" {
			printASTSingleLine = true
		} else if opt == "-X" {
			exitAfterParse = true

		} else if opt == "-S" {
			// TODO: this is a no-op in Miller 6 and above.
			// Comment this in more detail.

		} else if opt == "-F" {
			// TODO: this is a no-op in Miller 6 and above.
			// Comment this in more detail.

		} else {
			// This is inelegant. For error-proofing we advance argi already in our
			// loop (so individual if-statements don't need to). However,
			// ParseWriterOptions expects it unadvanced.
			wargi := argi - 1
			if cliutil.ParseWriterOptions(args, argc, &wargi, recordWriterOptions) {
				// This lets mlr main and mlr put have different output formats.
				// Nothing else to handle here.
				argi = wargi
			} else {
				transformerPutUsage(os.Stderr, true, 1)
			}
		}
	}

	// If they've used either of 'mlr put -f {filename}' or 'mlr put -e
	// {expression}' then that specifies their DSL expression. But if they've
	// done neither then we expect 'mlr put {expression}'.
	if len(dslStrings) == 0 {
		// Get the DSL string from the command line, after the flags
		if argi >= argc {
			transformerPutUsage(os.Stderr, true, 1)
		}
		dslString := args[argi]
		dslStrings = append(dslStrings, dslString)
		argi++
	}

	isFilter := verb == "filter"
	transformer, err := NewTransformerPut(
		dslStrings,
		isFilter,
		presets,
		echoDSLString,
		printASTAsTree,
		printASTMultiLine,
		printASTSingleLine,
		exitAfterParse,
		invertFilter,
		suppressOutputRecord,
		recordWriterOptions,
	)
	if err != nil {
		// Error message already printed out
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerPut struct {
	cstRootNode          *cst.RootNode
	runtimeState         *runtime.State
	callCount            int
	invertFilter         bool
	suppressOutputRecord bool
	executedBeginBlocks  bool
}

func NewTransformerPut(
	dslStrings []string,
	isFilter bool,
	presets []string,
	echoDSLString bool,
	printASTAsTree bool,
	printASTMultiLine bool,
	printASTSingleLine bool,
	exitAfterParse bool,
	invertFilter bool,
	suppressOutputRecord bool,
	recordWriterOptions *cliutil.TWriterOptions,
) (*TransformerPut, error) {

	cstRootNode := cst.NewEmptyRoot(recordWriterOptions)

	for _, dslString := range dslStrings {
		astRootNode, err := BuildASTFromStringWithMessage(dslString)
		if err != nil {
			// Error message already printed out
			return nil, err
		}

		if echoDSLString {
			fmt.Println("DSL EXPRESSION:")
			fmt.Println(dslString)
		}
		if printASTAsTree {
			fmt.Println("AST:")
			astRootNode.Print()
			fmt.Println()
		}
		if printASTMultiLine {
			astRootNode.PrintParex()
			fmt.Println()
		}
		if printASTSingleLine {
			astRootNode.PrintParexOneLine()
			fmt.Println()
		}
		if exitAfterParse {
			os.Exit(0)
		}

		err = cstRootNode.IngestAST(astRootNode, isFilter, false) // TODO: split out methods ...
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, err
		}
	}

	err := cstRootNode.Resolve()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	runtimeState := runtime.NewEmptyState()

	// E.g.
	//   mlr put -s sum=0
	// is like
	//   mlr put -s 'begin {@sum = 0}'
	if len(presets) > 0 {
		for _, preset := range presets {
			pair := strings.SplitN(preset, "=", 2)
			if len(pair) != 2 {
				return nil, errors.New(
					fmt.Sprintf(
						"Miller: missing \"=\" in preset expression \"%s\".",
						preset,
					),
				)
			}
			key := pair[0]
			svalue := pair[1]
			mvalue := types.MlrvalPointerFromInferredType(svalue)
			runtimeState.Oosvars.PutCopy(key, mvalue)
		}
	}

	return &TransformerPut{
		cstRootNode:          cstRootNode,
		runtimeState:         runtimeState,
		callCount:            0,
		invertFilter:         invertFilter,
		suppressOutputRecord: suppressOutputRecord,
		executedBeginBlocks:  false,
	}, nil
}

func BuildASTFromStringWithMessage(dslString string) (*dsl.AST, error) {
	astRootNode, err := BuildASTFromString(dslString)
	if err != nil {
		// Leave this out until we get better control over the error-messaging.
		// At present it's overly parser-internal, and confusing. :(
		// fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "%s: cannot parse DSL expression.\n",
			lib.MlrExeName())
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	} else {
		return astRootNode, nil
	}
}

// xxx note (package cycle) why not a dsl.AST constructor :(
// xxx maybe split out dsl into two packages ... and/or put the ast.go into miller/parsing -- ?
//   depends on TBD split-out of AST and CST ...
func BuildASTFromString(dslString string) (*dsl.AST, error) {
	theLexer := lexer.NewLexer([]byte(dslString))
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err != nil {
		return nil, err
	}
	astRootNode := interfaceAST.(*dsl.AST)
	return astRootNode, nil
}

func (this *TransformerPut) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.runtimeState.OutputChannel = outputChannel

	inrec := inrecAndContext.Record
	context := inrecAndContext.Context
	if !inrecAndContext.EndOfStream {

		// Execute the begin { ... } before the first input record
		this.callCount++
		if this.callCount == 1 {
			this.runtimeState.Update(nil, &context)
			err := this.cstRootNode.ExecuteBeginBlocks(this.runtimeState)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			this.executedBeginBlocks = true
		}

		this.runtimeState.Update(inrec, &context)

		// Execute the main block on the current input record
		outrec, err := this.cstRootNode.ExecuteMainBlock(this.runtimeState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !this.suppressOutputRecord {
			filterBool, isBool := this.runtimeState.FilterExpression.GetBoolValue()
			if !isBool {
				filterBool = false
			}
			wantToEmit := lib.BooleanXOR(filterBool, this.invertFilter)
			if wantToEmit {
				outputChannel <- types.NewRecordAndContext(
					outrec,
					&context,
				)
			}
		}

	} else {
		this.runtimeState.Update(nil, &context)

		// If there were no input records then we never executed the
		// begin-blocks. Do so now.
		if this.executedBeginBlocks == false {
			err := this.cstRootNode.ExecuteBeginBlocks(this.runtimeState)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		// Execute the end { ... } after the last input record
		err := this.cstRootNode.ExecuteEndBlocks(this.runtimeState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Send all registered OutputHandlerManager instances the end-of-stream
		// indicator.
		this.cstRootNode.ProcessEndOfStream()

		outputChannel <- types.NewEndOfStreamMarker(&context)
	}
}
