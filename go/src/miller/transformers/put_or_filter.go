package transformers

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"miller/clitypes"
	"miller/dsl"
	"miller/dsl/cst"
	"miller/lib"
	"miller/parsing/lexer"
	"miller/parsing/parser"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var PutSetup = transforming.TransformerSetup{
	Verb:         "put",
	ParseCLIFunc: transformerPutParseCLI,
	IgnoresInput: false,
}

var FilterSetup = transforming.TransformerSetup{
	Verb:         "filter",
	ParseCLIFunc: transformerPutParseCLI,
	IgnoresInput: false,
}

func transformerPutParseCLI(
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

	dslString := ""
	verbose := false
	printASTOnly := false
	printASTSingleLine := false
	invertFilter := false
	suppressOutputRecord := false
	needExpressionArg := true
	presets := make([]string, 0)

	// Parse local flags.
	//
	// Unlike other transformers, we can't use flagSet here. The syntax of 'mlr put'
	// and 'mlr filter' is they need to be able to take -f and/or -e more than
	// once, and Go flags can't handle that.

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerPutUsage(os.Stdout, 0, errorHandling, args[0], verb)
			return nil // help intentionally requested

		} else if args[argi] == "-f" {
			checkArgCountPut(args, argi, argc, 2)

			// Get a DSL string from the user-specified filename
			data, err := ioutil.ReadFile(args[argi+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s %s: cannot load DSL expression: ", args[0], verb)
				fmt.Println(err)
				return nil
			}
			if dslString != "" {
				dslString += "\n"
			}
			dslString += string(data)

			needExpressionArg = false
			argi += 2

		} else if args[argi] == "-e" {
			checkArgCountPut(args, argi, argc, 2)

			// Get a DSL string from the user-specified filename
			if dslString != "" {
				dslString += ";\n"
			}
			dslString += args[argi+1]

			needExpressionArg = false
			argi += 2

		} else if args[argi] == "-s" {
			// E.g.
			//   mlr put -s sum=0
			// is like
			//   mlr put -s 'begin {@sum = 0}'
			checkArgCountPut(args, argi, argc, 2)
			presets = append(presets, args[argi+1])

			argi += 2

		} else if args[argi] == "-x" {
			invertFilter = true
			argi++
		} else if args[argi] == "-q" {
			suppressOutputRecord = true
			argi++

			// TODO: move these to mlr auxents?
		} else if args[argi] == "-d" {
			printASTOnly = true
			printASTSingleLine = false
			argi++
		} else if args[argi] == "-D" {
			printASTOnly = true
			printASTSingleLine = true
			argi++

		} else if args[argi] == "-v" {
			verbose = true
			argi++

		} else if args[argi] == "-S" {
			// TODO: this is a no-op in Miller 6 and above.
			// Comment this in more detail.
			argi++

		} else if args[argi] == "-F" {
			// TODO: this is a no-op in Miller 6 and above.
			// Comment this in more detail.
			argi++

		} else {
			transformerPutUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
			os.Exit(1)
		}
	}

	// If they've used either of 'mlr put -f {filename}' or 'mlr put -e
	// {expression}' then that specifies their DSL expression. But if they've
	// done neither then we expect 'mlr put {expression}'.
	if needExpressionArg {
		// Get the DSL string from the command line, after the flags
		if argi >= argc {
			transformerPutUsage(os.Stderr, 1, flag.ExitOnError, args[0], verb)
			os.Exit(1)
		}
		dslString = args[argi]
		argi += 1
	}

	if printASTOnly {
		astRootNode, err := BuildASTFromStringWithMessage(dslString, false)
		if err == nil {
			if printASTSingleLine {
				astRootNode.PrintParexOneLine()
			} else {
				astRootNode.PrintParex()
			}
			os.Exit(0)
		} else {
			// error message already printed out
			os.Exit(1)
		}
	}

	isFilter := verb == "filter"
	transformer, err := NewTransformerPut(
		dslString,
		isFilter,
		presets,
		verbose,
		invertFilter,
		suppressOutputRecord,
	)
	if err != nil {
		// Error message already printed out
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// For flags with values, e.g. ["-n" "10"], while we're looking at the "-n"
// this let us see if the "10" slot exists.
func checkArgCountPut(args []string, argi int, argc int, n int) {
	if (argc - argi) < n {
		fmt.Fprintf(os.Stderr, "%s: option \"%s\" missing argument(s).\n", args[0], args[argi])
		transformerPutUsage(os.Stderr, 1, flag.ExitOnError, os.Args[0], "sort")
		os.Exit(1)
	}
}

func transformerPutUsage(
	o *os.File,
	exitCode int,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	argv0 string,
	verb string,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {DSL expression}\n", argv0, verb)
	fmt.Fprintf(o, "TODO: put detailed on-line help here.\n")
	fmt.Fprintf(o,
		` -f {file name} File containing a DSL expression.
 -x (default false) Prints records for which {expression} evaluates to false, not true,
    i.e. invert the sense of the filter expression.
 -q (default false) Does not include the modified record in the output stream.
    Useful for when all desired output is in begin and/or end blocks.
 -v (default false) Prints the expressions's AST (abstract syntax tree), which gives
    full transparency on the precedence and associativity rules of
    Miller's grammar, to stdout.
`)
}

// ----------------------------------------------------------------
type TransformerPut struct {
	astRootNode          *dsl.AST
	cstRootNode          *cst.RootNode
	cstState             *cst.State
	callCount            int64
	invertFilter         bool
	suppressOutputRecord bool
	executedBeginBlocks  bool
}

func NewTransformerPut(
	dslString string,
	isFilter bool,
	presets []string,
	verbose bool,
	invertFilter bool,
	suppressOutputRecord bool,
) (*TransformerPut, error) {

	astRootNode, err := BuildASTFromStringWithMessage(dslString, verbose)
	if err != nil {
		// Error message already printed out
		return nil, err
	}

	if verbose {
		fmt.Println("DSL EXPRESSION:")
		fmt.Println(dslString)
		fmt.Println("RAW AST:")
		astRootNode.Print()
		fmt.Println()
	}
	cstRootNode, err := cst.Build(astRootNode, isFilter)
	cstState := cst.NewEmptyState()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

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
			mvalue := types.MlrvalFromInferredType(svalue)
			cstState.Oosvars.PutCopy(&key, &mvalue)
		}
	}

	return &TransformerPut{
		astRootNode:          astRootNode,
		cstRootNode:          cstRootNode,
		cstState:             cstState,
		callCount:            0,
		invertFilter:         invertFilter,
		suppressOutputRecord: suppressOutputRecord,
		executedBeginBlocks:  false,
	}, nil
}

func BuildASTFromStringWithMessage(dslString string, verbose bool) (*dsl.AST, error) {
	astRootNode, err := BuildASTFromString(dslString)
	if err != nil {
		// Leave this out until we get better control over the error-messaging.
		// At present it's overly parser-internal, and confusing. :(
		// fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "%s: cannot parse DSL expression.\n",
			os.Args[0])
		if verbose {
			fmt.Fprintln(os.Stderr, dslString)
		}
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
	this.cstState.OutputChannel = outputChannel

	inrec := inrecAndContext.Record
	context := inrecAndContext.Context
	if inrec != nil {

		// Execute the begin { ... } before the first input record
		this.callCount++
		if this.callCount == 1 {
			this.cstState.Update(nil, &context)
			err := this.cstRootNode.ExecuteBeginBlocks(this.cstState)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			this.executedBeginBlocks = true
		}

		this.cstState.Update(inrec, &context)

		// Execute the main block on the current input record
		outrec, err := this.cstRootNode.ExecuteMainBlock(this.cstState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if !this.suppressOutputRecord {
			wantToEmit := lib.BooleanXOR(this.cstState.FilterResult, this.invertFilter)
			if wantToEmit {
				outputChannel <- types.NewRecordAndContext(
					outrec,
					&context,
				)
			}
		}

	} else {
		this.cstState.Update(nil, &context)

		// If there were no input records then we never executed the
		// begin-blocks. Do so now.
		if this.executedBeginBlocks == false {
			err := this.cstRootNode.ExecuteBeginBlocks(this.cstState)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		// Execute the end { ... } after the last input record
		err := this.cstRootNode.ExecuteEndBlocks(this.cstState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		outputChannel <- types.NewRecordAndContext(
			nil, // signals end of input record stream
			&context,
		)

	}
}
